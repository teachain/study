// Copyright 2017 The go-simplechain Authors
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

package core

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/core/rawdb"
	"github.com/simplechain-org/go-simplechain/core/types"
	"github.com/simplechain-org/go-simplechain/ethdb"
	"github.com/simplechain-org/go-simplechain/event"
	"github.com/simplechain-org/go-simplechain/log"
)

// ChainIndexerBackend defines the methods needed to process chain segments in
// the background and write the segment results into the database. These can be
// used to create filter blooms or CHTs.
type ChainIndexerBackend interface {
	// Reset initiates the processing of a new chain segment, potentially terminating
	// any partially completed operations (in case of a reorg).
	Reset(ctx context.Context, section uint64, prevHead common.Hash) error

	// Process crunches through the next header in the chain segment. The caller
	// will ensure a sequential order of headers.
	Process(ctx context.Context, header *types.Header) error

	// Commit finalizes the section metadata and stores it into the database.
	Commit() error
}

// ChainIndexerChain interface is used for connecting the indexer to a blockchain
type ChainIndexerChain interface {
	// CurrentHeader retrieves the latest locally known header.
	CurrentHeader() *types.Header

	// SubscribeChainHeadEvent subscribes to new head header notifications.
	SubscribeChainHeadEvent(ch chan<- ChainHeadEvent) event.Subscription
}

// ChainIndexer does a post-processing job for equally sized sections of the
// canonical chain (like BlooomBits and CHT structures). A ChainIndexer is
// connected to the blockchain through the event system by starting a
// ChainHeadEventLoop in a goroutine.
//
// Further child ChainIndexers can be added which use the output of the parent
// section indexer. These child indexers receive new head notifications only
// after an entire section has been finished or in case of rollbacks that might
// affect already finished sections.
type ChainIndexer struct {

	chainDb  ethdb.Database      // Chain database to index the data from

	//初步推测，indexdb和backend会一对一存在
	indexDb  ethdb.Database      // Prefixed table-view of the db to write index metadata into

	backend  ChainIndexerBackend // Background processor generating the index data content

	children []*ChainIndexer     // Child indexers to cascade chain updates to

	active    uint32          // Flag whether the event loop was started

	update    chan struct{}   // Notification channel that headers should be processed

	quit      chan chan error // Quit channel to tear down running goroutines

	ctx       context.Context

	ctxCancel func()

	//要处理的链片段中的区块数（一个片段中包含多少个区块）
	sectionSize uint64 // Number of blocks in a single chain segment to process

	//处理完成段之前的确认数
	confirmsReq uint64 // Number of confirmations before processing a completed segment

	//成功存入数据库的片段数
	storedSections uint64 // Number of sections successfully indexed into the database

	//已知已完成的片段数
	knownSections  uint64 // Number of sections known to be complete (block wise)

	//级联到子索引器的最后完成节的块号
	cascadedHead   uint64 // Block number of the last completed section cascaded to subindexers

	//检查点覆盖的节数(片段数)（检查点片段数）
	checkpointSections uint64      // Number of sections covered by the checkpoint

	//检查点的片段头哈希(注意到他还是片段中的最后一个哈希的值,但是是sections=checkpointSections-1的)
	checkpointHead     common.Hash // Section head belonging to the checkpoint

	//磁盘限制以防止大量升级占用资源
	throttling time.Duration // Disk throttling to prevent a heavy upgrade from hogging resources

	log  log.Logger

	lock sync.RWMutex
}

// NewChainIndexer creates a new chain indexer to do background processing on
// chain segments of a given size after certain number of confirmations passed.
// The throttling parameter might be used to prevent database thrashing.
func NewChainIndexer(chainDb ethdb.Database, indexDb ethdb.Database, backend ChainIndexerBackend, section, confirm uint64, throttling time.Duration, kind string) *ChainIndexer {
	c := &ChainIndexer{
		chainDb:     chainDb,
		indexDb:     indexDb,
		backend:     backend,
		update:      make(chan struct{}, 1),
		quit:        make(chan chan error),
		sectionSize: section,
		confirmsReq: confirm,
		throttling:  throttling,
		log:         log.New("type", kind),
	}
	// Initialize database dependent fields and start the updater
	c.loadValidSections()

	c.ctx, c.ctxCancel = context.WithCancel(context.Background())

	go c.updateLoop()

	return c
}

