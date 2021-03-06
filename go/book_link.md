##golang常用库链接##

https://github.com/gopherchina


###1、日志库###
<a href="https://github.com/cihub/seelog">https://github.com/cihub/seelog</a>

###2、NSQ：分布式的实时消息平台###

<a href="https://github.com/nsqio/nsq">前往查看源码</a>

NSQ是一个基于Go语言的分布式实时消息平台，它基于MIT开源协议发布，代码托管在GitHub，其当前最新版本是0.3.1版。NSQ可用于大规模系统中的实时消息服务，并且每天能够处理数亿级别的消息，其设计目标是为在分布式环境下运行的去中心化服务提供一个强大的基础架构。NSQ具有分布式、去中心化的拓扑结构，该结构具有无单点故障、故障容错、高可用性以及能够保证消息的可靠传递的特征。NSQ非常容易配置和部署，且具有最大的灵活性，支持众多消息协议。另外，官方还提供了拆箱即用Go和Python库。如果读者兴趣构建自己的客户端的话，还可以参考官方提供的协议规范。

NSQ是由四个重要组件构成：

* nsqd：一个负责接收、排队、转发消息到客户端的守护进程
* nsqlookupd：管理拓扑信息并提供最终一致性的发现服务的守护进程
* nsqadmin：一套Web用户界面，可实时查看集群的统计数据和执行各种各样的管理任务
* utilities：常见基础功能、数据流处理工具，如nsq_stat、nsq_tail、nsq_to_file、nsq_to_http、nsq_to_nsq、to_nsq

NSQ的主要特点如下:

* 具有分布式且无单点故障的拓扑结构 支持水平扩展，在无中断情况下能够无缝地添加集群节点
* 低延迟的消息推送，参见官方提供的性能说明文档
* 具有组合式的负载均衡和多播形式的消息路由
* 既擅长处理面向流（高吞吐量）的工作负载，也擅长处理面向Job的（低吞吐量）工作负载
* 消息数据既可以存储于内存中，也可以存储在磁盘中
* 实现了生产者、消费者自动发现和消费者自动连接生产者，参见nsqlookupd
* 支持安全传输层协议（TLS），从而确保了消息传递的安全性
* 具有与数据格式无关的消息结构，支持JSON、Protocol Buffers、MsgPack等消息格式
* 非常易于部署（几乎没有依赖）和配置（所有参数都可以通过命令行进行配置）
* 使用了简单的TCP协议且具有多种语言的客户端功能库
* 具有用于信息统计、管理员操作和实现生产者等的HTTP接口
* 为实时检测集成了统计数据收集器StatsD
* 具有强大的集群管理界面，参见nsqadmin


为了达到高效的分布式消息服务，NSQ实现了合理、智能的权衡，从而使得其能够完全适用于生产环境中，具体内容如下：

* 支持消息内存队列的大小设置，默认完全持久化（值为0），消息即可持久到磁盘也可以保存在内存中
* 保证消息至少传递一次,以确保消息可以最终成功发送
* 收到的消息是<font color="red">无序的</font>, 实现了松散订购
* 发现服务nsqlookupd具有最终一致性,消息最终能够找到所有Topic生产者

<font color="red">注意:因为在NSQ中消息是无序的，所以它并不能作为消息队列使用。</font>

###3、“令牌桶”算法进行限流###

<a href="https://github.com/juju/ratelimit">很简单轻量的包 </a>

<a href="https://en.wikipedia.org/wiki/Token_bucket">令牌桶算法</a>

常用的限流算法有两种：漏桶算法和令牌桶算法。

*  Leaky Bucket 漏桶算法思路很简单，水（请求）先进入到漏桶里，漏桶以一定的速度出水，当水流入速度过大会直接溢出，可以看出漏桶算法能强行限制数据的传输速率。
*  token bucket 令牌桶算法的原理是系统会以一个恒定的速度往桶里放入令牌，而如果请求需要被处理，则需要先从桶里获取一个令牌，当桶里没有令牌可取时，则拒绝服务。


####令牌桶工作参数####

工作过程包括3个阶段：产生令牌、消耗令牌和判断数据包是否通过。其中涉及到2个参数：令牌产生的速率CIR（Committed Information Rate）/EIR（Excess Information Rate）和令牌桶的大小CBS（Committed Burst Size）/EBS（Excess Burst Size）。

* 产生令牌：周期性的以速率CIR/EIR向令牌桶中增加令牌，桶中的令牌不断增多。如果桶中令牌数已到达CBS/EBS，则丢弃多余令牌。

* 消耗令牌：输入数据包会消耗桶中的令牌。在网络传输中，数据包的大小通常不一致。大的数据包相较于小的数据包消耗的令牌要多。

* 判断是否通过：输入数据包经过令牌桶后的结果包括输出的数据包和丢弃的数据包。当桶中的令牌数量可以满足数据包对令牌的需求，则将数据包输出，否则将其丢弃

####4、ETCD####

etcd 是一个分布式一致性k-v存储系统，可用于服务注册发现与共享配置

<a href="https://github.com/coreos/etcd">查看源码</a>


##rabbitmq##

<a href="github.com/streadway/amqp">rabbitmq-golang驱动</a>

##系统监控与告警##

https://github.com/bosun-monitor/bosun



模块说明

https://github.com/golang/go/wiki/Modules










