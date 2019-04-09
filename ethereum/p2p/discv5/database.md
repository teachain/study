



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

// Contains the node database, storing previously seen nodes and any collected
// metadata about them for QoS purposes.

package discv5

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/simplechain-org/go-simplechain/crypto"
	"github.com/simplechain-org/go-simplechain/log"
	"github.com/simplechain-org/go-simplechain/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	nodeDBNilNodeID      = NodeID{}       // Special node ID to use as a nil element.
	nodeDBNodeExpiration = 24 * time.Hour // Time after which an unseen node should be dropped.
	nodeDBCleanupCycle   = time.Hour      // Time period for running the expiration task.
)

// nodeDB stores all nodes we know about.
type nodeDB struct {
	lvl    *leveldb.DB   // Interface to the database itself
	self   NodeID        // Own node id to prevent adding it into the database
	runner sync.Once     // Ensures we can start at most one expirer
	quit   chan struct{} // Channel to signal the expiring thread to stop
}

// Schema layout for the node database
var (
	nodeDBVersionKey = []byte("version") // Version of the database to flush if changes
	nodeDBItemPrefix = []byte("n:")      // Identifier to prefix node entries with

	nodeDBDiscoverRoot          = ":discover"
	nodeDBDiscoverPing          = nodeDBDiscoverRoot + ":lastping"
	nodeDBDiscoverPong          = nodeDBDiscoverRoot + ":lastpong"
	nodeDBDiscoverFindFails     = nodeDBDiscoverRoot + ":findfail"
	nodeDBDiscoverLocalEndpoint = nodeDBDiscoverRoot + ":localendpoint"
	nodeDBTopicRegTickets       = ":tickets"
)

// newNodeDB creates a new node database for storing and retrieving infos about
// known peers in the network. If no path is given, an in-memory, temporary
// database is constructed.
func newNodeDB(path string, version int, self NodeID) (*nodeDB, error) {
	//从这里我们可以看到，当不提供path路径的时候，直接使用leveldb内存来存储数据
	if path == "" {
		return newMemoryNodeDB(self)
	}
	//提供有效的path的时候，我们使用leveldb的文件来存储数据
	return newPersistentNodeDB(path, version, self)
}
//内存数据库
// newMemoryNodeDB creates a new in-memory node database without a persistent
// backend.
func newMemoryNodeDB(self NodeID) (*nodeDB, error) {
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		return nil, err
	}
	return &nodeDB{
		lvl:  db,
		self: self,
		quit: make(chan struct{}),
	}, nil
}
//文件数据库
// newPersistentNodeDB creates/opens a leveldb backed persistent node database,
// also flushing its contents in case of a version mismatch.
func newPersistentNodeDB(path string, version int, self NodeID) (*nodeDB, error) {
	opts := &opt.Options{OpenFilesCacheCapacity: 5}
	//这里实际的意思应该是，我需要打开一个目录，里面任我创建文件，读写文件。
	db, err := leveldb.OpenFile(path, opts)
	//如果里面的文件有损，我们试着恢复
	if _, iscorrupted := err.(*errors.ErrCorrupted); iscorrupted {
		db, err = leveldb.RecoverFile(path, nil)
	}
	//如果出错，返回错误
	if err != nil {
		return nil, err
	}
	// The nodes contained in the cache correspond to a certain protocol version.
	// Flush all nodes if the version doesn't match.
	currentVer := make([]byte, binary.MaxVarintLen64)
	currentVer = currentVer[:binary.PutVarint(currentVer, int64(version))]
    //如果到了，这里，表示上面的没有错误
    //我们检测一下，看是否存在版本号
	blob, err := db.Get(nodeDBVersionKey, nil)
	switch err {
	case leveldb.ErrNotFound:
	    //表示我们没有写过，就是还没有初始化，那我们就把版本号写入数据库中
		// Version not found (i.e. empty cache), insert it
		if err := db.Put(nodeDBVersionKey, currentVer, nil); err != nil {
			//如果写入有错误，关闭数据库，返回错误信息
			db.Close()
			return nil, err
		}
	case nil:
	    //表示已经初始化过了，那么我们比对一下版本号
		// Version present, flush if different
		if !bytes.Equal(blob, currentVer) {
		   //版本号不一致，我们关闭数据库
			db.Close()
			//将目录下的所有文件全部删除
			if err = os.RemoveAll(path); err != nil {
				return nil, err
			}
			//重建(递归调用本方法)
			//这样它就会走上面那个分支，然后到方法的后面去返回nodedb对象
			return newPersistentNodeDB(path, version, self)
		}
	}
	return &nodeDB{
		lvl:  db,
		self: self,
		quit: make(chan struct{}),
	}, nil
}