// AddCheckpoint adds a checkpoint. Sections are never processed and the chain
// is not expected to be available before this point. The indexer assumes that
// the backend has sufficient information available to process subsequent sections.
//
// Note: knownSections == 0 and storedSections == checkpointSections until
// syncing reaches the checkpoint
func (c *ChainIndexer) AddCheckpoint(section uint64, shead common.Hash) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Short circuit if the given checkpoint is below than local's.
	if c.checkpointSections >= section+1 || section < c.storedSections {
		return
	}
	//
	c.checkpointSections = section + 1

	c.checkpointHead = shead

	c.setSectionHead(section, shead)

	c.setValidSections(section + 1)
}

// Start creates a goroutine to feed chain head events into the indexer for
// cascading background processing. Children do not need to be started, they
// are notified about new events by their parents.
func (c *ChainIndexer) Start(chain ChainIndexerChain) {
	events := make(chan ChainHeadEvent, 10)
	sub := chain.SubscribeChainHeadEvent(events)

	go c.eventLoop(chain.CurrentHeader(), events, sub)
}

// Close tears down all goroutines belonging to the indexer and returns any error
// that might have occurred internally.
func (c *ChainIndexer) Close() error {
	var errs []error

	c.ctxCancel()

	// Tear down the primary update loop
	errc := make(chan error)
	c.quit <- errc
	if err := <-errc; err != nil {
		errs = append(errs, err)
	}
	// If needed, tear down the secondary event loop
	if atomic.LoadUint32(&c.active) != 0 {
		c.quit <- errc
		if err := <-errc; err != nil {
			errs = append(errs, err)
		}
	}
	// Close all children
	for _, child := range c.children {
		if err := child.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	// Return any failures
	switch {
	case len(errs) == 0:
		return nil

	case len(errs) == 1:
		return errs[0]

	default:
		return fmt.Errorf("%v", errs)
	}
}

//处理一个新块的逻辑
// eventLoop is a secondary - optional - event loop of the indexer which is only
// started for the outermost indexer to push chain head events into a processing
// queue.
func (c *ChainIndexer) eventLoop(currentHeader *types.Header, events chan ChainHeadEvent, sub event.Subscription) {
	// Mark the chain indexer as active, requiring an additional teardown
	atomic.StoreUint32(&c.active, 1)

	defer sub.Unsubscribe()

	// Fire the initial new head event to start any outstanding processing
	c.newHead(currentHeader.Number.Uint64(), false)

	var (
		prevHeader = currentHeader
		prevHash   = currentHeader.Hash()
	)
	for {
		select {
		case errc := <-c.quit:
			// Chain indexer terminating, report no failure and abort
			errc <- nil
			return

		case ev, ok := <-events:
			// Received a new event, ensure it's not nil (closing) and update
			if !ok {
				errc := <-c.quit
				errc <- nil
				return
			}
		    // 我们需要注意到，这个事件的特殊性，当快速同步的时候，其实它是批量插入数据的，
		    // 它只会把最新的区块作为事件通知到
		    // 这里我们不能假设它是每个区块都会有通知到。
		    // 我们这里只能说得到了新区块的通知，
			header := ev.Block.Header()

			log.Error("eventLoop","number",header.Number.Uint64())

			if header.ParentHash != prevHash {
				// Reorg to the common ancestor if needed (might not exist in light sync mode, skip reorg then)
				// TODO(karalabe, zsfelfoldi): This seems a bit brittle, can we detect this case explicitly?

				if rawdb.ReadCanonicalHash(c.chainDb, prevHeader.Number.Uint64()) != prevHash {
					if h := rawdb.FindCommonAncestor(c.chainDb, prevHeader, header); h != nil {
						c.newHead(h.Number.Uint64(), true)
					}
				}
			}
			//正常走
			c.newHead(header.Number.Uint64(), false)

			//上一个区块，哈希
			prevHeader, prevHash = header, header.Hash()
		}
	}
}

//主要是变化knownSections
//有新的区块头的时候就调用这个方法
// newHead notifies the indexer about new chain heads and/or reorgs.
func (c *ChainIndexer) newHead(head uint64, reorg bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	//链分叉，重组
	// If a reorg happened, invalidate all sections until that point
	if reorg {
		// Revert the known section number to the reorg point
		known := (head + 1) / c.sectionSize
		fmt.Println("known=",known)
		stored := known
		if known < c.checkpointSections {
			known = 0
		}
		if stored < c.checkpointSections {
			stored = c.checkpointSections
		}

		if known < c.knownSections {

			fmt.Println("333 known=",known,"c.knownSections=",c.knownSections)

			c.knownSections = known


		}
		fmt.Println("2222 known=",known,"c.knownSections=",c.knownSections)
		fmt.Println("2222 stored=",stored,"c.storedSections=",c.storedSections)
		// Revert the stored sections from the database to the reorg point
		if stored < c.storedSections {
			c.setValidSections(stored)
		}
		// Update the new head number to the finalized section end and notify children
		head = known * c.sectionSize

		fmt.Println("c.cascadedHead=",c.cascadedHead,"head=",head)

		if head < c.cascadedHead {
			c.cascadedHead = head
			for _, child := range c.children {
				child.newHead(c.cascadedHead, true)
			}
		}
		return
	}
	//没有发生重组，顺序产生的
	// No reorg, calculate the number of newly known sections and update if high enough
	var sections uint64
	//区块头大于要求的确认数
	if head >= c.confirmsReq {
		//用确认数（也就是当前的区块高度先减去要求的确认数），然后在计算片段数
		sections = (head + 1 - c.confirmsReq) / c.sectionSize
		//fmt.Println("newHead","sections",sections,"checkpointSections",c.checkpointSections)
		//如果片段数小于检查片段数，那么就让它为0，让它不好去处理，再等等
		if sections < c.checkpointSections {
			//fmt.Println("newHead wait","sections",sections ,"c.checkpointSections",c.checkpointSections)
			sections = 0
		}
		//如果片段数已经大于检查片段数，表明你该去处理数据了
		if sections > c.knownSections {
			if c.knownSections < c.checkpointSections {
				// syncing reached the checkpoint, verify section head
				syncedHead := rawdb.ReadCanonicalHash(c.chainDb, c.checkpointSections*c.sectionSize-1)
				if syncedHead != c.checkpointHead {
					c.log.Error("Synced chain does not match checkpoint", "number", c.checkpointSections*c.sectionSize-1, "expected", c.checkpointHead, "synced", syncedHead)
					return
				}
			}
			fmt.Println("newHead update","c.knownSections make to=",sections,"header",head)
			c.knownSections = sections

			select {
			case c.update <- struct{}{}:
			default:
			}
		}
	}
}

//处理一个片段的逻辑
// updateLoop is the main event loop of the indexer which pushes chain segments
// down into the processing backend.
func (c *ChainIndexer) updateLoop() {
	var (
		updating bool
		updated  time.Time
	)

	for {
		select {
		case errc := <-c.quit:
			// Chain indexer terminating, report no failure and abort
			errc <- nil
			return

		case <-c.update:
			// Section headers completed (or rolled back), update the index
			c.lock.Lock()
			if c.knownSections > c.storedSections {
				//fmt.Println("updateLoop","c.knownSections",c.knownSections,"storedSections=",c.storedSections)
				// Periodically print an upgrade log message to the user
				if time.Since(updated) > 8*time.Second {
					if c.knownSections > c.storedSections+1 {
						updating = true
						c.log.Info("Upgrading chain index", "percentage", c.storedSections*100/c.knownSections)
					}
					updated = time.Now()
				}
				// Cache the current section count and head to allow unlocking the mutex
				c.verifyLastHead()
				section := c.storedSections
				var oldHead common.Hash
				if section > 0 {
					//取上一个片段的最后一个区块的哈希
					oldHead = c.SectionHead(section - 1)
				}
				// Process the newly defined section in the background
				c.lock.Unlock()
				//处理当前的分片数据
				newHead, err := c.processSection(section, oldHead)
				if err != nil {
					select {
					case <-c.ctx.Done():
						<-c.quit <- nil
						return
					default:
					}
					fmt.Println("Section processing failed",err)
					c.log.Error("Section processing failed", "error", err)
				}
				c.lock.Lock()

				// If processing succeeded and no reorgs occurred, mark the section completed
				if err == nil && (section == 0 || oldHead == c.SectionHead(section-1)) {

					//处理成功了，就设置分片和它的最后一个区块的哈希
					c.setSectionHead(section, newHead)

					//保存的是处理了多少个片段
					c.setValidSections(section + 1)

					if c.storedSections == c.knownSections && updating {
						updating = false
						c.log.Info("Finished upgrading chain index")
						log.Error("updateLoop Finished upgrading chain index")
					}

					c.cascadedHead = c.storedSections*c.sectionSize - 1

					for _, child := range c.children {
						c.log.Trace("Cascading chain index update", "head", c.cascadedHead)
						child.newHead(c.cascadedHead, false)
					}
				} else {
					// If processing failed, don't retry until further notification
					c.log.Debug("Chain index processing failed", "section", section, "err", err)
					c.verifyLastHead()


					fmt.Println("now updateLoop","c.knownSections make to=",c.storedSections,"err=",err)
					c.knownSections = c.storedSections
				}
			}
			// If there are still further sections to process, reschedule
			if c.knownSections > c.storedSections {
				time.AfterFunc(c.throttling, func() {
					select {
					case c.update <- struct{}{}:
					default:
					}
				})
			}
			c.lock.Unlock()
		}
	}
}

// processSection processes an entire section by calling backend functions while
// ensuring the continuity of the passed headers. Since the chain mutex is not
// held while processing, the continuity can be broken by a long reorg, in which
// case the function returns with an error.
func (c *ChainIndexer) processSection(section uint64, lastHead common.Hash) (common.Hash, error) {
	c.log.Trace("Processing new chain section", "section", section)
    //fmt.Println("Processing new chain section", "section", section)
	// Reset and partial processing

	if err := c.backend.Reset(c.ctx, section, lastHead); err != nil {
		c.setValidSections(0)
		return common.Hash{}, err
	}
    //section是从开始的
	for number := section * c.sectionSize; number < (section+1)*c.sectionSize; number++ {
		//读出经典链的哈希
		hash := rawdb.ReadCanonicalHash(c.chainDb, number)
		if hash == (common.Hash{}) {
			//哈希为空，报错
			return common.Hash{}, fmt.Errorf("canonical block #%d unknown", number)
		}
		//根据区块和哈希，获取区块头
		header := rawdb.ReadHeader(c.chainDb, hash, number)

		//如果区块不存在，报错
		if header == nil {
			return common.Hash{}, fmt.Errorf("block #%d [%x…] not found", number, hash[:4])
		} else if header.ParentHash != lastHead {
			//如果区块的前一个哈希对不上，报错
			return common.Hash{}, fmt.Errorf("chain reorged during section processing")
		}
		//处理区块头
		if err := c.backend.Process(c.ctx, header); err != nil {
			return common.Hash{}, err
		}
		//log.Error("processSection","number",number,"section",section,"sectionSize",c.sectionSize)
		//上一个区块的哈希（链的处理方法，构造如此，所以必须这样子操作）
		lastHead = header.Hash()
	}
	//提交
	if err := c.backend.Commit(); err != nil {
		return common.Hash{}, err
	}
	//fmt.Println("processSection","section",section,"hash",lastHead.String())
	//返回最后一个区块头（严格按顺序处理）
	return lastHead, nil
}

// verifyLastHead compares last stored section head with the corresponding block hash in the
// actual canonical chain and rolls back reorged sections if necessary to ensure that stored
// sections are all valid
func (c *ChainIndexer) verifyLastHead() {
	for c.storedSections > 0 && c.storedSections > c.checkpointSections {
		if c.SectionHead(c.storedSections-1) == rawdb.ReadCanonicalHash(c.chainDb, c.storedSections*c.sectionSize-1) {
			return
		}
		c.setValidSections(c.storedSections - 1)
	}
}

// Sections returns the number of processed sections maintained by the indexer
// and also the information about the last header indexed for potential canonical
// verifications.
func (c *ChainIndexer) Sections() (uint64, uint64, common.Hash) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.verifyLastHead()
	//返回片段数，（片段中最后一个区块）区块高度，（片段中最后一个区块）区块的哈希值
	return c.storedSections, c.storedSections*c.sectionSize - 1, c.SectionHead(c.storedSections - 1)
}

