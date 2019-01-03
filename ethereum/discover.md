### **核心数据结构：**

Table负责以太坊的节点发现，Table采用kademlia（KAD）算法进行节点发现

- Table维护一个网络节点列表，此列表为当前可用节点，供上层使用
- 由于NodeID经过sha3生成出的Hash为256位。列表有256-1=255(v5版本，非v5版本是17个桶)项，其中-1是因为刨除了当前节点（本机）
- 列表的每一项为一个节点桶（bucket）,每个桶中最多放16个节点
- 列表的第i项代表据当前节点（本机）距离为i+1的网络节点集合（难道是因为i从0开始？）

#### 其中节点间距离定义如下：

- 节点NodeID(512位)会先用sha3算法生成一个256位hash。算两个节点的256位hash的XOR值，节点距离定义为此XOR值的1位最高位的位数。（例如：0010 0000 1000 0101 位XOR值的化 那么这两个节点的距离为14）
- 此处的NodeID为网络节点公钥（512位）
- **注意：**这里的节点距离与机器的物理距离无关，这个距离仅仅是逻辑上的一种约定



### **发现算法思路**

1.先随机一个目标节点的NodeID

2.在列表中以相对NodeID的“距离”为指标，由近及远查找此待连节点“附近”的节点。并将这些节点放入“附近”节点集合

3.向目标节点的“附近”节点集合中的每个节点发送FindNode消息

4.若在目标节点的”附近”没有搜到节点，则返回步骤1

5.否则等待600ms后跳转到步骤2

### **节点的状态：**

Pending：

- 挂起状态，每个新发现的节点或通过代码添加的节点的初始状态
- 在新增节点时会向此节点发送Ping消息，已查看是否在线

**Alive：**此状态说明Pong消息已收到，此节点在线。

**Evicted：**由于对当前节点某个距离的桶最多只允许存在16个节点，若在此距离发现的新节点正好超过了限额，则新节点保留，桶中最老的节点会被唤出，进入此状态

这是一个节点的状态转换过程。

发现一个新节点，这个时候这个节点我们定义它为pending状态，然后我们给这个节点发送一个ping消息，在收到它的回复信息pong之前的这个时间内，它的状态都是pending,一旦收到pong消息，那么它的状态就转换为Alive状态，当它从桶中被换出，它就进入Evicted状态。



### **如何冷启动：**

由于初始化时节点列表为空，所以不可能找到目标节点的所谓附近节点。这就需要一些初始种子节点进行连接。在eth客户端启动时会添加6个种子节点，这些节点的NodeID、ip、端口被硬编码在params/bootnodes.go中。

 



### **节点发现协议：**

- 协议：

  - FindNode：节点查询协议，向目标节点询问其临近节点列表
  - Neighbours：响应FindNode消息，当某节点接到其它节点发来的FindNode消息时，会回送Neighbours消息，其中携带了此节点的附近节点
  - Ping：用来查看节点是否存活。对于缺失NodeID的节点，也可用来询问其NodeID
  - Pong：对Ping消息的响应

- 在上面的算法描述中，当前节点会向随机选定的节点的附近节点集合中每个节点发送FindNode消息，表示希望查这些节点的附近节点

- 这些节点在接收到FindNode消息会需要向发送节点回送Neighbours消息，并在此消息内含有所在节点的附近节点集合

- 当前节点在收到回送的Neighbours消息后，会将Neighbours中所携带的节点加入到自己的网络节点列表中，并对这些携带节点发送Ping消息

- 等到Pong的消息到达，证明此节点存活，并加入到节点列表中



  最大的数据包大小为1280个字节



  针对节点

  收到对端发送过来的ping消息

  1、回复一个pong消息

  2、记录下这个节点的收到ping消息的时间

  3、判断这个节点上次的pong消息到现在是否已经超过一定的时间了，如果超过，那么发送一个ping消息，如果没有超过，那么我们直接把它加入table中。



v5

```
ping

pong

findnode

neighbors

findnodeHash

topicRegister

topicQuery

topicNodes
```

state变化

unknow--->handleping(),ping()-->verifywait



### discover

1、udp.go定义了消息码

```
	pingPacket = iota + 1

	pongPacket

	findnodePacket

	neighborsPacket
```

各条消息的结构体，以及发送消息的方法，解析消息的方法，处理消息的方法。

udp.go 里启动了启动了两个goroutine。

1. 其中一个goroutine负责读取udp数据包。


discv5 一共有以下8条消息

1. ping
2. pong
3. findnode
4. neighbors
5. findnodeHash
6. topicRegister
7. topicQuery
8. topicNodes

然后network就根据节点的state来决定消息由谁来处理

一共有以下这些状态对象：

1. unknown
2. verifyinit
3. verifywait
4. remoteverifywait
5. known
6. contested
7. unresponsive



模拟一下场景

1、一个从来没有发过消息的节点A，发送消息过来到本节点local，这个时候节点A我们标记它的状态(state)为

unknown,这个时候消息就交由unknown的handle来处理，这个时候

它处理ping消息以外，其他消息都不处理，只是返回它的下一个状态，还是unknown，

当处理的是ping消息的时候，它的状态就流转到verifywait状态（下文接着），我们这里先讨论它处理ping消息的

过程：

根据Node和ping.Topics获取ticket,注意到获取ticket的时候，它是根据节点来的，local从数据库中针对节点A取到

了一个lastIssuedTicket，也就是可以认为lastIssuedTicket是绑定在节点A身上的（也就是说如果有个节点B进

来，它也有一个独立的lastIssuedTicket），然后根据这个lastIssuedTicket给节点A发行一个ticket。

这个ticket里包含了主题们以及他们能够注册的时间，以及这个这个节点的ticket流水号。

回复一个pong消息，并向节点A发送一个ping消息，此时节点A的状态是verifywait，那么当pong消息返回时。我

们需要进入第二步

2、当pong消息返回时，节点A所在的状态为verifywait，这是它交由verifywait的handle来处理。

它处理完pong消息之后，将进入known状态，而known状态的enter方法不为空，所以将调用known.enter()





这里我们从首次启动开始讨论：

当系统启动以后，首先获取种子节点，这个时候他们的状态应该都是unkown,所以根据代码的逻辑，这个时候将他

们的状态过渡到verifyinit，然后verifyinit存在enter方法，所以它进入了enter方法，这个时候主要做的工作是，

给种子节点发送ping消息。这个时候根据代码的逻辑，那么当pong消息返回的时候，它将由verifyinit的handle方

法来处理。处理通过了handleKnownPong方法，并且进入remoteverifywait状态。这个怎么来理解呢

我们知道它的处理流程是这样子的：

1. 节点A发送一个ping消息给节点B
2. 节点B收到ping消息之后给节点A发送Pong消息
3. 节点B发送ping消息给节点A
4. 节点A收到ping消息以后，给节点B发送pong消息

我们从节点A这一端来看节点B的状态，一开始，B的状态为unkown,然后我们调用了transition之后节点B就进入了

verifyinit状态，接着节点A给节点B发送了ping消息，然后当节点B发送给A pong消息以后，数据包将由verifyinit

的handle方法来处理。处理通过了handleKnownPong方法，并且进入remoteverifywait状态。注意这个时候，节

点B会紧跟发送一个ping消息给节点A,对节点A而言节点B的状态是remoteverifywait状态，所以这个时候的ping消

息将由remoteverifywait的handle方法来处理，处理就是给节点B发送一条pong消息，并维持remoteverifywait

状态。









