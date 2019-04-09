

从以下的方法我们知道，节点的数据我们是首先从数据库文件中读取，然后加上种子节点(bootnodes),然后把它加入k桶（可以认为它是一个运行时的一个仓库）。

```
func (tab *Table) loadSeedNodes() {
	//从数据库中查询节点信息
	seeds := tab.db.querySeeds(seedCount, seedMaxAge)
	//从数据库中查询出来的节点和种子节点合在一起
	seeds = append(seeds, tab.nursery...)
	//把它添加到k桶中
	for i := range seeds {
		seed := seeds[i]
		age := log.Lazy{Fn: func() interface{} { return time.Since(tab.db.lastPongReceived(seed.ID)) }}
		log.Debug("Found seed node in database", "id", seed.ID, "addr", seed.addr(), "age", age)
		tab.add(seed)
	}
}

//添加入k桶的逻辑
func (tab *Table) add(n *Node) {
	tab.mutex.Lock()
	defer tab.mutex.Unlock()
    //根据哈希决定它应该放在哪个桶
	b := tab.bucket(n.sha)
	//所以这里的逻辑就是节点能不能放桶里，不能放桶里就放替换列表里
	if !tab.bumpOrAdd(b, n) {
		// Node is not in table. Add it to the replacement list.
		//把节点加入到替换列表里
		tab.addReplacement(b, n)
	}
}

//就一个逻辑，n在k桶中吗
//在，那么把它的位置挪到第一的位置，告诉你它在桶中
//不在，就啥也不干，就告诉你它不在桶中
func (b *bucket) bump(n *Node) bool {
	for i := range b.entries {
		if b.entries[i].ID == n.ID {
			// move it to the front
			copy(b.entries[1:], b.entries[:i])
			b.entries[0] = n
			return true
		}
	}
	return false
}

func (tab *Table) bumpOrAdd(b *bucket, n *Node) bool {
	//在桶中，毫无疑问了，返回true就好了
	if b.bump(n) {
		return true
	}
	//检测一下桶是否还有位置，或者检测这个ip能不能添加进来
	//这里还要特别注意tab.addIP(b, n.IP)，也就是说只要添加进了桶里，肯定这步是true的
	//就是肯定加入了ip列表
	if len(b.entries) >= bucketSize || !tab.addIP(b, n.IP) {
		return false
	}
	//有位置，那么就直接放在第一个位置上
	b.entries, _ = pushNode(b.entries, n, bucketSize)
	//如果这个节点在replacements中，那么就从replacements删除它
	b.replacements = deleteNode(b.replacements, n)
	//设置它添加进来的时间
	n.addedAt = time.Now()
	//如果有回调方法，那么调用回调方法
	if tab.nodeAddedHook != nil {
		tab.nodeAddedHook(n)
	}
	//返回放入桶里了
	return true
}

//随机一个非空的桶，取桶里最后一个元素(可看做一个队列)
// nodeToRevalidate returns the last node in a random, non-empty bucket.
func (tab *Table) nodeToRevalidate() (n *Node, bi int) {
	tab.mutex.Lock()
	defer tab.mutex.Unlock()

	for _, bi = range tab.rand.Perm(len(tab.buckets)) {
		b := tab.buckets[bi]
		if len(b.entries) > 0 {
			last := b.entries[len(b.entries)-1]
			return last, bi
		}
	}
	return nil, 0
}

// doRevalidate checks that the last node in a random bucket is still live
// and replaces or deletes the node if it isn't.
func (tab *Table) doRevalidate(done chan<- struct{}) {
	
	//告诉监听done的goroutine,我这边验证已经完毕
	defer func() { done <- struct{}{} }()

	last, bi := tab.nodeToRevalidate()
	if last == nil {
		// No non-empty bucket found.
		return
	}

	// Ping the selected node and wait for a pong.
	err := tab.net.ping(last.ID, last.addr())

	tab.mutex.Lock()
	defer tab.mutex.Unlock()
	b := tab.buckets[bi]
	//这里进行ping pong之后没有问题(意思是活得好好的)，那么就把节点放入桶里
	//当然肯定是放在第一个位置。这里的bump并没有检查它的返回值，因为不需要
	//因为last就是从桶里拿出来的，肯定在桶里
	if err == nil {
		// The node responded, move it to the front.
		log.Trace("Revalidated node", "b", bi, "id", last.ID)
		b.bump(last)
		return
	}
	// No reply received, pick a replacement or delete the node if there aren't
	// any replacements.
	if r := tab.replace(b, last); r != nil {
		log.Trace("Replaced dead node", "b", bi, "id", last.ID, "ip", last.IP, "r", r.ID, "rip", r.IP)
	} else {
		log.Trace("Removed dead node", "b", bi, "id", last.ID, "ip", last.IP)
	}
}
```