// AddChildIndexer adds a child ChainIndexer that can use the output of this one
func (c *ChainIndexer) AddChildIndexer(indexer *ChainIndexer) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.children = append(c.children, indexer)

	//这里我们要弄清楚stored和known的关系
	// Cascade any pending updates to new children too
	sections := c.storedSections
	if c.knownSections < sections {
		// if a section is "stored" but not "known" then it is a checkpoint without
		// available chain data so we should not cascade it yet
		sections = c.knownSections
	}
	if sections > 0 {
		indexer.newHead(sections*c.sectionSize-1, false)
	}
}

//从数据库中读出存储的片段数
// loadValidSections reads the number of valid sections from the index database
// and caches is into the local state.
func (c *ChainIndexer) loadValidSections() {
	//读count
	//读取已经成功处理的片段
	data, _ := c.indexDb.Get([]byte("count"))
	//uint64,8个字节
	if len(data) == 8 {
		c.storedSections = binary.BigEndian.Uint64(data)
	}
}

//因为它是个数，所以一般是sections+1，因为sections从0开始
//成功处理了多少个片段（严格顺序）
//把合法的片段数存入数据库
// setValidSections writes the number of valid sections to the index database
func (c *ChainIndexer) setValidSections(sections uint64) {
	//写count
	// Set the current number of valid sections in the database
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], sections)
	c.indexDb.Put([]byte("count"), data[:])

	//这是异常的情况才走
	// Remove any reorged sections, caching the valids in the mean time
	for c.storedSections > sections {
		c.storedSections--
		//一个片段一个头哈希值（一个片段的最后一个哈希值）
		c.removeSectionHead(c.storedSections)
	}
	//基本上就是一致的
	c.storedSections = sections // needed if new > old
}


