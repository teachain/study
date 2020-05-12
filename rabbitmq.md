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

### 为什么RabbitMQ有Queue，还要有Exchange?

exchange就像一个导购一样

```
假设你在Apple商店里边，先要买耳机。 店里就会有人过来问你:"需要什么?" 你告诉他你需要买耳机，然后他就把你带到他的同事的柜台前的排队队列之后等待。因为很多其他人也在买东西，销售员正在处理队列前面的那个消费者。 如果这个时候，另外一个人进店了，刚才招呼你的人会同样询问对方需要什么帮助。刚进来的人需要修下手机，被找呼的人带到了另外一个修理手机的柜台等待了。

这个例子中问你需要什么的人就是exchange, 他会根据需要把你路由到恰当的队列中排队等待。在队列的后面有很多员工，也就是对应队列的worker, 或者消费者。一次处理一个请求，基于先进先出的原则。也可能会根据最先到的人做一个简单轮询。

如果店里没有导流的服务员，那么你就需要来回在每个柜台前来回问是否能帮到你，直到找到你需要办理业务的柜台后开始排队。

当然，导航苹果商店的工作不复杂，但在应用程序中，你可能有很多队列，服务不同类型的请求，基于路由和绑定具有交换路由消息的键来说非常有帮助。 发布者只需要关心添加正确的路由密匙，而消费者只需要关心用正确的绑定密匙创建正确的队列，就可以做到"我对这些消息感兴趣。"
```



exchange的类型

* Direct exchange
* Fanout exchange
* Topic exchange
* Headers exchange



### Direct exchange(直接交换)

***直接交换基于消息路由密钥将消息传递到队列。直接交换是消息单播路由的理想选择***，下面是它的工作原理：

- 队列使用路由密钥K绑定到交换机
- 当具有路由键R的新消息到达直接交换时，如果K = R，则交换会将其路由到队列

直接交换通常用于以轮循方式在多个工作程序（同一应用程序的实例）之间分配任务。这样做时，重要的是要了解，在AMQP 0-9-1中，消息在使用者之间进行负载均衡，而不是队列之间进行负载均衡。

一条消息只会被一个消费者处理，一条消息只会进入一个队列。

##Fanout exchange（扇出交换）

***扇出交换机将消息路由到与其绑定的所有队列，并且路由键将被忽略。***如果将N个队列绑定到扇出交换，则将新消息发布到该交换时，会将消息的副本传递到所有N个队列。扇出交换机非常适合消息的广播路由。

因为扇出交换将消息的副本发送到绑定到它的每个队列，所以它的用例非常相似：

- 大型多人在线（MMO）游戏可以将其用于排行榜更新或其他全球性事件
- 体育新闻网站可以使用扇出交换以近乎实时的方式向移动客户端分发得分更新
- 分布式系统可以广播各种状态和配置更新
- 群组聊天可以使用扇出交换在参与者之间分发消息（尽管AMQP没有内置的在线状态概念，因此XMPP可能是更好的选择）

一条消息只会被多个消费者处理，一条消息会进入多个队列。



### Topic exchange

主题根据消息路由键和用于将队列绑定到交换机的模式之间的匹配将消息路由到一个或多个队列。主题交换类型通常用于实现各种发布/订阅模式变体。主题交换通常用于消息的多播路由。

主题交流有非常广泛的用例集。每当问题涉及多个使用者/应用程序，这些使用者/应用程序有选择地选择他们希望接收的消息类型时，应考虑使用主题交换。

示例使用：

- 分发与特定地理位置有关的数据，例如销售点
- 由多个工作人员完成的后台任务处理，每个工作人员都可以处理特定的任务集
- 股票价格更新（以及其他种类的财务数据的更新）
- 涉及分类或标记的新闻更新（例如，仅针对特定运动或团队）
- 云中各种服务的编排
- 分布式体系结构/特定于操作系统的软件构建或打包，其中每个构建器只能处理一个体系结构或OS



### Headers exchange

标头交换旨在用于在多个属性上路由，这些属性比路由键更容易表示为消息标头。标头交换忽略路由键属性。相反，用于路由的属性取自headers属性。如果标头的值等于绑定时指定的值，则认为消息匹配。

可以使用多个标题进行匹配，将队列绑定到标题交换。在这种情况下，代理需要从应用程序开发人员那里获得另一条信息，即，它应该考虑具有任何匹配的标头的消息，还是全部匹配的消息？这就是“ x-match”绑定参数的作用。当“ x-match”参数设置为“ any”时，仅一个匹配的标头值就足够了。或者，将“ x-match”设置为“ all”要求所有值必须匹配。

标头交换可以看作是“类固醇的直接交换”。由于它们基于标头值进行路由，因此可以将它们用作直接交换，而路由密钥不必是字符串。例如，它可以是整数或哈希（字典）。



[队列](https://www.rabbitmq.com/queues.html)在AMQP 0-9-1模式非常类似于其他MESSAGE-和任务排队系统队列：它们存储由应用程序使用的消息。队列与交换共享一些属性，但也具有一些其他属性：

- 名称
- 持久（队列将在代理重新启动后幸存）
- 独占（仅由一个连接使用，并且该连接关闭时队列将被删除）
- 自动删除（至少有一个使用方的队列在最后一个使用方退订时被删除）
- 参数（可选；由插件和特定于代理的功能使用，例如消息TTL，队列长度限制等）

必须先声明队列，然后才能使用它。声明队列将导致它创建（如果尚不存在）。如果队列已经存在并且其属性与声明中的相同，则该声明将无效。当现有队列属性与声明中的属性不同时，将引发代码为406（PRECONDITION_FAILED）的通道级异常。

### [队列名称](https://www.rabbitmq.com/tutorials/amqp-concepts.html#queue-names)

应用程序可以选择队列名称，也可以要求代理为它们生成名称。队列名称最多可以包含255个字节的UTF-8字符。AMQP 0-9-1代理可以代表应用程序生成唯一的队列名称。要使用此功能，请传递一个空字符串作为队列名称参数。生成的名称将与队列声明响应一起返回给客户端。

队列名称以“ amq”开头。保留供经纪人内部使用。尝试使用违反此规则的名称声明队列将导致通道级异常，其应答代码为



注意：**队列和exchange进行绑定时，要特别注意队列名称和routekey的一对一关系，如果不是的话，后面绑定的routekey会被前面绑定的routekey覆盖掉。 **