当从桶里取出来的节点进行ping之后得到err,那么

1、检查这个桶是否为空或者最后一个元素并不是该元素，直接返回nil

2、如果此时替换列表还是空的，我们直接从桶里将last节点删除，返回返回nil

3、如果都不是上面的情况，那么就从替换列表里随机取一个节点，把该节点放入桶的最后一个位置（有利于探测是否存活），然后从替换列表中移除。

```
// replace removes n from the replacement list and replaces 'last' with it if it is the
// last entry in the bucket. If 'last' isn't the last entry, it has either been replaced
// with someone else or became active.
func (tab *Table) replace(b *bucket, last *Node) *Node {
	if len(b.entries) == 0 || b.entries[len(b.entries)-1].ID != last.ID {
		// Entry has moved, don't replace it.
		return nil
	}
	// Still the last entry.
	if len(b.replacements) == 0 {
		tab.deleteInBucket(b, last)
		return nil
	}
	r := b.replacements[tab.rand.Intn(len(b.replacements))]
	b.replacements = deleteNode(b.replacements, r)
	b.entries[len(b.entries)-1] = r
	tab.removeIP(b, last.IP)
	return r
}
```



将ip加入列表中，一个tab中的ips,一个是bucket中的ips

```
func (tab *Table) addIP(b *bucket, ip net.IP) bool {
	if netutil.IsLAN(ip) {
		return true
	}
	if !tab.ips.Add(ip) {
		log.Debug("IP exceeds table limit", "ip", ip)
		return false
	}
	if !b.ips.Add(ip) {
		log.Debug("IP exceeds bucket limit", "ip", ip)
		//这里保证了一个ip在tab中就肯定在bucket中
		tab.ips.Remove(ip)
		return false
	}
	return true
}

func (tab *Table) removeIP(b *bucket, ip net.IP) {
	if netutil.IsLAN(ip) {
		return
	}
	tab.ips.Remove(ip)
	b.ips.Remove(ip)
}
```



```
// copyLiveNodes adds nodes from the table to the database if they have been in the table
// longer then minTableTime.
func (tab *Table) copyLiveNodes() {
	tab.mutex.Lock()
	defer tab.mutex.Unlock()
	
	now := time.Now()
	for _, b := range tab.buckets {
		for _, n := range b.entries {
			//如果节点在桶中的时间已经超过5分钟了，那么就将该节点数据持久化到数据库中
			if now.Sub(n.addedAt) >= seedMinTableTime {
				tab.db.updateNode(n)
			}
		}
	}
}
```







//这里我们看到，持久化到数据库中的节点也是有时间限制的，如果这个节点

//的pong时间是一天以前的话，就从数据库中删除。其实就是保持一天之内，收到pong的节点。

```
// expireNodes iterates over the database and deletes all nodes that have not
// been seen (i.e. received a pong from) for some allotted time.
func (db *nodeDB) expireNodes() error {
	threshold := time.Now().Add(-nodeDBNodeExpiration)

	// Find discovered nodes that are older than the allowance
	it := db.lvl.NewIterator(nil, nil)
	defer it.Release()

	for it.Next() {
		// Skip the item if not a discovery node
		id, field := splitKey(it.Key())
		if field != nodeDBDiscoverRoot {
			continue
		}
		// Skip the node if not expired yet (and not self)
		if !bytes.Equal(id[:], db.self[:]) {
			if seen := db.lastPongReceived(id); seen.After(threshold) {
				continue
			}
		}
		// Otherwise delete all associated information
		db.deleteNode(id)
	}
	return nil
}
```