// makeKey generates the leveldb key-blob from a node id and its particular
// field of interest.
func makeKey(id NodeID, field string) []byte {
	//要检测一下，id不能是一个空的id,必须是一个有效的id
	if bytes.Equal(id[:], nodeDBNilNodeID[:]) {
		return []byte(field)
	}
	//返回  "n:"+id+field,+表示拼接的意思
	return append(nodeDBItemPrefix, append(id[:], field...)...)
}

//根据一个key区分id和field
// splitKey tries to split a database key into a node id and a field part.
func splitKey(key []byte) (id NodeID, field string) {
    //首先是这个key必须是n:开头的，否则不是node的key
	// If the key is not of a node, return it plainly
	if !bytes.HasPrefix(key, nodeDBItemPrefix) {
		return NodeID{}, string(key)
	}
	// Otherwise split the id and field
	//先去掉前缀n：
	item := key[len(nodeDBItemPrefix):]
	//然后根据id的定长特性(NodeID是64字节的)，取出id
	copy(id[:], item[:len(id)])
	//去掉id部分，自然就是field部分
	field = string(item[len(id):])
    //返回id和field
	return id, field
}

// fetchInt64 retrieves an integer instance associated with a particular
// database key.
//从数据库中获取一个key对应的value(key=>value)
func (db *nodeDB) fetchInt64(key []byte) int64 {
	blob, err := db.lvl.Get(key, nil)
	if err != nil {
		return 0
	}
	val, read := binary.Varint(blob)
	if read <= 0 {
		return 0
	}
	return val
}
//将key=>value存入数据库中
// storeInt64 update a specific database entry to the current time instance as a
// unix timestamp.
func (db *nodeDB) storeInt64(key []byte, n int64) error {
	blob := make([]byte, binary.MaxVarintLen64)
	blob = blob[:binary.PutVarint(blob, n)]
	return db.lvl.Put(key, blob, nil)
}
//存入一块数据(key=>value)
func (db *nodeDB) storeRLP(key []byte, val interface{}) error {
	blob, err := rlp.EncodeToBytes(val)
	if err != nil {
		return err
	}
	return db.lvl.Put(key, blob, nil)
}
//读取一块数据(key=>value)
func (db *nodeDB) fetchRLP(key []byte, val interface{}) error {
	blob, err := db.lvl.Get(key, nil)
	if err != nil {
		return err
	}
	err = rlp.DecodeBytes(blob, val)
	if err != nil {
		log.Warn(fmt.Sprintf("key %x (%T) %v", key, val, err))
	}
	return err
}

//节点数据读取
//根据id从数据库中读取一个节点的数据
//Node存入数据库中的key为:n:+id+:discover
// node retrieves a node with a given id from the database.
func (db *nodeDB) node(id NodeID) *Node {
	var node Node
	if err := db.fetchRLP(makeKey(id, nodeDBDiscoverRoot), &node); err != nil {
		return nil
	}
	node.sha = crypto.Keccak256Hash(node.ID[:])
	return &node
}
//节点数据写入
// updateNode inserts - potentially overwriting - a node into the peer database.
func (db *nodeDB) updateNode(node *Node) error {
	return db.storeRLP(makeKey(node.ID, nodeDBDiscoverRoot), node)
}
//节点数据删除
//只要是:n:+id开头的都删除掉，就是有关于此id的数据都删除
// deleteNode deletes all information/keys associated with a node.
func (db *nodeDB) deleteNode(id NodeID) error {
	deleter := db.lvl.NewIterator(util.BytesPrefix(makeKey(id, "")), nil)
	for deleter.Next() {
		if err := db.lvl.Delete(deleter.Key(), nil); err != nil {
			return err
		}
	}
	return nil
}

