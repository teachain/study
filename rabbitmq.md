## RabbitMQ##

 RabbitMQ是一个由erlang开发的AMQP（Advanced Message Queue Protocol）的开源实现。
 
 <a href="http://www.rabbitmq.com">官方网站</a>
 
  对于一个数据从Producer到Consumer的正确传递，还有三个概念需要明确：exchanges, queues and bindings。

 * Exchanges are where producers publish their messages.

 * Queues are where the messages end up and are received by consumers

 * Bindings are how the messages get routed from the exchange to particular queues.

 
 Connection： 就是一个TCP的连接。Producer和Consumer都是通过TCP连接到RabbitMQ Server的。以后我们可以看到，程序的起始处就是建立这个TCP连接。

 Channels： 虚拟连接。它建立在上述的TCP连接中。数据流动都是在Channel中进行的。也就是说，一般情况是程序起始建立TCP连接，第二步就是建立这个Channel。
 
  默认情况下，如果Message 已经被某个Consumer正确的接收到了，那么该Message就会被从queue中移除。当然也可以让同一个Message发送到很多的Consumer。
    如果一个queue没被任何的Consumer Subscribe（订阅），那么，如果这个queue有数据到达，那么这个数据会被cache，不会被丢弃。当有Consumer时，这个数据会被立即发送到这个Consumer，这个数据被Consumer正确收到时，这个数据就被从queue中删除。


使用命令启动服务器

```
#start the server
#$install_rabbitmq_path表示rabbitmq所在根目录

//进入rabbitmq sbin目录
cd $install_rabbitmq_path/sbin

//执行可执行文件
./rabbitmq-server

#以守护进程的方式启动服务器
./rabbitmq-server -detached

```


```
#使用rabbitMQ管理页面
./rabbitmq-plugins enable rabbitmq_management
./rabbitmq-server
```

```
//打印un-acked Messages
sudo rabbitmqctl list_queues name messages_ready messages_unacknowledged  

```


producer(Producer只能发送到exchange,routing_key就是指定的queue名字-单发的时候。)

```
   //编程套路
   //1、创建一个tcp链接
   //2、创建一个channel
   //3、声明一个queue
   //4、publish一条message
   
```

consumer

```
   //编程套路
   //1、创建一个tcp链接
   //2、创建一个channel
   //3、声明一个queue
   //4、subscribe(订阅)
   //5、无限循环接收

```


##Message acknowledgment 消息确认##

* no-ack 每次Consumer接到数据后，而不管是否处理完成，RabbitMQ Server会立即把这个Message标记为完成，然后从queue中删除了。
* ack  在处理数据后发送的ack，就是告诉RabbitMQ数据已经被接收，处理完成，RabbitMQ可以去安全的删除它了。

<font color="red"> 为了保证数据不被丢失，RabbitMQ支持消息确认机制，即acknowledgments。为了保证数据能被正确处理而不仅仅是被Consumer收到，那么我们不能采用no-ack。而应该是在处理完数据后发送ack。</font>

 如果Consumer退出了但是没有发送ack，那么RabbitMQ就会把这个Message发送到下一个Consumer。这样就保证了在Consumer异常退出的情况下数据也不会丢失。
 
 默认情况下，消息确认是打开的（enabled）
 
##Message durability消息持久化##

<font color="red"> 为了保证在RabbitMQ退出或者crash了数据仍没有丢失，需要将queue和Message都要持久化。</font>

* queue的持久化需要在声明时指定durable=True：也就是queue_declare的时候指定。

* 需要持久化Message，即在Publish的时候指定一个properties


 为了数据不丢失，我们采用了：

* 在数据处理结束后发送ack，这样RabbitMQ Server会认为Message Deliver 成功。
* 持久化queue，可以防止RabbitMQ Server 重启或者crash引起的数据丢失。
* 持久化Message，可以防止RabbitMQ Server 重启或者crash引起的数据丢失。

但是还是不能够保证100%不丢失数据。一种可能的方案是在系统panic时或者异常重启时或者断电时，应该给各个应用留出时间去flush cache，保证每个应用都能exit gracefully。


##Fair dispatch 公平分发##
默认状态下，RabbitMQ将第n个Message分发给第n个Consumer。当然n是取余后的。它不管Consumer是否还有unacked Message，只是按照这个默认机制进行分发。

 通过 (channel的)basic.qos 方法设置prefetch_count=1 。这样RabbitMQ就会使得每个Consumer在同一个时间点最多处理一个Message。换句话说，在接收到该Consumer的ack前，他它不会将新的Message分发给它。
 
 <font color="red">注意，这种方法可能会导致queue满。当然，这种情况下你可能需要添加更多的Consumer，或者创建更多的virtualHost来细化你的设计。</font>
 

##分发到多Consumer（Publish/Subscribe）##

有三种类型的Exchange：direct, topic 和fanout。fanout就是广播模式，会将所有的Message都放到它所知道的queue中

通过exchange，而不是routing_key来publish Message