table的loop循环

```
// loop schedules refresh, revalidate runs and coordinates shutdown.
func (tab *Table) loop() {
	var (
		revalidate     = time.NewTimer(tab.nextRevalidateTime())
		refresh        = time.NewTicker(refreshInterval) //半个小时
		copyNodes      = time.NewTicker(copyNodesInterval)//30秒
		revalidateDone = make(chan struct{})
		refreshDone    = make(chan struct{})           // where doRefresh reports completion
		waiting        = []chan struct{}{tab.initDone} // holds waiting callers while doRefresh runs
	)
	defer refresh.Stop()
	defer revalidate.Stop()
	defer copyNodes.Stop()

	// Start initial refresh.
	go tab.doRefresh(refreshDone)

loop:
	for {
		select {
		case <-refresh.C:
		    //半个小时过去了，设置一下随机种子
			tab.seedRand()
			//if refreshDone == nil 这样保证了doRefresh的单次调用，也即是说不会同时有两个goroutine
			//在跑tab.doRefresh(refreshDone)
			if refreshDone == nil {
				refreshDone = make(chan struct{})
				go tab.doRefresh(refreshDone)
			}
		case req := <-tab.refreshReq:
			waiting = append(waiting, req)
			//if refreshDone == nil 这样保证了doRefresh的单次调用，也即是说不会同时有两个goroutine
			//在跑tab.doRefresh(refreshDone)
			if refreshDone == nil {
				refreshDone = make(chan struct{})
				go tab.doRefresh(refreshDone)
			}
		case <-refreshDone:
		    刷新完毕了，就把waiting和refreshDone都置于初始状态
			for _, ch := range waiting {
				close(ch)
			}
			//这样不管是定时到点了还是刷新请求都可以继续了
			waiting, refreshDone = nil, nil
		case <-revalidate.C:
			//验证节点存活，10秒之内随机，它就完成了节点的存活刷新，死节点的移除，用替代列表中的节点填充k桶
			go tab.doRevalidate(revalidateDone)
		case <-revalidateDone:
			//每次如果doRevalidate已经验证完，那么就设置一下定时器
			revalidate.Reset(tab.nextRevalidateTime())
		case <-copyNodes.C:
			//30秒来一下，持久化符合条件的节点到数据库中
			log.Info("copyLiveNodes","now",time.Now().Unix())
			go tab.copyLiveNodes()
		case <-tab.closeReq:
			break loop
		}
	}

	if tab.net != nil {
		tab.net.close()
	}
	if refreshDone != nil {
		<-refreshDone
	}
	for _, ch := range waiting {
		close(ch)
	}
	tab.db.close()
	close(tab.closed)
}
```





要建立tcp连接的时候，就从k桶中随机选取节点连接

```
// ReadRandomNodes fills the given slice with random nodes from the
// table. It will not write the same node more than once. The nodes in
// the slice are copies and can be modified by the caller.
func (tab *Table) ReadRandomNodes(buf []*Node) (n int) {
	if !tab.isInitDone() {
		return 0
	}
	tab.mutex.Lock()
	defer tab.mutex.Unlock()

	// Find all non-empty buckets and get a fresh slice of their entries.
	var buckets [][]*Node
	for _, b := range tab.buckets {
		if len(b.entries) > 0 {
			buckets = append(buckets, b.entries[:])
		}
	}
	if len(buckets) == 0 {
		return 0
	}
	// Shuffle the buckets.
	for i := len(buckets) - 1; i > 0; i-- {
		j := tab.rand.Intn(len(buckets))
		buckets[i], buckets[j] = buckets[j], buckets[i]
	}
	// Move head of each bucket into buf, removing buckets that become empty.
	var i, j int
	for ; i < len(buf); i, j = i+1, (j+1)%len(buckets) {
		b := buckets[j]
		buf[i] = &(*b[0])
		buckets[j] = b[1:]
		if len(b) == 1 {
			buckets = append(buckets[:j], buckets[j+1:]...)
		}
		if len(buckets) == 0 {
			break
		}
	}
	return i + 1
}
```