// ensureExpirer is a small helper method ensuring that the data expiration
// mechanism is running. If the expiration goroutine is already running, this
// method simply returns.
//
// The goal is to start the data evacuation only after the network successfully
// bootstrapped itself (to prevent dumping potentially useful seed nodes). Since
// it would require significant overhead to exactly trace the first successful
// convergence, it's simpler to "ensure" the correct state when an appropriate
// condition occurs (i.e. a successful bonding), and discard further events.
//只运行一次的定时清理器
func (db *nodeDB) ensureExpirer() {
	db.runner.Do(func() { go db.expirer() })
}

// expirer should be started in a go routine, and is responsible for looping ad
// infinitum and dropping stale data from the database.
//断续器，一个小时执行一次，我们需要注意的是，当应用启动一个小时以后才会执行一次清理工作。
func (db *nodeDB) expirer() {
	tick := time.NewTicker(nodeDBCleanupCycle)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if err := db.expireNodes(); err != nil {
				log.Error(fmt.Sprintf("Failed to expire nodedb items: %v", err))
			}
		case <-db.quit:
			return
		}
	}
}
//查找节点发现的节点数据
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
		//确定必须是节点发现的数据
		if field != nodeDBDiscoverRoot {
		    //不是就略过
			continue
		}
		// Skip the node if not expired yet (and not self)
		// 不是本节点
		if !bytes.Equal(id[:], db.self[:]) {
		    //它的最近的pong消息的时间在一天之内，我们就保留
			if seen := db.lastPong(id); seen.After(threshold) {
			    //所以这里我们就略过
				continue
			}
		}
		//这里就是在一天之内的节点数据，我们将他们从数据库中删除
		// Otherwise delete all associated information
		db.deleteNode(id)
	}
	return nil
}
//在这里我们发现它的节点信息和节点的ping,pong时间是分开存储的
//ping消息的记录时间的key是 n:+id+:discover+:lastping
// lastPing retrieves the time of the last ping packet send to a remote node,
// requesting binding.
func (db *nodeDB) lastPing(id NodeID) time.Time {
	return time.Unix(db.fetchInt64(makeKey(id, nodeDBDiscoverPing)), 0)
}

// updateLastPing updates the last time we tried contacting a remote node.
func (db *nodeDB) updateLastPing(id NodeID, instance time.Time) error {
	return db.storeInt64(makeKey(id, nodeDBDiscoverPing), instance.Unix())
}
//pong消息的记录时间的key是 n:+id+:discover+:lastpong
// lastPong retrieves the time of the last successful contact from remote node.
func (db *nodeDB) lastPong(id NodeID) time.Time {
	return time.Unix(db.fetchInt64(makeKey(id, nodeDBDiscoverPong)), 0)
}

// updateLastPong updates the last time a remote node successfully contacted.
func (db *nodeDB) updateLastPong(id NodeID, instance time.Time) error {
	return db.storeInt64(makeKey(id, nodeDBDiscoverPong), instance.Unix())
}
//查找邻居节点失败次数的key是n:+id+:discover+:findfail
// findFails retrieves the number of findnode failures since bonding.
func (db *nodeDB) findFails(id NodeID) int {
	return int(db.fetchInt64(makeKey(id, nodeDBDiscoverFindFails)))
}

