##1、本地缓存系统##


####<font color="red">原理</font>####

在内存中维护一个cache。查询时，首先查询cache中是否已经缓存查询结果，如果已经缓存，直接返回缓存结果，如果没有缓存，将查询结果Load到cache中，然后返回结果。



##2、分布式缓存系统##

在设计分布式缓存系统的时候，需要让key分布在不同的缓存节点上。当某节点收到查询请求时，如果该key归属于本节点，则在本节点获取查询结果，如果该key归属于其他节点，则本节点向归属节点获取查询结果。

从这里我们可以看出分布式缓存系统是在本地系统上进行功能扩展，也就是当key归属于本节点时，操作就如同本地缓存系统。当key不属于本节点时，本节点就想归属节点获取查询结果（这个就如同人家向你问路，如果你认识所问的地方，你就告诉人家，如果你不认识，你是个热心的人，百度得到结果了告诉人家）。

我们需要解决一下问题

1. 如何判断查询的key归属哪个节点

2. 如何从其他节点获取数据

第1个问题实际上是一个路由问题，即给定key，路由到某个节点。这里忽略具体实现，将问题抽象出来，由2个接口表示（对应于groupcache的ProtoGetter和PeerPicker）：

```
type Peer interface{
    Get(key string) Value
}

type Router interface{
    Route(key string) Peer
}
```

接口Peer表示一个远端节点，Get方法从远端节点查询数据；接口Router表示一个路由器，Route方法将给定key路由到其归属的远端节点，如果key归属于本节点，Route方法返回nil。基于此，改写Group结构体和Group.load方法

```
type Group struct {
    router Router
    cache  Cache
}
func (g *Group) load(key string) Value {
    peer := g.router.Route(key)
    if peer == nil {
        val := g.getLocally(key)
        g.cache.Add(key, val) //store result in cache
        return val
    }
    return peer.Get(g.name, key)
}
```

至此，一个分布式缓存框架搭建完成了。


2.1. Cache的内存结构

Cache对象是内存中的一个容器，用来存放查询的结果。我们将Cache对象实现为一个List。另外，为了提升访问效率，用一个map结构来索引key及其对应的值：

```
type ListCache struct {
    MaxEntries int
    l list.List
    m map[string]*list.Element
}
```

指定LRU为缓存的淘汰策略：即当Cache满了需要淘汰数据时，优先淘汰最老的数据。我们约定，链表中越靠近表头数据越新，越靠近表尾数据越老。因此在访问cache时，需要把命中的数据移至表头：
     
 ```
 type entry struct {
    key string
    val Value
}
func (c *ListCache) Get(key string) (Value, bool) {
    if e, hit := c.m[key]; hit {
        c.l.MoveToFront(e)
        return e.Value.(*entry).val, hit
    }
    return nil, false
}
func (c *ListCache) Add(key string, val Value) {
    if e, ok := c.m[key]; ok {
        c.l.MoveToFront(e)
    } else {
        ee := c.l.PushFront(&entry{key, val})
        c.m[key] = ee
        if c.l.Len() > c.MaxEntries {
            oldest := c.l.Back()
            c.l.Remove(oldest)
            delete(c.m, oldest.Value.(*entry).key)
        }
    }
}
 
 ```
 
 2.2. 路由器的实现——一致性哈希
 
路由器的基本功能是将key分布到不同的节点。很多方法可以实现这一点，最常用的方法是计算哈希。例如对于每次访问，可以按如下算法计算其哈希值：
         
         h= Hash(key) % N

这个算式计算每个key应该被路由到哪个节点，其中N为节点总数，所有节点按照0 – (N-1)编号。

 这个算法的问题在于容错性和扩展性不好。假设有一台服务器宕机了，那么为了填补空缺，要将宕机的服务器从编号列表中移除，后面的服务器按顺序前移一位并将其编号值减一，此时每个key就要按h = Hash(key) % (N-1)重新计算；同样，如果新增了一台服务器，虽然原有服务器编号不用改变，但是要按h = Hash(key) % (N+1)重新计算哈希值。因此系统中一旦有服务器变更，大量的key会被重定位到不同的服务器从而造成大量的缓存不命中。

            一致性哈希算法就是解决这个问题的一种哈希方案。