//以下三个方法实现了增删改查（增和改用同一个方法）

// SectionHead retrieves the last block hash of a processed section from the
// index database.
func (c *ChainIndexer) SectionHead(section uint64) common.Hash {
	log.Error("SectionHead","section",section,"checkpoint",c.checkpointSections)
	//读取shead
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], section)
	hash, _ := c.indexDb.Get(append([]byte("shead"), data[:]...))
	if len(hash) == len(common.Hash{}) {
		return common.BytesToHash(hash)
	}
	return common.Hash{}
}

//我们可以这么假设，我们假定一个片段为4096个区块，这样子我们就可以一个片段一个片段地处理（确定了片段大小，我们就一直以这个大小来处理）
//假设我们现在section的值为0，那么我们就将这个片段（0~4095，共4096个块）的最后一个区块（也就是高度（number）为4095）这个块的哈希写入数据库中
//将已处理的片段的最后一个块哈希写入索引数据库
// setSectionHead writes the last block hash of a processed section to the index
// database.
func (c *ChainIndexer) setSectionHead(section uint64, hash common.Hash) {
	log.Info("setSectionHead","section",section,"hash",hash.String())
	//写入shead
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], section)
	c.indexDb.Put(append([]byte("shead"), data[:]...), hash.Bytes())
}

//从数据库中移除片段数对应的哈希值（片段中最后一个块的哈希值）
// removeSectionHead removes the reference to a processed section from the index
// database.
func (c *ChainIndexer) removeSectionHead(section uint64) {
	//删除shead
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], section)
	c.indexDb.Delete(append([]byte("shead"), data[:]...))
}
