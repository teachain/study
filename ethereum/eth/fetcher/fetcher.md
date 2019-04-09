```
// Copyright 2015 The go-simplechain Authors
// This file is part of the go-simplechain library.
//
// The go-simplechain library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-simplechain library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-simplechain library. If not, see <http://www.gnu.org/licenses/>.

// Package fetcher contains the block announcement based synchronisation.
package fetcher

import (
	"errors"
	"math/rand"
	"time"

	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/consensus"
	"github.com/simplechain-org/go-simplechain/core/types"
	"github.com/simplechain-org/go-simplechain/log"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

const (
	arriveTimeout = 500 * time.Millisecond // Time allowance before an announced block is explicitly requested
	gatherSlack   = 100 * time.Millisecond // Interval used to collate almost-expired announces with fetches
	fetchTimeout  = 5 * time.Second        // Maximum allotted time to return an explicitly requested block
	maxUncleDist  = 7                      // Maximum allowed backward distance from the chain head
	maxQueueDist  = 32                     // Maximum allowed distance from the chain head to queue
	hashLimit     = 256                    // Maximum number of unique blocks a peer may have announced
	blockLimit    = 64                     // Maximum number of unique blocks a peer may have delivered
)

var (
	errTerminated = errors.New("terminated")
)

//从本地区块链中检索一个区块（根据区块哈希获取区块）
// blockRetrievalFn is a callback type for retrieving a block from the local chain.
type blockRetrievalFn func(common.Hash) *types.Block

//发送一个区块头部检索请求
// headerRequesterFn is a callback type for sending a header retrieval request.
type headerRequesterFn func(common.Hash) error

//发送一个区块体检索请求
// bodyRequesterFn is a callback type for sending a body retrieval request.
type bodyRequesterFn func([]common.Hash) error

//区块头部校验函数（快速传播时使用）
// headerVerifierFn is a callback type to verify a block's header for fast propagation.
type headerVerifierFn func(header *types.Header) error

//区块广播函数
// blockBroadcasterFn is a callback type for broadcasting a block to connected peers.
type blockBroadcasterFn func(block *types.Block, propagate bool)

//获取区块高度函数
// chainHeightFn is a callback type to retrieve the current chain height.
type chainHeightFn func() uint64

//将区块插入链函数
// chainInsertFn is a callback type to insert a batch of blocks into the local chain.
type chainInsertFn func(types.Blocks) (int, error)

//节点掉线函数
// peerDropFn is a callback type for dropping a peer detected as malicious.
type peerDropFn func(id string)

// announce is the hash notification of the availability of a new block in the
// network.
//公告，声称，就是广而告之，我有新区块了，快来获取取吧
type announce struct {
	hash   common.Hash   // Hash of the block being announced
	number uint64        // Number of the block being announced (0 = unknown | old protocol)
	header *types.Header // Header of the block partially reassembled (new protocol)
	time   time.Time     // Timestamp of the announcement

	origin string // Identifier of the peer originating the notification

	//获取区块头
	fetchHeader headerRequesterFn // Fetcher function to retrieve the header of an announced block

	//获取区块体
	fetchBodies bodyRequesterFn   // Fetcher function to retrieve the body of an announced block
}

// headerFilterTask represents a batch of headers needing fetcher filtering.
type headerFilterTask struct {
	peer    string          // The source peer of block headers
	headers []*types.Header // Collection of headers to filter
	time    time.Time       // Arrival time of the headers
}

// bodyFilterTask represents a batch of block bodies (transactions and uncles)
// needing fetcher filtering.
type bodyFilterTask struct {
	peer         string                 // The source peer of block bodies
	transactions [][]*types.Transaction // Collection of transactions per block bodies
	uncles       [][]*types.Header      // Collection of uncles per block bodies
	time         time.Time              // Arrival time of the blocks' contents
}

// inject represents a schedules import operation.
type inject struct {
	origin string
	block  *types.Block
}

// Fetcher is responsible for accumulating block announcements from various peers
// and scheduling them for retrieval.
type Fetcher struct {
	// Various event channels
	notify chan *announce
	inject chan *inject

	blockFilter  chan chan []*types.Block
	headerFilter chan chan *headerFilterTask
	bodyFilter   chan chan *bodyFilterTask

	done chan common.Hash
	quit chan struct{}

	// Announce states
	announces  map[string]int              // Per peer announce counts to prevent memory exhaustion
	announced  map[common.Hash][]*announce // Announced blocks, scheduled for fetching
	fetching   map[common.Hash]*announce   // Announced blocks, currently fetching
	fetched    map[common.Hash][]*announce // Blocks with headers fetched, scheduled for body retrieval
	completing map[common.Hash]*announce   // Blocks with headers, currently body-completing


	// Block cache
	//收到的区块队列（通过网络过来的）
	queue  *prque.Prque            // Queue containing the import operations (block number sorted)
	queues map[string]int          // Per peer block counts to prevent memory exhaustion
	queued map[common.Hash]*inject // Set of already queued blocks (to dedupe imports)

	// Callbacks
	getBlock       blockRetrievalFn   // Retrieves a block from the local chain
	verifyHeader   headerVerifierFn   // Checks if a block's headers have a valid proof of work
	broadcastBlock blockBroadcasterFn // Broadcasts a block to connected peers
	chainHeight    chainHeightFn      // Retrieves the current chain's height
	insertChain    chainInsertFn      // Injects a batch of blocks into the chain
	dropPeer       peerDropFn         // Drops a peer for misbehaving

	// Testing hooks
	announceChangeHook func(common.Hash, bool) // Method to call upon adding or deleting a hash from the announce list
	queueChangeHook    func(common.Hash, bool) // Method to call upon adding or deleting a block from the import queue
	fetchingHook       func([]common.Hash)     // Method to call upon starting a block (eth/61) or header (eth/62) fetch
	completingHook     func([]common.Hash)     // Method to call upon starting a block body fetch (eth/62)
	importedHook       func(*types.Block)      // Method to call upon successful block import (both eth/61 and eth/62)
}

// New creates a block fetcher to retrieve blocks based on hash announcements.
func New(getBlock blockRetrievalFn, verifyHeader headerVerifierFn, broadcastBlock blockBroadcasterFn, chainHeight chainHeightFn, insertChain chainInsertFn, dropPeer peerDropFn) *Fetcher {
	return &Fetcher{
		notify:         make(chan *announce),
		inject:         make(chan *inject),
		blockFilter:    make(chan chan []*types.Block),
		headerFilter:   make(chan chan *headerFilterTask),
		bodyFilter:     make(chan chan *bodyFilterTask),
		done:           make(chan common.Hash),
		quit:           make(chan struct{}),
		announces:      make(map[string]int),
		announced:      make(map[common.Hash][]*announce),
		fetching:       make(map[common.Hash]*announce),
		fetched:        make(map[common.Hash][]*announce),
		completing:     make(map[common.Hash]*announce),
		queue:          prque.New(),
		queues:         make(map[string]int),
		queued:         make(map[common.Hash]*inject),
		getBlock:       getBlock,
		verifyHeader:   verifyHeader,
		broadcastBlock: broadcastBlock,
		chainHeight:    chainHeight,
		insertChain:    insertChain,
		dropPeer:       dropPeer,
	}
}

// Start boots up the announcement based synchroniser, accepting and processing
// hash notifications and block fetches until termination requested.
func (f *Fetcher) Start() {
	go f.loop()
}

// Stop terminates the announcement based synchroniser, canceling all pending
// operations.
func (f *Fetcher) Stop() {
	close(f.quit)
}

// 得到NewBlockHashesMsg消息以后，就调用该方法
// 意思是有新区块了，给你这个区块的哈希，它在哪里，它的高度是多少，
// 用什么方法请求区块头，用什么方法请求区块体，都告诉你
// Notify announces the fetcher of the potential availability of a new block in
// the network.
func (f *Fetcher) Notify(peer string, hash common.Hash, number uint64, time time.Time,
	headerFetcher headerRequesterFn, bodyFetcher bodyRequesterFn) error {
	block := &announce{
		hash:        hash,
		number:      number,
		time:        time,
		origin:      peer,
		fetchHeader: headerFetcher,
		fetchBodies: bodyFetcher,
	}
	select {
	case f.notify <- block:
		return nil
	case <-f.quit:
		return errTerminated
	}
}

//也就是收到NewBlockMsg之后执行的流程
//这里是直接收到整个区块的流程
// Enqueue tries to fill gaps the the fetcher's future import queue.
// 想当于Enqueue-->enqueue-->f.queue.Push-->loop-->insert
func (f *Fetcher) Enqueue(peer string, block *types.Block) error {
	op := &inject{
		origin: peer,
		block:  block,
	}
	select {
	case f.inject <- op:
		return nil
	case <-f.quit:
		return errTerminated
	}
}

//发送获取区块头请求以后，得到响应以后，调用该方法
// FilterHeaders extracts all the headers that were explicitly requested by the fetcher,
// returning those that should be handled differently.
func (f *Fetcher) FilterHeaders(peer string, headers []*types.Header, time time.Time) []*types.Header {

	log.Trace("Filtering headers", "peer", peer, "headers", len(headers))

	// Send the filter channel to the fetcher
	filter := make(chan *headerFilterTask)

	select {
	case f.headerFilter <- filter:
	case <-f.quit:
		return nil
	}
	// Request the filtering of the header list
	select {
	case filter <- &headerFilterTask{peer: peer, headers: headers, time: time}:
	case <-f.quit:
		return nil
	}
	// Retrieve the headers remaining after filtering
	select {
	case task := <-filter:
		return task.headers
	case <-f.quit:
		return nil
	}
}

//发送请求区块体以后，得到响应以后，调用该方法
// FilterBodies extracts all the block bodies that were explicitly requested by
// the fetcher, returning those that should be handled differently.
func (f *Fetcher) FilterBodies(peer string, transactions [][]*types.Transaction, uncles [][]*types.Header, time time.Time) ([][]*types.Transaction, [][]*types.Header) {
	log.Trace("Filtering bodies", "peer", peer, "txs", len(transactions), "uncles", len(uncles))

	// Send the filter channel to the fetcher
	filter := make(chan *bodyFilterTask)

	select {
	case f.bodyFilter <- filter:
	case <-f.quit:
		return nil, nil
	}
	// Request the filtering of the body list
	select {
	case filter <- &bodyFilterTask{peer: peer, transactions: transactions, uncles: uncles, time: time}:
	case <-f.quit:
		return nil, nil
	}
	// Retrieve the bodies remaining after filtering
	select {
	case task := <-filter:
		return task.transactions, task.uncles
	case <-f.quit:
		return nil, nil
	}
}

// Loop is the main fetcher loop, checking and processing various notification
// events.
func (f *Fetcher) loop() {
	// Iterate the block fetching until a quit is requested
	fetchTimer := time.NewTimer(0)
	completeTimer := time.NewTimer(0)

	for {
		// Clean up any expired block fetches
		for hash, announce := range f.fetching {
			if time.Since(announce.time) > fetchTimeout {
				f.forgetHash(hash)
			}
		}
		//获取本地区块链的高度
		// Import any queued blocks that could potentially fit
		height := f.chainHeight()
		//队列不为空时
		for !f.queue.Empty() {
			//出队
			op := f.queue.PopItem().(*inject)
            //区块哈希
			hash := op.block.Hash()
            //队列有改变，调用钩子函数
			if f.queueChangeHook != nil {
				f.queueChangeHook(hash, false)
			}
			// If too high up the chain or phase, continue later
			//获取区块的高度
			number := op.block.NumberU64()
			//期望是得到下一个区块，如果不是并且它是一个未来块，那么我们将它放回队列中
			if number > height+1 {
				//将区块放回队列中
				f.queue.Push(op, -float32(number))
				//测试函数不为空，调用钩子测试函数
				if f.queueChangeHook != nil {
					f.queueChangeHook(hash, true)
				}
				//退出当前的for循环
				break
			}
			// Otherwise if fresh and still unknown, try and import
			//区块太老，是一个很前面的区块或者是已知的区块
			// （f.getBlock(hash) != nil也就是表示本地已经存在该区块了）
			//那么我们就忽略该区块
			if number+maxUncleDist < height || f.getBlock(hash) != nil {
				f.forgetBlock(hash)
				continue
			}
			//否则我们就走插入流程
			f.insert(op.origin, op.block)
		}
		// Wait for an outside event to occur
		select {
		case <-f.quit:
			// Fetcher terminating, abort all operations
			return

		case notification := <-f.notify:
			// A block was announced, make sure the peer isn't DOSing us
			propAnnounceInMeter.Mark(1)

		    //预计算来自这个节点的通知
			count := f.announces[notification.origin] + 1

			//如果这个节点的通知数大于给定的限制数（hashLimit     = 256 ）,认为它是dos攻击
			if count > hashLimit {
				log.Debug("Peer exceeded outstanding announces", "peer", notification.origin, "limit", hashLimit)
				propAnnounceDOSMeter.Mark(1)
				break
			}
			// If we have a valid block number, check that it's potentially useful
			if notification.number > 0 {
				//maxUncleDist  = 7
				//maxQueueDist  = 32
				//定义有效的区块为以当前区块为基点，-7<x<32为有效区块
				if dist := int64(notification.number) - int64(f.chainHeight()); dist < -maxUncleDist || dist > maxQueueDist {
					log.Debug("Peer discarded announcement", "peer", notification.origin, "number", notification.number, "hash", notification.hash, "distance", dist)
					propAnnounceDropMeter.Mark(1)
					break
				}
			}
			//如果是已知的话，显然不需要记录下来
			// All is well, schedule the announce if block's not yet downloading
			if _, ok := f.fetching[notification.hash]; ok {
				break
			}
			if _, ok := f.completing[notification.hash]; ok {
				break
			}
			//改变这个节点的通知数
			f.announces[notification.origin] = count

			//将此通知记录下来（也就是说将该通知放入了我（fetcher）的收集通知的列表中）
			f.announced[notification.hash] = append(f.announced[notification.hash], notification)

			//钩子函数，用于测试时观察
			if f.announceChangeHook != nil && len(f.announced[notification.hash]) == 1 {
				f.announceChangeHook(notification.hash, true)
			}
			//如果只存在一个通知
			if len(f.announced) == 1 {
				f.rescheduleFetch(fetchTimer)
			}

		case op := <-f.inject:
			// A direct block insertion was requested, try and fill any pending gaps
			propBroadcastInMeter.Mark(1)
			f.enqueue(op.origin, op.block)

		case hash := <-f.done:
			// A pending import finished, remove all traces of the notification
			f.forgetHash(hash)
			f.forgetBlock(hash)

		case <-fetchTimer.C:
			// At least one block's timer ran out, check for needing retrieval
			request := make(map[string][]common.Hash)

			for hash, announces := range f.announced {
				//arriveTimeout = 500 * time.Millisecond
				//gatherSlack   = 100 * time.Millisecond
				//如果说时间已经过去了400毫秒
				if time.Since(announces[0].time) > arriveTimeout-gatherSlack {

					// Pick a random peer to retrieve from, reset all others
					announce := announces[rand.Intn(len(announces))]

					f.forgetHash(hash)

					// If the block still didn't arrive, queue for fetching
					if f.getBlock(hash) == nil {
						request[announce.origin] = append(request[announce.origin], hash)
						f.fetching[hash] = announce
					}
				}
			}
			// Send out all block header requests
			for peer, hashes := range request {

				log.Trace("Fetching scheduled headers", "peer", peer, "list", hashes)

				// Create a closure of the fetch and schedule in on a new thread
				fetchHeader, hashes := f.fetching[hashes[0]].fetchHeader, hashes

				go func() {
					//钩子函数，测试使用
					if f.fetchingHook != nil {
						f.fetchingHook(hashes)
					}
					for _, hash := range hashes {
						headerFetchMeter.Mark(1)
						//发送区块头请求
						fetchHeader(hash) // Suboptimal, but protocol doesn't allow batch header retrievals
					}
				}()
			}
			// Schedule the next fetch if blocks are still pending
			f.rescheduleFetch(fetchTimer)

		case <-completeTimer.C:
			// At least one header's timer ran out, retrieve everything
			request := make(map[string][]common.Hash)

			for hash, announces := range f.fetched {

				// Pick a random peer to retrieve from, reset all others
				announce := announces[rand.Intn(len(announces))]

				f.forgetHash(hash)

				// If the block still didn't arrive, queue for completion
				if f.getBlock(hash) == nil {
					request[announce.origin] = append(request[announce.origin], hash)
					f.completing[hash] = announce
				}
			}
			// Send out all block body requests
			for peer, hashes := range request {

				log.Trace("Fetching scheduled bodies", "peer", peer, "list", hashes)

				// Create a closure of the fetch and schedule in on a new thread
				if f.completingHook != nil {
					f.completingHook(hashes)
				}
				bodyFetchMeter.Mark(int64(len(hashes)))

				//请求区块体
				go f.completing[hashes[0]].fetchBodies(hashes)

			}
			// Schedule the next fetch if blocks are still pending
			f.rescheduleComplete(completeTimer)

		case filter := <-f.headerFilter:
			// Headers arrived from a remote peer. Extract those that were explicitly
			// requested by the fetcher, and return everything else so it's delivered
			// to other parts of the system.
			var task *headerFilterTask
			select {
			case task = <-filter:
			case <-f.quit:
				return
			}

			headerFilterInMeter.Mark(int64(len(task.headers)))

			// Split the batch of headers into unknown ones (to return to the caller),
			// known incomplete ones (requiring body retrievals) and completed blocks.
			unknown, incomplete, complete := []*types.Header{}, []*announce{}, []*types.Block{}

			for _, header := range task.headers {

				hash := header.Hash()

				// Filter fetcher-requested headers from other synchronisation algorithms
				if announce := f.fetching[hash]; announce != nil && announce.origin == task.peer && f.fetched[hash] == nil && f.completing[hash] == nil && f.queued[hash] == nil {
					//对不上号，当然不处理哦
					// If the delivered header does not match the promised number, drop the announcer
					if header.Number.Uint64() != announce.number {
						log.Trace("Invalid block number fetched", "peer", announce.origin, "hash", header.Hash(), "announced", announce.number, "provided", header.Number)
						f.dropPeer(announce.origin)
						f.forgetHash(hash)
						continue
					}
					//本地还没有该区块
					// Only keep if not imported by other means
					if f.getBlock(hash) == nil {

						//注意这里我们把区块头记录下来了
						announce.header = header

						announce.time = task.time
						//这if执行的逻辑是指，当这个区块中没有交易（就是这个区块打包的时候，人家都没有发生过交易）
						//这个时候我们就没有必要再去要求取区块体了（因为它为空，去取没有啥意义）
						// If the block is empty (header only), short circuit into the final import queue
						if header.TxHash == types.DeriveSha(types.Transactions{}) && header.UncleHash == types.CalcUncleHash([]*types.Header{}) {
							log.Trace("Block empty, skipping body retrieval", "peer", announce.origin, "number", header.Number, "hash", header.Hash())

							block := types.NewBlockWithHeader(header)
							block.ReceivedAt = task.time

							complete = append(complete, block)
							f.completing[hash] = announce
							continue
						}
						//表示我需要去请求获取区块体
						// Otherwise add to the list of blocks needing completion
						incomplete = append(incomplete, announce)
					} else {
						log.Trace("Block already imported, discarding header", "peer", announce.origin, "number", header.Number, "hash", header.Hash())
						f.forgetHash(hash)
					}
				} else {
					// Fetcher doesn't know about it, add to the return list
					unknown = append(unknown, header)
				}
			}
			headerFilterOutMeter.Mark(int64(len(unknown)))
			select {
			case filter <- &headerFilterTask{headers: unknown, time: task.time}:
			case <-f.quit:
				return
			}
			// Schedule the retrieved headers for body completion
			for _, announce := range incomplete {
				hash := announce.header.Hash()
				if _, ok := f.completing[hash]; ok {
					continue
				}
				//加入到需要获取区块体的队列中
				f.fetched[hash] = append(f.fetched[hash], announce)

				if len(f.fetched) == 1 {
					f.rescheduleComplete(completeTimer)
				}
			}
			//这里插入的是那些区块体为空的区块（没有人转账，发合约，调合约），所以他的区块体就为空
			// Schedule the header-only blocks for import
			for _, block := range complete {
				if announce := f.completing[block.Hash()]; announce != nil {
					f.enqueue(announce.origin, block)
				}
			}

		case filter := <-f.bodyFilter:
			// Block bodies arrived, extract any explicitly requested blocks, return the rest
			var task *bodyFilterTask
			select {
			case task = <-filter:
			case <-f.quit:
				return
			}
			bodyFilterInMeter.Mark(int64(len(task.transactions)))

			blocks := []*types.Block{}
			for i := 0; i < len(task.transactions) && i < len(task.uncles); i++ {
				// Match up a body to any possible completion request
				matched := false

				for hash, announce := range f.completing {
					if f.queued[hash] == nil {
						txnHash := types.DeriveSha(types.Transactions(task.transactions[i]))
						uncleHash := types.CalcUncleHash(task.uncles[i])

						if txnHash == announce.header.TxHash && uncleHash == announce.header.UncleHash && announce.origin == task.peer {
							// Mark the body matched, reassemble if still unknown
							matched = true

							if f.getBlock(hash) == nil {
								block := types.NewBlockWithHeader(announce.header).WithBody(task.transactions[i], task.uncles[i])
								block.ReceivedAt = task.time

								blocks = append(blocks, block)
							} else {
								f.forgetHash(hash)
							}
						}
					}
				}
				if matched {
					task.transactions = append(task.transactions[:i], task.transactions[i+1:]...)
					task.uncles = append(task.uncles[:i], task.uncles[i+1:]...)
					i--
					continue
				}
			}

			bodyFilterOutMeter.Mark(int64(len(task.transactions)))
			select {
			case filter <- task:
			case <-f.quit:
				return
			}
			// Schedule the retrieved blocks for ordered import
			for _, block := range blocks {
				if announce := f.completing[block.Hash()]; announce != nil {
					f.enqueue(announce.origin, block)
				}
			}
		}
	}
}

// rescheduleFetch resets the specified fetch timer to the next announce timeout.
func (f *Fetcher) rescheduleFetch(fetch *time.Timer) {
	// Short circuit if no blocks are announced
	if len(f.announced) == 0 {
		return
	}
	// Otherwise find the earliest expiring announcement
	earliest := time.Now()
	for _, announces := range f.announced {
		if earliest.After(announces[0].time) {
			earliest = announces[0].time
		}
	}
	fetch.Reset(arriveTimeout - time.Since(earliest))
}

// rescheduleComplete resets the specified completion timer to the next fetch timeout.
func (f *Fetcher) rescheduleComplete(complete *time.Timer) {
	// Short circuit if no headers are fetched
	if len(f.fetched) == 0 {
		return
	}
	// Otherwise find the earliest expiring announcement
	earliest := time.Now()
	for _, announces := range f.fetched {
		if earliest.After(announces[0].time) {
			earliest = announces[0].time
		}
	}
	complete.Reset(gatherSlack - time.Since(earliest))
}

// enqueue schedules a new future import operation, if the block to be imported
// has not yet been seen.
func (f *Fetcher) enqueue(peer string, block *types.Block) {
	hash := block.Hash()

	// Ensure the peer isn't DOSing us
	count := f.queues[peer] + 1
	if count > blockLimit {
		log.Debug("Discarded propagated block, exceeded allowance", "peer", peer, "number", block.Number(), "hash", hash, "limit", blockLimit)
		propBroadcastDOSMeter.Mark(1)
		f.forgetHash(hash)
		return
	}
	// Discard any past or too distant blocks
	if dist := int64(block.NumberU64()) - int64(f.chainHeight()); dist < -maxUncleDist || dist > maxQueueDist {
		log.Debug("Discarded propagated block, too far away", "peer", peer, "number", block.Number(), "hash", hash, "distance", dist)
		propBroadcastDropMeter.Mark(1)
		f.forgetHash(hash)
		return
	}
	// Schedule the block for future importing
	if _, ok := f.queued[hash]; !ok {
		op := &inject{
			origin: peer,
			block:  block,
		}
		f.queues[peer] = count
		f.queued[hash] = op
		f.queue.Push(op, -float32(block.NumberU64()))
		if f.queueChangeHook != nil {
			f.queueChangeHook(op.block.Hash(), true)
		}
		log.Debug("Queued propagated block", "peer", peer, "number", block.Number(), "hash", hash, "queued", f.queue.Size())
	}
}

// insert spawns a new goroutine to run a block insertion into the chain. If the
// block's number is at the same height as the current import phase, it updates
// the phase states accordingly.
func (f *Fetcher) insert(peer string, block *types.Block) {
	hash := block.Hash()

	// Run the import on a new thread
	log.Debug("Importing propagated block", "peer", peer, "number", block.Number(), "hash", hash)
	go func() {
		defer func() { f.done <- hash }()

		// If the parent's unknown, abort insertion
		parent := f.getBlock(block.ParentHash())
		if parent == nil {
			log.Debug("Unknown parent of propagated block", "peer", peer, "number", block.Number(), "hash", hash, "parent", block.ParentHash())
			return
		}
		// Quickly validate the header and propagate the block if it passes
		switch err := f.verifyHeader(block.Header()); err {
		case nil:
			// All ok, quickly propagate to our peers
			propBroadcastOutTimer.UpdateSince(block.ReceivedAt)
			go f.broadcastBlock(block, true)

		case consensus.ErrFutureBlock:
			// Weird future block, don't fail, but neither propagate

		default:
			// Something went very wrong, drop the peer
			log.Debug("Propagated block verification failed", "peer", peer, "number", block.Number(), "hash", hash, "err", err)
			f.dropPeer(peer)
			return
		}
		//插入区块
		// Run the actual import and log any issues
		if _, err := f.insertChain(types.Blocks{block}); err != nil {
			log.Debug("Propagated block import failed", "peer", peer, "number", block.Number(), "hash", hash, "err", err)
			return
		}
		// If import succeeded, broadcast the block
		propAnnounceOutTimer.UpdateSince(block.ReceivedAt)

		//广播出去，我有新块了
		go f.broadcastBlock(block, false)

		// Invoke the testing hook if needed
		if f.importedHook != nil {
			f.importedHook(block)
		}
	}()
}

// forgetHash removes all traces of a block announcement from the fetcher's
// internal state.
func (f *Fetcher) forgetHash(hash common.Hash) {
	// Remove all pending announces and decrement DOS counters
	for _, announce := range f.announced[hash] {
		f.announces[announce.origin]--
		if f.announces[announce.origin] == 0 {
			delete(f.announces, announce.origin)
		}
	}
	delete(f.announced, hash)
	if f.announceChangeHook != nil {
		f.announceChangeHook(hash, false)
	}
	// Remove any pending fetches and decrement the DOS counters
	if announce := f.fetching[hash]; announce != nil {
		f.announces[announce.origin]--
		if f.announces[announce.origin] == 0 {
			delete(f.announces, announce.origin)
		}
		delete(f.fetching, hash)
	}

	// Remove any pending completion requests and decrement the DOS counters
	for _, announce := range f.fetched[hash] {
		f.announces[announce.origin]--
		if f.announces[announce.origin] == 0 {
			delete(f.announces, announce.origin)
		}
	}
	delete(f.fetched, hash)

	// Remove any pending completions and decrement the DOS counters
	if announce := f.completing[hash]; announce != nil {
		f.announces[announce.origin]--
		if f.announces[announce.origin] == 0 {
			delete(f.announces, announce.origin)
		}
		delete(f.completing, hash)
	}
}

// forgetBlock removes all traces of a queued block from the fetcher's internal
// state.
func (f *Fetcher) forgetBlock(hash common.Hash) {
	if insert := f.queued[hash]; insert != nil {
		f.queues[insert.origin]--
		if f.queues[insert.origin] == 0 {
			delete(f.queues, insert.origin)
		}
		delete(f.queued, hash)
	}
}

```