table中的重头戏

```
// doRefresh performs a lookup for a random target to keep buckets
// full. seed nodes are inserted if the table is empty (initial
// bootstrap or discarded faulty peers).
func (tab *Table) doRefresh(done chan struct{}) {

	defer close(done)

	// Load nodes from the database and insert
	// them. This should yield a few previously seen nodes that are
	// (hopefully) still alive.
	tab.loadSeedNodes()

	// Run self lookup to discover new neighbor nodes.
	tab.lookup(tab.self.ID, false)

	// The Kademlia paper specifies that the bucket refresh should
	// perform a lookup in the least recently used bucket. We cannot
	// adhere to this because the findnode target is a 512bit value
	// (not hash-sized) and it is not easily possible to generate a
	// sha3 preimage that falls into a chosen bucket.
	// We perform a few lookups with a random target instead.
	for i := 0; i < 3; i++ {
		var target NodeID
		crand.Read(target[:])
		tab.lookup(target, false)
	}
}
```





```
func (tab *Table) lookup(targetID NodeID, refreshIfEmpty bool) []*Node {
	var (
		target         = crypto.Keccak256Hash(targetID[:])
		asked          = make(map[NodeID]bool)
		seen           = make(map[NodeID]bool)
		reply          = make(chan []*Node, alpha)
		pendingQueries = 0
		result         *nodesByDistance
	)
	// don't query further if we hit ourself.
	// unlikely to happen often in practice.
	asked[tab.self.ID] = true

	for {
		tab.mutex.Lock()
		// generate initial result set
		result = tab.closest(target, bucketSize)
		tab.mutex.Unlock()
		if len(result.entries) > 0 || !refreshIfEmpty {
			break
		}
		// The result set is empty, all nodes were dropped, refresh.
		// We actually wait for the refresh to complete here. The very
		// first query will hit this case and run the bootstrapping
		// logic.
		<-tab.refresh()
		refreshIfEmpty = false
	}

	for {
		// ask the alpha closest nodes that we haven't asked yet
		for i := 0; i < len(result.entries) && pendingQueries < alpha; i++ {
			n := result.entries[i]
			if !asked[n.ID] {
				asked[n.ID] = true
				pendingQueries++
				go tab.findnode(n, targetID, reply)
			}
		}
		if pendingQueries == 0 {
			// we have asked all closest nodes, stop the search
			break
		}
		//这里的代码有点意思，如果说pendingQueries不为0，说明他肯定去询问了一个节点，要求邻居节点数据
		//这里却只需要等待一个回来。
		// wait for the next reply
		for _, n := range <-reply {
			if n != nil && !seen[n.ID] {
				seen[n.ID] = true
				result.push(n, bucketSize)
			}
		}
		pendingQueries--
	}
	return result.entries
}
```





```
func (tab *Table) findnode(n *Node, targetID NodeID, reply chan<- []*Node) {
	fails := tab.db.findFails(n.ID)
	r, err := tab.net.findnode(n.ID, n.addr(), targetID)
	if err != nil || len(r) == 0 {
		fails++
		tab.db.updateFindFails(n.ID, fails)
		log.Trace("Findnode failed", "id", n.ID, "failcount", fails, "err", err)
		if fails >= maxFindnodeFailures {
			log.Trace("Too many findnode failures, dropping", "id", n.ID, "failcount", fails)
			tab.delete(n)
		}
	} else if fails > 0 {
		tab.db.updateFindFails(n.ID, fails-1)
	}

	// Grab as many nodes as possible. Some of them might not be alive anymore, but we'll
	// just remove those again during revalidation.
	//我们发现findnode如果能够得到结果，那么就必然试着将它写入k桶
	for _, n := range r {
		tab.add(n)
	}
	reply <- r
}
```