// updateFindFails updates the number of findnode failures since bonding.
func (db *nodeDB) updateFindFails(id NodeID, fails int) error {
	return db.storeInt64(makeKey(id, nodeDBDiscoverFindFails), int64(fails))
}
//本地key是n:+id+:discover+:localendpoint
// localEndpoint returns the last local endpoint communicated to the
// given remote node.
func (db *nodeDB) localEndpoint(id NodeID) *rpcEndpoint {
	var ep rpcEndpoint
	if err := db.fetchRLP(makeKey(id, nodeDBDiscoverLocalEndpoint), &ep); err != nil {
		return nil
	}
	return &ep
}

func (db *nodeDB) updateLocalEndpoint(id NodeID, ep rpcEndpoint) error {
	return db.storeRLP(makeKey(id, nodeDBDiscoverLocalEndpoint), &ep)
}

// querySeeds retrieves random nodes to be used as potential seed nodes
// for bootstrapping.
func (db *nodeDB) querySeeds(n int, maxAge time.Duration) []*Node {
	var (
		now   = time.Now()
		nodes = make([]*Node, 0, n)
		it    = db.lvl.NewIterator(nil, nil)
		id    NodeID
	)
	defer it.Release()

seek:
    //这里可以这么理解，假设n等于10，那么我们就最多尝试n*5也就是50次
	for seeks := 0; len(nodes) < n && seeks < n*5; seeks++ {
		// Seek to a random entry. The first byte is incremented by a
		// random amount each time in order to increase the likelihood
		// of hitting all existing nodes in very small databases.
		
		ctr := id[0]
		
		rand.Read(id[:])
		
		id[0] = ctr + id[0]%16
		
		//寻找一个随机条目。第一个字节每次递增一个随机数，以增加命中非常小数据库中所有现有节点的可能性。
		//意思即是数据量小的时候也能容易找到节点
		it.Seek(makeKey(id, nodeDBDiscoverRoot))

		n := nextNode(it)
		
		//找不到，继续尝试
		if n == nil {
			id[0] = 0
			continue seek // iterator exhausted
		}
		//是本节点，继续尝试
		if n.ID == db.self {
			continue seek
		}
		//时间离得太久，也不要，继续尝试
		if now.Sub(db.lastPong(n.ID)) > maxAge {
			continue seek
		}
		//如果已经在集合里，也不要，继续尝试
		for i := range nodes {
			if nodes[i].ID == n.ID {
				continue seek // duplicate
			}
		}
		//把选好的节点数据放入集合中
		nodes = append(nodes, n)
	}
	return nodes
}
//获取
//主题注册ticket的key是 n:+id+:tickets
func (db *nodeDB) fetchTopicRegTickets(id NodeID) (issued, used uint32) {
	key := makeKey(id, nodeDBTopicRegTickets)
	blob, _ := db.lvl.Get(key, nil)
	if len(blob) != 8 {
		return 0, 0
	}
	//发行时间
	issued = binary.BigEndian.Uint32(blob[0:4])
	//使用事件
	used = binary.BigEndian.Uint32(blob[4:8])
	return
}
//保存或更新
func (db *nodeDB) updateTopicRegTickets(id NodeID, issued, used uint32) error {
	key := makeKey(id, nodeDBTopicRegTickets)
	blob := make([]byte, 8)
	binary.BigEndian.PutUint32(blob[0:4], issued)
	binary.BigEndian.PutUint32(blob[4:8], used)
	return db.lvl.Put(key, blob, nil)
}

// reads the next node record from the iterator, skipping over other
// database entries.
//从迭代器里将节点信息读取出来
func nextNode(it iterator.Iterator) *Node {
	for end := false; !end; end = !it.Next() {
		id, field := splitKey(it.Key())
		if field != nodeDBDiscoverRoot {
			continue
		}
		var n Node
		if err := rlp.DecodeBytes(it.Value(), &n); err != nil {
			log.Warn(fmt.Sprintf("invalid node %x: %v", id, err))
			continue
		}
		return &n
	}
	return nil
}

// close flushes and closes the database files.
func (db *nodeDB) close() {
	close(db.quit)
	db.lvl.Close()
}

```

