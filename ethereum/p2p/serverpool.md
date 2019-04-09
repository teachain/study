```
// Copyright 2016 The go-simplechain Authors
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

// Package les implements the Light Simplechain Subprotocol.
package p2p

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/simplechain-org/go-simplechain/common/mclock"
	"github.com/simplechain-org/go-simplechain/ethdb"
	"github.com/simplechain-org/go-simplechain/log"
	"github.com/simplechain-org/go-simplechain/p2p/discover"
	"github.com/simplechain-org/go-simplechain/rlp"
)

const (
	// After a connection has been ended or timed out, there is a waiting period
	// before it can be selected for connection again.
	// waiting period = base delay * (1 + random(1))
	// base delay = shortRetryDelay for the first shortRetryCnt times after a
	// successful connection, after that longRetryDelay is applied
	shortRetryCnt   = 5
	shortRetryDelay = time.Second * 5
	longRetryDelay  = time.Minute * 10
	// maxNewEntries is the maximum number of newly discovered (never connected) nodes.
	// If the limit is reached, the least recently discovered one is thrown out.
	maxNewEntries = 1000
	// maxKnownEntries is the maximum number of known (already connected) nodes.
	// If the limit is reached, the least recently connected one is thrown out.
	// (not that unlike new entries, known entries are persistent)
	maxKnownEntries = 1000
	// target for simultaneously connected servers
	targetServerCount = 5
	// target for servers selected from the known table
	// (we leave room for trying new ones if there is any)
	targetKnownSelect = 3
	// after dialTimeout, consider the server unavailable and adjust statistics
	dialTimeout = time.Second * 30
	// targetConnTime is the minimum expected connection duration before a server
	// drops a client without any specific reason
	targetConnTime = time.Minute * 10
	// new entry selection weight calculation based on most recent discovery time:
	// unity until discoverExpireStart, then exponential decay with discoverExpireConst
	discoverExpireStart = time.Minute * 20
	discoverExpireConst = time.Minute * 20
	// known entry selection weight is dropped by a factor of exp(-failDropLn) after
	// each unsuccessful connection (restored after a successful one)
	failDropLn = 0.1
	// known node connection success and quality statistics have a long term average
	// and a short term value which is adjusted exponentially with a factor of
	// pstatRecentAdjust with each dial/connection and also returned exponentially
	// to the average with the time constant pstatReturnToMeanTC
	pstatReturnToMeanTC = time.Hour
	// node address selection weight is dropped by a factor of exp(-addrFailDropLn) after
	// each unsuccessful connection (restored after a successful one)
	addrFailDropLn = math.Ln2
	// responseScoreTC and delayScoreTC are exponential decay time constants for
	// calculating selection chances from response times and block delay times
	responseScoreTC = time.Millisecond * 100
	delayScoreTC    = time.Second * 5
	timeoutPow      = 10
	// initStatsWeight is used to initialize previously unknown peers with good
	// statistics to give a chance to prove themselves
	initStatsWeight = 1

	dbKey string = "topicSearch"
)

// connReq represents a request for peer connection.
type connReq struct {
	p      *Peer
	ip     net.IP
	port   uint16
	result chan *PoolEntry
}

// disconnReq represents a request for peer disconnection.
type disconnReq struct {
	entry   *PoolEntry
	stopped bool
	done    chan struct{}
}

// registerReq represents a request for peer registration.
type registerReq struct {
	entry *PoolEntry
	done  chan struct{}
}

// serverPool implements a pool for storing and selecting newly discovered and already
// known light server nodes. It received discovered nodes, stores statistics about
// known nodes and takes care of always having enough good quality servers connected.
type ServerPool struct {
	db ethdb.Database

	dbKey []byte

	server *Server

	quit chan struct{}

	wg *sync.WaitGroup

	connWg sync.WaitGroup

	topic discover.Topic

	discSetPeriod chan time.Duration

	discNodes chan *discover.Node

	discLookups chan bool

	entries map[discover.NodeID]*PoolEntry //一个peer经过connCh的处理之后，再放入entries中

	timeout, enableRetry chan *PoolEntry

	adjustStats chan poolStatAdjust

	connCh chan *connReq //一个peer中的加入，先经过该chan的处理

	disconnCh chan *disconnReq

	registerCh chan *registerReq

	knownQueue, newQueue poolEntryQueue

	knownSelect, newSelect *weightedRandomSelect

	knownSelected, newSelected int

	fastDiscover bool

	running bool

	locker sync.Mutex
}

// newServerPool creates a new serverPool instance
func NewServerPool(db ethdb.Database, quit chan struct{}, wg *sync.WaitGroup) *ServerPool {
	pool := &ServerPool{
		db:           db,
		quit:         quit,
		wg:           wg,
		entries:      make(map[discover.NodeID]*PoolEntry),
		timeout:      make(chan *PoolEntry, 1),
		adjustStats:  make(chan poolStatAdjust, 100),
		enableRetry:  make(chan *PoolEntry, 1),
		connCh:       make(chan *connReq),
		disconnCh:    make(chan *disconnReq),
		registerCh:   make(chan *registerReq),
		knownSelect:  newWeightedRandomSelect(),
		newSelect:    newWeightedRandomSelect(),
		fastDiscover: true,
		dbKey:        []byte(dbKey),
		running:      false,
	}
	//可理解为保存在数据库中的
	pool.knownQueue = newPoolEntryQueue(maxKnownEntries, pool.removeEntry)
	pool.newQueue = newPoolEntryQueue(maxNewEntries, pool.removeEntry)
	return pool
}

func (pool *ServerPool) Start(server *Server, registerTopic discover.Topic, searchTopics []discover.Topic) {
	pool.locker.Lock()
	defer pool.locker.Unlock()
	if pool.running {
		return
	}
	pool.running = true
	pool.server = server
	pool.wg.Add(1)
	pool.loadNodes()

	if pool.server.Network != nil {
		pool.discSetPeriod = make(chan time.Duration, 1)
		pool.discNodes = make(chan *discover.Node, 100)
		pool.discLookups = make(chan bool, 100)
		//register only one
		go pool.server.Network.RegisterTopic(registerTopic, pool.quit)
        //search some
		for i := range searchTopics {
			go pool.server.Network.SearchTopic(searchTopics[i], pool.discSetPeriod, pool.discNodes, pool.discLookups)
		}
	}
	pool.checkDial()
	go pool.eventLoop()
}

//这里表示的是经过了二阶段握手以后的连接
// connect should be called upon any incoming connection. If the connection has been
// dialed by the server pool recently, the appropriate pool entry is returned.
// Otherwise, the connection should be rejected.
// Note that whenever a connection has been accepted and a pool entry has been returned,
// disconnect should also always be called.
func (pool *ServerPool) Connect(p *Peer, ip net.IP, port uint16) *PoolEntry {
	log.Debug("Connect new entry", "enode", p.rw.id)
	//连接请求
	req := &connReq{p: p, ip: ip, port: port, result: make(chan *PoolEntry, 1)}
	select {
	case pool.connCh <- req: //交给池子的连接通道处理
	case <-pool.quit:
		return nil
	}
	//等connCh那边处理完毕以后，会将结果返回放入到result这chan中
	//然后我们就等着result中得到元素然后返回
	return <-req.result
}

// registered should be called after a successful handshake
// 看英文注释，也就是成功握手之后（第三阶段的握手），一定要调用该方法
func (pool *ServerPool) Registered(entry *PoolEntry) {
	log.Debug("Registered new entry", "enode", entry.id)
	req := &registerReq{entry: entry, done: make(chan struct{})}
	select {
	case pool.registerCh <- req:
	case <-pool.quit:
		return
	}
	<-req.done
}

// disconnect should be called when ending a connection. Service quality statistics
// can be updated optionally (not updated if no registration happened, in this case
// only connection statistics are updated, just like in case of timeout)
func (pool *ServerPool) Disconnect(entry *PoolEntry) {
	stopped := false
	select {
	case <-pool.quit:
		stopped = true
	default:
	}

	log.Debug("Disconnected old entry", "enode", entry.id)

	req := &disconnReq{entry: entry, stopped: stopped, done: make(chan struct{})}

	// Block until disconnection request is served.
	pool.disconnCh <- req //交给disconnCh处理

	<-req.done //等待处理结果
}

const (
	pseBlockDelay = iota
	pseResponseTime
	pseResponseTimeout
)

// poolStatAdjust records are sent to adjust peer block delay/response time statistics
type poolStatAdjust struct {
	adjustType int
	entry      *PoolEntry
	time       time.Duration
}

//注意这个方法在fetcher中使用
// adjustBlockDelay adjusts the block announce delay statistics of a node
func (pool *ServerPool) adjustBlockDelay(entry *PoolEntry, time time.Duration) {
	if entry == nil {
		return
	}
	pool.adjustStats <- poolStatAdjust{pseBlockDelay, entry, time}
}

//注意这个方法在fetcher中使用
// adjustResponseTime adjusts the request response time statistics of a node
func (pool *ServerPool) adjustResponseTime(entry *PoolEntry, time time.Duration, timeout bool) {
	if entry == nil {
		return
	}
	if timeout {
		pool.adjustStats <- poolStatAdjust{pseResponseTimeout, entry, time}
	} else {
		pool.adjustStats <- poolStatAdjust{pseResponseTime, entry, time}
	}
}

// eventLoop handles pool events and mutex locking for all internal functions
func (pool *ServerPool) eventLoop() {
	lookupCnt := 0
	var convTime mclock.AbsTime
	if pool.discSetPeriod != nil {
		pool.discSetPeriod <- time.Millisecond * 100
	}

	// disconnect updates service quality statistics depending on the connection time
	// and disconnection initiator.
	disconnect := func(req *disconnReq, stopped bool) {
		// Handle peer disconnection requests.
		entry := req.entry
		if entry.state == psRegistered {
			connAdjust := float64(mclock.Now()-entry.regTime) / float64(targetConnTime)
			if connAdjust > 1 {
				connAdjust = 1
			}
			if stopped {
				// disconnect requested by ourselves.
				entry.connectStats.add(1, connAdjust)
			} else {
				// disconnect requested by server side.
				entry.connectStats.add(connAdjust, 1)
			}
		}

		//连接断开后，它的状态就设置为psNotConnected了
		entry.state = psNotConnected

		if entry.knownSelected {
			pool.knownSelected--
		} else {
			pool.newSelected--
		}
		pool.setRetryDial(entry)
		pool.connWg.Done()
		close(req.done)
	}

	for {
		select {
		case entry := <-pool.timeout:
			if !entry.removed {
				pool.checkDialTimeout(entry)
			}

		case entry := <-pool.enableRetry:
			if !entry.removed {
				entry.delayedRetry = false
				pool.updateCheckDial(entry)
			}

		case adj := <-pool.adjustStats:
			switch adj.adjustType {
			case pseBlockDelay:
				adj.entry.delayStats.add(float64(adj.time), 1)
			case pseResponseTime:
				adj.entry.responseStats.add(float64(adj.time), 1)
				adj.entry.timeoutStats.add(0, 1)
			case pseResponseTimeout:
				adj.entry.timeoutStats.add(1, 1)
			}

		case node := <-pool.discNodes: //这里是得到了主题节点以后的处理
			entry := pool.findOrNewNode(discover.NodeID(node.ID), node.IP, node.TCP)
			pool.updateCheckDial(entry)

		case conv := <-pool.discLookups: //这是查找的策略
			if conv { //这表示找到了主题节点
				if lookupCnt == 0 {
					convTime = mclock.Now()
				}
				lookupCnt++
				//前面采用快速查找
				//随着时间的推移，演进，我们要修改策略
				//如果找到到的数量已经到达50个了，或者已经超时一个分钟了，我们让它的等待时间长一些
				//就是主题查找的频率降低一些
				if pool.fastDiscover && (lookupCnt == 50 || time.Duration(mclock.Now()-convTime) > time.Minute) {
					pool.fastDiscover = false
					if pool.discSetPeriod != nil {
						pool.discSetPeriod <- time.Minute
					}
				}
			}
		case req := <-pool.connCh: //从连接通道（connCh）中取出连接请求
			// Handle peer connection requests.
			//从池中试着取，要特别注意entry它是一个指针
			entry := pool.entries[req.p.ID()]

			//如果它不在池中
			if entry == nil {
				//当走这一步时，就不可能走下一步（entry.state == psConnected）
				//当它new出来以后，就会把它加入entries中
				entry = pool.findOrNewNode(req.p.ID(), req.ip, req.port)
			}
			//如果这个实体的状态是已经连接，或者是已经注册，那么就直接返回
			if entry.state == psConnected || entry.state == psRegistered {
				req.result <- nil
				continue
			}
			//因为上面已经判断了它的状态，所以这里肯定就是一个新连接了的设置了
			//当然不一定是一个新节点
			//也就是说如果它已经
			pool.connWg.Add(1)
			entry.peer = req.p
			entry.state = psConnected
			//因为到这里已经是使用tcp连接了，所以这里的port就对应tcp的端口
			//ip就自不用说了。
			addr := &poolEntryAddress{
				ip:       req.ip,
				port:     req.port,
				lastSeen: mclock.Now(), //记录下最近见到的时间
			}
			entry.lastConnected = addr //最近连接使用的地址，因为我们是使用nodeID来识别一个节点的
			entry.addr = make(map[string]*poolEntryAddress)
			//strKey就是使用了addr的ip和端口
			entry.addr[addr.strKey()] = addr

			entry.addrSelect = *newWeightedRandomSelect()

			entry.addrSelect.update(addr)

			req.result <- entry

		case req := <-pool.registerCh:
			// Handle peer registration requests.
			entry := req.entry
			entry.state = psRegistered
			entry.regTime = mclock.Now()
			if !entry.known {
				pool.newQueue.remove(entry)
				entry.known = true
			}
			pool.knownQueue.setLatest(entry)
			entry.shortRetry = shortRetryCnt
			close(req.done)

		case req := <-pool.disconnCh:
			// Handle peer disconnection requests.
			//连接断开
			disconnect(req, req.stopped)

		case <-pool.quit:
			if pool.discSetPeriod != nil {
				close(pool.discSetPeriod)
			}

			// Spawn a goroutine to close the disconnCh after all connections are disconnected.
			go func() {
				pool.connWg.Wait()
				close(pool.disconnCh)
			}()

			// Handle all remaining disconnection requests before exit.
			for req := range pool.disconnCh {
				disconnect(req, true)
			}
			pool.saveNodes()
			pool.wg.Done()
			return
		}
	}
}

func (pool *ServerPool) findOrNewNode(id discover.NodeID, ip net.IP, port uint16) *PoolEntry {

	now := mclock.Now()

	entry := pool.entries[id]

	if entry == nil {
		log.Debug("Discovered new entry", "id", id)
		entry = &PoolEntry{
			id:         id,
			addr:       make(map[string]*poolEntryAddress),
			addrSelect: *newWeightedRandomSelect(),
			shortRetry: shortRetryCnt,
		}
		pool.entries[id] = entry
		// initialize previously unknown peers with good statistics to give a chance to prove themselves
		entry.connectStats.add(1, initStatsWeight)
		entry.delayStats.add(0, initStatsWeight)
		entry.responseStats.add(0, initStatsWeight)
		entry.timeoutStats.add(0, initStatsWeight)
	}
	entry.lastDiscovered = now
	addr := &poolEntryAddress{
		ip:   ip,
		port: port,
	}
	if a, ok := entry.addr[addr.strKey()]; ok {
		addr = a
	} else {
		entry.addr[addr.strKey()] = addr
	}
	addr.lastSeen = now
	entry.addrSelect.update(addr)
	if !entry.known {
		pool.newQueue.setLatest(entry)
	}
	return entry
}

//从数据库中取出已知的队列的节点，恢复到已知队列中
// loadNodes loads known nodes and their statistics from the database
func (pool *ServerPool) loadNodes() {
	enc, err := pool.db.Get(pool.dbKey)
	if err != nil {
		return
	}
	var list []*PoolEntry
	err = rlp.DecodeBytes(enc, &list)
	if err != nil {
		log.Debug("Failed to decode node list", "err", err)
		return
	}
	for _, e := range list {
		log.Debug("Loaded server stats", "id", e.id, "fails", e.lastConnected.fails,
			"conn", fmt.Sprintf("%v/%v", e.connectStats.avg, e.connectStats.weight),
			"delay", fmt.Sprintf("%v/%v", time.Duration(e.delayStats.avg), e.delayStats.weight),
			"response", fmt.Sprintf("%v/%v", time.Duration(e.responseStats.avg), e.responseStats.weight),
			"timeout", fmt.Sprintf("%v/%v", e.timeoutStats.avg, e.timeoutStats.weight))

		//加入池中
		pool.entries[e.id] = e

		//这将恢复它原来的顺序
		pool.knownQueue.setLatest(e)

		pool.knownSelect.update((*knownEntry)(e))
	}
}

//把已知队列中的节点保存到数据库中
//方便下次启动时使用
// saveNodes saves known nodes and their statistics into the database. Nodes are
// ordered from least to most recently connected.
func (pool *ServerPool) saveNodes() {
	//已知队列
	list := make([]*PoolEntry, len(pool.knownQueue.queue))
	//队列，它可是有顺序的
	for i := range list {
		list[i] = pool.knownQueue.fetchOldest()
	}
	//将切片进行编码
	enc, err := rlp.EncodeToBytes(list)
	if err == nil {
		//将数据保存到数据库
		pool.db.Put(pool.dbKey, enc)
	}
}

// removeEntry removes a pool entry when the entry count limit is reached.
// Note that it is called by the new/known queues from which the entry has already
// been removed so removing it from the queues is not necessary.
func (pool *ServerPool) removeEntry(entry *PoolEntry) {
	pool.newSelect.remove((*discoveredEntry)(entry))
	pool.knownSelect.remove((*knownEntry)(entry))
	entry.removed = true
	delete(pool.entries, entry.id)
}

// setRetryDial starts the timer which will enable dialing a certain node again
func (pool *ServerPool) setRetryDial(entry *PoolEntry) {
	delay := longRetryDelay
	if entry.shortRetry > 0 {
		entry.shortRetry--
		delay = shortRetryDelay
	}
	delay += time.Duration(rand.Int63n(int64(delay) + 1))
	entry.delayedRetry = true
	go func() {
		select {
		case <-pool.quit:
		case <-time.After(delay):
			select {
			case <-pool.quit:
			case pool.enableRetry <- entry:
			}
		}
	}()
}

// updateCheckDial is called when an entry can potentially be dialed again. It updates
// its selection weights and checks if new dials can/should be made.
//当一个节点有可能被重新建立连接的收，就应该调用该接口
//它会更新选择权重，并且检查新的连接是否应该建立
func (pool *ServerPool) updateCheckDial(entry *PoolEntry) {
	pool.newSelect.update((*discoveredEntry)(entry))
	pool.knownSelect.update((*knownEntry)(entry))
	pool.checkDial()
}

// checkDial checks if new dials can/should be made. It tries to select servers both
// based on good statistics and recent discovery.
func (pool *ServerPool) checkDial() {
	fillWithKnownSelects := !pool.fastDiscover
	for pool.knownSelected < targetKnownSelect {
		entry := pool.knownSelect.choose()
		if entry == nil {
			fillWithKnownSelects = false
			break
		}
		pool.dial((*PoolEntry)(entry.(*knownEntry)), true)
	}
	for pool.knownSelected+pool.newSelected < targetServerCount {
		entry := pool.newSelect.choose()
		if entry == nil {
			break
		}
		pool.dial((*PoolEntry)(entry.(*discoveredEntry)), false)
	}
	if fillWithKnownSelects {
		// no more newly discovered nodes to select and since fast discover period
		// is over, we probably won't find more in the near future so select more
		// known entries if possible
		for pool.knownSelected < targetServerCount {
			entry := pool.knownSelect.choose()
			if entry == nil {
				break
			}
			pool.dial((*PoolEntry)(entry.(*knownEntry)), true)
		}
	}
}

// dial initiates a new connection
func (pool *ServerPool) dial(entry *PoolEntry, knownSelected bool) {
	if pool.server == nil || entry.state != psNotConnected {
		return
	}
	entry.state = psDialed
	entry.knownSelected = knownSelected
	if knownSelected {
		pool.knownSelected++
	} else {
		pool.newSelected++
	}
	addr := entry.addrSelect.choose().(*poolEntryAddress)
	log.Debug("Dialing new peer", "lesaddr", entry.id.String()+"@"+addr.strKey(), "set", len(entry.addr), "known", knownSelected)
	entry.dialed = addr
	go func() {
		//这里是让它加入列表中，让它被重新选择，然后发起dial操作
		pool.server.AddPeer(discover.NewNode(entry.id, addr.ip, addr.port, addr.port))
		select {
		case <-pool.quit:
		case <-time.After(dialTimeout):
			select {
			case <-pool.quit:
			case pool.timeout <- entry:
			}
		}
	}()
}

// checkDialTimeout checks if the node is still in dialed state and if so, resets it
// and adjusts connection statistics accordingly.
func (pool *ServerPool) checkDialTimeout(entry *PoolEntry) {
	if entry.state != psDialed {
		return
	}
	log.Debug("Dial timeout", "lesaddr", entry.id.String()+"@"+entry.dialed.strKey())
	entry.state = psNotConnected
	if entry.knownSelected {
		pool.knownSelected--
	} else {
		pool.newSelected--
	}
	entry.connectStats.add(0, 1)
	entry.dialed.fails++
	pool.setRetryDial(entry)
}

const (
	psNotConnected = iota
	psDialed
	psConnected
	psRegistered
)

// poolEntry represents a server node and stores its current state and statistics.
type PoolEntry struct {
	peer                  *Peer
	id                    discover.NodeID
	addr                  map[string]*poolEntryAddress
	lastConnected, dialed *poolEntryAddress
	addrSelect            weightedRandomSelect

	lastDiscovered              mclock.AbsTime
	known, knownSelected        bool
	connectStats, delayStats    poolStats
	responseStats, timeoutStats poolStats
	state                       int
	regTime                     mclock.AbsTime

	//在队列中的索引
	queueIdx int

	removed bool

	delayedRetry bool
	shortRetry   int
}

func (e *PoolEntry) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{e.id, e.lastConnected.ip, e.lastConnected.port, e.lastConnected.fails, &e.connectStats, &e.delayStats, &e.responseStats, &e.timeoutStats})
}

func (e *PoolEntry) DecodeRLP(s *rlp.Stream) error {
	var entry struct {
		ID                         discover.NodeID
		IP                         net.IP
		Port                       uint16
		Fails                      uint
		CStat, DStat, RStat, TStat poolStats
	}
	if err := s.Decode(&entry); err != nil {
		return err
	}
	addr := &poolEntryAddress{ip: entry.IP, port: entry.Port, fails: entry.Fails, lastSeen: mclock.Now()}
	e.id = entry.ID
	e.addr = make(map[string]*poolEntryAddress)
	e.addr[addr.strKey()] = addr
	e.addrSelect = *newWeightedRandomSelect()
	e.addrSelect.update(addr)
	e.lastConnected = addr
	e.connectStats = entry.CStat
	e.delayStats = entry.DStat
	e.responseStats = entry.RStat
	e.timeoutStats = entry.TStat
	e.shortRetry = shortRetryCnt
	e.known = true
	return nil
}

// discoveredEntry implements wrsItem
type discoveredEntry PoolEntry

// Weight calculates random selection weight for newly discovered entries
func (e *discoveredEntry) Weight() int64 {
	if e.state != psNotConnected || e.delayedRetry {
		return 0
	}
	t := time.Duration(mclock.Now() - e.lastDiscovered)
	if t <= discoverExpireStart {
		return 1000000000
	}
	return int64(1000000000 * math.Exp(-float64(t-discoverExpireStart)/float64(discoverExpireConst)))
}

// knownEntry implements wrsItem
type knownEntry PoolEntry

// Weight calculates random selection weight for known entries
func (e *knownEntry) Weight() int64 {
	if e.state != psNotConnected || !e.known || e.delayedRetry {
		return 0
	}
	return int64(1000000000 * e.connectStats.recentAvg() * math.Exp(-float64(e.lastConnected.fails)*failDropLn-e.responseStats.recentAvg()/float64(responseScoreTC)-e.delayStats.recentAvg()/float64(delayScoreTC)) * math.Pow(1-e.timeoutStats.recentAvg(), timeoutPow))
}

// poolEntryAddress is a separate object because currently it is necessary to remember
// multiple potential network addresses for a pool entry. This will be removed after
// the final implementation of v5 discovery which will retrieve signed and serial
// numbered advertisements, making it clear which IP/port is the latest one.
type poolEntryAddress struct {
	ip       net.IP
	port     uint16
	lastSeen mclock.AbsTime // last time it was discovered, connected or loaded from db
	fails    uint           // connection failures since last successful connection (persistent)
}

func (a *poolEntryAddress) Weight() int64 {
	t := time.Duration(mclock.Now() - a.lastSeen)
	return int64(1000000*math.Exp(-float64(t)/float64(discoverExpireConst)-float64(a.fails)*addrFailDropLn)) + 1
}

func (a *poolEntryAddress) strKey() string {
	return a.ip.String() + ":" + strconv.Itoa(int(a.port))
}

// poolStats implement statistics for a certain quantity with a long term average
// and a short term value which is adjusted exponentially with a factor of
// pstatRecentAdjust with each update and also returned exponentially to the
// average with the time constant pstatReturnToMeanTC
type poolStats struct {
	sum, weight, avg, recent float64
	lastRecalc               mclock.AbsTime
}

// init initializes stats with a long term sum/update count pair retrieved from the database
func (s *poolStats) init(sum, weight float64) {
	s.sum = sum
	s.weight = weight
	var avg float64
	if weight > 0 {
		avg = s.sum / weight
	}
	s.avg = avg
	s.recent = avg
	s.lastRecalc = mclock.Now()
}

// recalc recalculates recent value return-to-mean and long term average
func (s *poolStats) recalc() {
	now := mclock.Now()
	s.recent = s.avg + (s.recent-s.avg)*math.Exp(-float64(now-s.lastRecalc)/float64(pstatReturnToMeanTC))
	if s.sum == 0 {
		s.avg = 0
	} else {
		if s.sum > s.weight*1e30 {
			s.avg = 1e30
		} else {
			s.avg = s.sum / s.weight
		}
	}
	s.lastRecalc = now
}

// add updates the stats with a new value
func (s *poolStats) add(value, weight float64) {
	s.weight += weight
	s.sum += value * weight
	s.recalc()
}

// recentAvg returns the short-term adjusted average
func (s *poolStats) recentAvg() float64 {
	s.recalc()
	return s.recent
}

func (s *poolStats) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{math.Float64bits(s.sum), math.Float64bits(s.weight)})
}

func (s *poolStats) DecodeRLP(st *rlp.Stream) error {
	var stats struct {
		SumUint, WeightUint uint64
	}
	if err := st.Decode(&stats); err != nil {
		return err
	}
	s.init(math.Float64frombits(stats.SumUint), math.Float64frombits(stats.WeightUint))
	return nil
}

// poolEntryQueue keeps track of its least recently accessed entries and removes
// them when the number of entries reaches the limit
type poolEntryQueue struct {

	//用map来模拟队列
	queue map[int]*PoolEntry // known nodes indexed by their latest lastConnCnt value

	newPtr int //队列头指针

	oldPtr int //队列尾指针

	maxCnt int //最大数量

	removeFromPool func(*PoolEntry) //从池中移除的函数，目前看，只有队列满了才会删除
}

//创建一个池子队列
// newPoolEntryQueue returns a new poolEntryQueue
func newPoolEntryQueue(maxCnt int, removeFromPool func(*PoolEntry)) poolEntryQueue {
	return poolEntryQueue{queue: make(map[int]*PoolEntry), maxCnt: maxCnt, removeFromPool: removeFromPool}
}

//获取最少使用的实体
// fetchOldest returns and removes the least recently accessed entry
func (q *poolEntryQueue) fetchOldest() *PoolEntry {
	if len(q.queue) == 0 {
		return nil
	}
	for {
		//因为有可能中途删除了一些的，比如调换了位置
		if e := q.queue[q.oldPtr]; e != nil {
			delete(q.queue, q.oldPtr)
			q.oldPtr++
			return e
		}
		q.oldPtr++
	}
}

//从队列中移除
// remove removes an entry from the queue
func (q *poolEntryQueue) remove(entry *PoolEntry) {
	if q.queue[entry.queueIdx] == entry {
		delete(q.queue, entry.queueIdx)
	}
}

// setLatest adds or updates a recently accessed entry. It also checks if an old entry
// needs to be removed and removes it from the parent pool too with a callback function.
func (q *poolEntryQueue) setLatest(entry *PoolEntry) {
	//也就是说如果它原来在队列中，我们要先把它从队列中移除（因为要修改它所在的位置）
	//我们先将它移除
	if q.queue[entry.queueIdx] == entry {
		delete(q.queue, entry.queueIdx)
	} else {
		//判断一下是否还有空位，如果没有，删除队尾
		if len(q.queue) == q.maxCnt {
			e := q.fetchOldest()
			q.remove(e)
			q.removeFromPool(e)
		}
	}
	//设置为最近使用的，也就是说设置为队列的头部
	entry.queueIdx = q.newPtr
	//放入队列中
	q.queue[entry.queueIdx] = entry
	//头指针往前移
	q.newPtr++
}

```

