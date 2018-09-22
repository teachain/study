##ETCD##

<font color="red">它是一个键值存储仓库,用于配置共享和服务发现</font>

A highly-available key value store for shared configuration and service discovery.

它具有以下4个特点

* 简单：基于HTTP+JSON的API让你用curl命令就可以轻松使用。

* 安全：可选SSL客户认证机制。

* 快速：每个实例每秒支持一千次写操作。

* 可信：使用Raft算法充分实现了分布式。


###etcd集群模式###

集群的搭建可分为三种方式

* 静态配置

* etcd发现

* dns发现



配置项说明

* --name
etcd集群中的节点名，这里可以随意，可区分且不重复就行  

* --listen-peer-urls
监听的用于节点之间通信的url，可监听多个，集群内部将通过这些url进行数据交互(如选举，数据同步等)

* --initial-advertise-peer-urls 
建议用于节点之间通信的url，节点间将以该值进行通信。

* --listen-client-urls
监听的用于客户端通信的url,同样可以监听多个。

* --advertise-client-urls
建议使用的客户端通信url,该值用于etcd代理或etcd成员与etcd节点通信。

* --initial-cluster-token etcd-cluster-1
节点的token值，设置该值后集群将生成唯一id,并为每个节点也生成唯一id,当使用相同配置文件再启动一个集群时，只要该token值不一样，etcd集群就不会相互影响。

* --initial-cluster
也就是集群中所有的initial-advertise-peer-urls 的合集

* --initial-cluster-state new
新建集群的标志

###1、静态配置##

静态配置主要预先将集群的配置信息分配好，然后将集群分布启动，集群将根据配置信息组成集群。这里按如下的配置信息分别启动三个etcd。


./etcd --name infra0 --initial-advertise-peer-urls http://10.0.1.111:2380 \
  --listen-peer-urls http://10.0.1.111:2380 \
  --listen-client-urls http://10.0.1.111:2379,http://127.0.0.1:2379 \
  --advertise-client-urls http://10.0.1.111:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster infra0=http://10.0.1.111:2380,infra1=http://10.0.1.109:2380,infra2=http://10.0.1.110:2380 \
  --initial-cluster-state new

./etcd --name infra1 --initial-advertise-peer-urls http://10.0.1.109:2380 \
  --listen-peer-urls http://10.0.1.109:2380 \
  --listen-client-urls http://10.0.1.109:2379,http://127.0.0.1:2379 \
  --advertise-client-urls http://10.0.1.109:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster infra0=http://10.0.1.111:2380,infra1=http://10.0.1.109:2380,infra2=http://10.0.1.110:2380 \
  --initial-cluster-state new

./etcd --name infra2 --initial-advertise-peer-urls http://10.0.1.110:2380 \
  --listen-peer-urls http://10.0.1.110:2380 \
  --listen-client-urls http://10.0.1.110:2379,http://127.0.0.1:2379 \
  --advertise-client-urls http://10.0.1.110:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster infra0=http://10.0.1.111:2380,infra1=http://10.0.1.109:2380,infra2=http://10.0.1.110:2380 \
  --initial-cluster-state new
  
  
按如上配置分别启动集群，启动集群后，将会进入集群选举状态，若出现大量超时，则需要检查主机的防火墙是否关闭，或主机之间是否能通过2380端口通信，集群建立后通过以下命令检查集群状态。