简单来说，一致性哈希将整个哈希值空间组织成一个虚拟的圆环，如假设某哈希函数H的值空间为0 - 232-1（即哈希值是一个32位无符号整形），整个空间按顺时针方向组织。0和232-1在零点中方向重合。将各个节点使用H进行一个哈希，确定其在哈希环上的位置。

将数据key使用相同的函数H计算出哈希值h，根据h确定此数据在环上的位置，从该位置<font color="red">沿环顺时针“行走”</font>，第一台遇到的节点就是其应该路由到的节点。


要在一个节点上启动上述的分布式缓存系统，包含如下步骤：

1.    初始化一个Group对象，指定Loader

2.    初始化一个ListCache对象，并将其赋值给Group对象的cache

3.    初始化一个HashRouter对象，并将其赋值给Group对象的router

4.    为集群的每个节点初始化一个HTTPPeer对象，并将它们添加到HashRouter对象。将本节点对应的HTTPPeer对象赋值给HashRouter对象的self

5.    http监听指定端口，并指定处理函数为HTTPPeer对象的ServeHTTP方法

可以做一些适当的封装，让此缓存系统使用更加方便。


保证对相同key的访问会被发送到相同的服务器。当某一台主机宕机了，不可能说都没有影响，我们要做的是把影响降到最低。


##一致性哈希算法##

抽象的说法就是我们假设有这样的一个给定的哈希函数H,我们计算节点的位置和计算key的位置都使用哈希函数H。

简单来说，一致性哈希将整个哈希值空间组织成一个虚拟的圆环，如假设某哈希函数H的值空间为0 - (2的32次方减1)（即哈希值是一个32位无符号整形）.整个空间按顺时针方向组织。0和(2的32次方减1)在零点中方向重合。

下一步将各个服务器使用H进行一个哈希，具体可以选择服务器的ip或主机名作为关键字进行哈希，这样每台机器就能确定其在哈希环上的位置。

我们使用主机名或者服务器的ip进行hash,计算出该节点的哈希值，然后根据这个哈希值来放置节点（形象的说法）。

接下来使用如下算法定位数据访问到相应服务器：将数据key使用相同的函数H计算出哈希值h，通根据h确定此数据在环上的位置，从此位置沿环顺时针“行走”，第一台遇到的服务器就是其应该定位到的服务器。


###容错性与可扩展性分析###

* 一般的，在一致性哈希算法中，如果一台服务器不可用，则受影响的数据仅仅是此服务器到其环空间中前一台服务器（即顺着逆时针方向行走遇到的第一台服务器）之间数据，其它不会受到影响。
* 一般的，在一致性哈希算法中，如果增加一台服务器，则受影响的数据仅仅是新服务器到其环空间中前一台服务器（即顺着逆时针方向行走遇到的第一台服务器）之间数据，其它不会受到影响。

综上所述，一致性哈希算法对于节点的增减都只需重定位环空间中的一小部分数据，具有较好的容错性和可扩展性。

###虚拟节点###
一致性哈希算法在服务节点太少时，容易因为节点分部不均匀而造成数据倾斜问题。例如我们的系统中有两台服务器Server 1和Server 2，在环上它们靠得很近。

此时必然造成大量数据集中到Server 1上，而只有极少量会定位到Server 2上。为了解决这种数据倾斜问题，一致性哈希算法引入了虚拟节点机制，即对每一个服务节点计算多个哈希，每个计算结果位置都放置一个此服务节点，称为虚拟节点。具体做法可以在服务器ip或主机名的后面增加编号来实现。例如上面的情况，我们决定为每台服务器计算三个虚拟节点，于是可以分别计算“Memcached Server 1#1”、“Memcached Server 1#2”、“Memcached Server 1#3”、“Memcached Server 2#1”、“Memcached Server 2#2”、“Memcached Server 2#3”的哈希值，于是形成六个虚拟节点

同时数据定位算法不变，只是多了一步虚拟节点到实际节点的映射，例如定位到“Memcached Server 1#1”、“Memcached Server 1#2”、“Memcached Server 1#3”三个虚拟节点的数据均定位到Server 1上。这样就解决了服务节点少时数据倾斜的问题。在实际应用中，通常将虚拟节点数设置为32甚至更大，因此即使很少的服务节点也能做到相对均匀的数据分布。

