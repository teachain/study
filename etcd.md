# ETCD

比较给力的链接 

https://github.com/coreos/etcd/releases   etcd 发行版本

https://www.cnblogs.com/xigang8068/p/5786027.html ectd 服务注册与发现

http://thesecretlivesofdata.com/raft/  raft的原理。

```
A highly-available key value store for shared configuration and service discovery.
```

服务发现（Service Discovery）要解决的是分布式系统中最常见的问题之一，即在同一个分布式集群中的进程或服务如何才能找到对方并建立连接。从本质上说，服务发现就是想要了解集群中是否有进程在监听udp或tcp端口，并且通过名字就可以进行查找和连接。要解决服务发现的问题，需要有下面三大支柱，缺一不可。

- 一个强一致性、高可用的服务存储目录。基于Raft算法的etcd天生就是这样一个强一致性高可用的服务存储目录。
- 一种注册服务和监控服务健康状态的机制。用户可以在etcd中注册服务，并且对注册的服务设置`key TTL`，定时保持服务的心跳以达到监控健康状态的效果。
- 一种查找和连接服务的机制。通过在etcd指定的主题下注册的服务也能在对应的主题下查找到。为了确保连接，我们可以在每个服务机器上都部署一个proxy模式的etcd，这样就可以确保能访问etcd集群的服务都能互相连接。



```
*设置一个key的value
curl -s http://127.0.0.1:2379/v2/keys/message -X PUT -d value="Hello world"
```

```
*获取一个key的value
curl -s http://127.0.0.1:2379/v2/keys/message
```

```
*改变一个key的value
curl -s http://127.0.0.1:2379/v2/keys/message -X PUT -d value="Hello etcd"
```

```
*删除一个key节点
curl -s http://127.0.0.1:2379/v2/keys/message -X DELETE
```

```
*使用ttl（即设置一个key的值并给这个key加一个生命周期，当超过这个时间该值没有被访问则自动被删除）
curl -s http://127.0.0.1:2379/v2/keys/foo -X PUT -d value=bar -d ttl=5
```

```
*watch一个值的变化
curl -s http://127.0.0.1:2379/v2/keys/foo?wait=true

该命令调用之后会阻塞进程，直到这个值发生变化才能返回，当改变一个key的值，或者删除等操作发生时，该等待就会返回结果。
```

```
*创建一个目录

curl -s http://127.0.0.1:2379/v2/keys/dir -X PUT -d dir=true
```

```
*列举一个目录
curl -s http://127.0.0.1:2379/v2/keys/dir
```

```
*递归列举一个目录
curl -s http://127.0.0.1:2379/v2/keys/dir?recursive=true
```

```
监控一个目录下的所有key的变化，包括子目录的。可以使用命令：

curl -s http://127.0.0.1:2379/v2/keys/dir?recursive=true&wait=true
```

```
*删除一个目录

curl -s http://127.0.0.1:2379/v2/keys/dir?dir=true -X DELETE
```

### etcd一般部署集群推荐奇数个节点（小数服从多数）

etcd作为一个高可用键值存储系统，天生就是为集群化而设计的。由于Raft算法在做决策时需要多数节点的投票，所以etcd一般部署集群推荐奇数个节点，推荐的数量为3、5或者7个节点构成一个集群。





worker启动时向etcd注册自己的信息,并设置一个过期时间TTL,每隔一段时间更新这个TTL,如果该worker挂掉了,这个TTL就会expire. master则监听`workers/`这个etcd directory, 根据检测到的不同action来增加, 更新, 或删除worker.