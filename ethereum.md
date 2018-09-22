http://ethdoc.cn/introduction/index.html

http://blog.luoyuanhang.com/2018/05/02/eth-basis-block-concepts/ 很好的释疑的文章

#### 分叉链

分叉链指的是基于同一个软件的，增减一些功能后单独部署的一个区块链。

分叉链是最简单的技术，开发团只需要修改一些参数与功能，就可以发布出一个新的项目。

#### 侧链

二是侧链指的是与主链相平行的单独一个区块链，但是它和主链之间可以通过相互了解的协议互联。作为主链的补充，侧链可以提供一些主链不能提供的功能。但是这个互联对共识机制有要求，而且侧链必须有与主链相当的算力才能保证侧链货币的安全性。

侧链是借助于主链的接口，研发一个适合自身技术要求或自身功能的区块链项目。由侧链必须拥有与主链相当的算力才能保证侧链货币的安全性，可知其算力要求比较高，安全性有隐患，未来的投入比较大。

#### 子链

子链指的是在主链的平台来派生出来的具有其他功能的区块链。这些子链不能单独存在，必须通过主链提供的基础设施才能运行，并且免费获得主链的全部用户。

它是基于用户需求不同，而派发出来的区块链，但是它又不能独立存在于主链之外，必须基于主链才能运行，同时也可以获得主链的全部用户，以降低其宣传难度，提高用户量。

一是各个子链之间拥有灵活的交互功能。一个子链可以使用另外一个子链提供的资源（比如分布式文件系统），也就是说我要实现某一个功能，我自己设计的子链不具备这个功能，但是通过子链的交互，我最终还是能够实现这个功能。

源码阅读

p2p/discover/udp.go

pingPacket = iota + 1 // zero is 'reserved' //发送包
pongPacket//（ping的响应包）
findnodePacket //查找节点包请求包
neighborsPacket//邻居节点包（查找节点包响应包）



```
func (t *udp) sendPing(toid NodeID, toaddr *net.UDPAddr, callback func()) <-chan error {
	req := &ping{
		Version:    4, //这里写死了
		From:       t.ourEndpoint,
		To:         makeEndpoint(toaddr, 0), // TODO: maybe use known TCP port from DB
		Expiration: uint64(time.Now().Add(expiration).Unix()),
	}
	packet, hash, err := encodePacket(t.priv, pingPacket, req)
	if err != nil {
		errc := make(chan error, 1)
		errc <- err
		return errc
	}
	errc := t.pending(toid, pongPacket, func(p interface{}) bool {
		ok := bytes.Equal(p.(*pong).ReplyTok, hash)
		if ok && callback != nil {
			callback()
		}
		return ok
	})
	t.write(toaddr, req.name(), packet)
	return errc
}
```



以太坊 P2P 网络是一个完全加密的网络，提供 UDP 和 TCP 两种连接方式，

主网默认 TCP 通信端口是 30303，

推荐的 UDP 发现端口为 30301。

节点发现

节点发现是任何区块链节点接入区块链 P2P 网络的第一步。 这与你孤身一人去陌生地方旅游一样，如果没有地图和导航，那你只能拽附近的人问路，“拽附近的人问路”的这个动作就可以理解成节点发现。

节点发现可分为初始节点发现，和启动后节点发现。初始节点发现就是说你的全节点是刚下载的，第一次运行，什么节点数据都没有。启动后发现表示正在运行的钱包已经能跟随网络动态维护可用节点



##### NodeTable类负责以太坊的节点发现，NodeTable采用kademlia（KAD）算法进行节点发现

- NodeTable维护一个网络节点列表，此列表为当前可用节点，供上层使用
- 由于NodeID经过sha3生成出的Hash为256位。列表有256-1=255项，其中-1是因为刨除了当前节点（本机）
- 列表的每一项为一个节点桶（NodeBucket）,每个桶中最多放16个节点
- 列表的第i项代表据当前节点（本机）距离为i+1的网络节点集合



### **节点探索算法：**

#### 其中节点间距离定义如下：

节点距离定义为此XOR值的1位最高位的位数。

（位是1的最高位，比如0010 0000 1000 0101 它的最高位是14，一共是16位），那么它的距离就定义为14。

- 节点NodeID(512位)会先用sha3算法生成一个256位hash。算两个节点的256位hash的XOR值，节点距离定义为此XOR值的1位最高位的位数。（例如：0010 0000 1000 0101 位XOR值的化 那么这两个节点的距离为14）
- 此处的NodeID为网络节点公钥（512位）
- **注意：**这里的节点距离与机器的物理距离无关，这个距离仅仅是逻辑上的一种约定





### **节点的状态：**

- ##### Pending：

  - 挂起状态，每个新发现的节点或通过代码添加的节点的初始状态
  - 在新增节点时会向此节点发送Ping消息，已查看是否在线

- **Alive：**此状态说明Pong消息已收到，此节点在线。

- **Evicted：**由于对当前节点某个距离的桶最多只允许存在16个节点，若在此距离发现的新节点正好超过了限额，则新节点保留，桶中最老的节点会被唤出，进入此状态

在go-ethereum中

在eth客户端启动时会添加5个种子节点，这些节点的NodeID、ip、端口被硬编码在params/bootnodes.go

```
package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var MainnetBootnodes = []string{
	// Ethereum Foundation Go Bootnodes
	"enode://a979fb575495b8d6db44f750317d0f4622bf4c2aa3365d6af7c284339968eef29b69ad0dce72a4d8db5ebb4968de0e3bec910127f134779fbcb0cb6d3331163c@52.16.188.185:30303", // IE
	"enode://3f1d12044546b76342d59d4a05532c14b85aa669704bfe1f864fe079415aa2c02d743e03218e57a33fb94523adb54032871a6c51b2cc5514cb7c7e35b3ed0a99@13.93.211.84:30303",  // US-WEST
	"enode://78de8a0916848093c73790ead81d1928bec737d565119932b98c6b100d944b7a95e94f847f689fc723399d2e31129d182f7ef3863f2b4c820abbf3ab2722344d@191.235.84.50:30303", // BR
	"enode://158f8aab45f6d19c6cbf4a089c2670541a8da11978a2f90dbf6a502a4a3bab80d288afdbeb7ec0ef6d92de563767f3b1ea9e8e334ca711e9f8e2df5a0385e8e6@13.75.154.138:30303", // AU
	"enode://1118980bf48b0a3640bdba04e0fe78b1add18e1cd99bf22d53daac1fd9972ad650df52176e7c7d89d1114cfef2bc23a2959aa54998a46afcf7d91809f0855082@52.74.57.123:30303",  // SG

	// Ethereum Foundation C++ Bootnodes
	"enode://979b7fa28feeb35a4741660a16076f1943202cb72b6af70d327f053e248bab9ba81760f39d0701ef1d8f89cc1fbd2cacba0710a12cd5314d5e0c9021aa3637f9@5.1.83.226:30303", // DE
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
	"enode://30b7ab30a01c124a6cceca36863ece12c4f5fa68e3ba9b0b51407ccc002eeed3b3102d20a88f1c1d3c3154e2449317b8ef95090e77b312d5cc39354f86d5d606@52.176.7.10:30303",    // US-Azure geth
	"enode://865a63255b3bb68023b6bffd5095118fcc13e79dcf014fe4e47e065c350c7cc72af2e53eff895f11ba1bbb6a2b33271c1116ee870f266618eadfc2e78aa7349c@52.176.100.77:30303",  // US-Azure parity
	"enode://6332792c4a00e3e4ee0926ed89e0d27ef985424d97b6a45bf0f23e51f0dcb5e66b875777506458aea7af6f9e4ffb69f43f3778ee73c81ed9d34c51c4b16b0b0f@52.232.243.152:30303", // Parity
	"enode://94c15d1b9e2fe7ce56e458b9a3b672ef11894ddedd0c6f247e0f1d3487f52b66208fb4aeb8179fce6e3a749ea93ed147c37976d67af557508d199d9594c35f09@192.81.208.223:30303", // @gpip
}

// RinkebyBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network.
var RinkebyBootnodes = []string{
	"enode://a24ac7c5484ef4ed0c5eb2d36620ba4e4aa13b8c84684e1b4aab0cebea2ae45cb4d375b77eab56516d34bfbd3c1a833fc51296ff084b770b94fb9028c4d25ccf@52.169.42.101:30303", // IE
	"enode://343149e4feefa15d882d9fe4ac7d88f885bd05ebb735e547f12e12080a9fa07c8014ca6fd7f373123488102fe5e34111f8509cf0b7de3f5b44339c9f25e87cb8@52.3.158.184:30303",  // INFURA
	"enode://b6b28890b006743680c52e64e0d16db57f28124885595fa03a562be1d2bf0f3a1da297d56b13da25fb992888fd556d4c1a27b1f39d531bde7de1921c90061cc6@159.89.28.211:30303", // AKASHA
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	"enode://06051a5573c81934c9554ef2898eb13b33a34b94cf36b202b69fde139ca17a85051979867720d4bdae4323d4943ddf9aeeb6643633aa656e0be843659795007a@35.177.226.168:30303",
	"enode://0cc5f5ffb5d9098c8b8c62325f3797f56509bff942704687b6530992ac706e2cb946b90a34f1f19548cd3c7baccbcaea354531e5983c7d1bc0dee16ce4b6440b@40.118.3.223:30304",
	"enode://1c7a64d76c0334b0418c004af2f67c50e36a3be60b5e4790bdac0439d21603469a85fad36f2473c9a80eb043ae60936df905fa28f1ff614c3e5dc34f15dcd2dc@40.118.3.223:30306",
	"enode://85c85d7143ae8bb96924f2b54f1b3e70d8c4d367af305325d30a61385a432f247d2c75c45c6b4a60335060d072d7f5b35dd1d4c45f76941f62a4f83b6e75daaf@40.118.3.223:30307",
}
```



### **节点发现协议：**

- 协议：
  - FindNode：节点查询协议，向目标节点询问其临近节点列表
  - Neighbours：响应FindNode消息，当某节点接到其它节点发来的FindNode消息时，会回送Neighbours消息，其中携带了此节点的附近节点
  - Ping：用来查看节点是否存活。对于缺失NodeID的节点，也可用来询问其NodeID
  - Pong：对Ping消息的响应

- - 当前节点会向随机选定的节点的附近节点集合中每个节点发送FindNode消息，表示希望查这些节点的附近节点（A对B说，把你的邻居告诉我一下）
  - 这些节点在接收到FindNode消息会需要向发送节点回送Neighbours消息，并在此消息内含有所在节点的附近节点集合（B对A说，我的邻居是...）
  - 当前节点在收到回送的Neighbours消息后，会将Neighbours中所携带的节点加入到自己的网络节点列表中，并对这些携带节点发送Ping消息
  - 等到Pong的消息到达，证明此节点存活，并加入到节点列表中

geth/config.go->utils.SetNodeConfig(ctx, &cfg.Node)->SetP2PConfig->

        setNodeKey(ctx, cfg)
    
        setNAT(ctx, cfg)
    setListenAddress(ctx, cfg)
    setBootstrapNodes(ctx, cfg)
    setBootstrapNodesV5(ctx, cfg)



命令行支持的参数基本都定义在cmd/utils/flags.go文件中

单个节点的默认配置在node/defaults.go文件中

cmd/geth/main.go

geth-->makeFullNode()

geth-->startNode()

geth-->node.wait()



node/node.go

node.start()会调用p2s.Server{},并且将p2s.Server.Start()



ethdb数据库相关操作，利用leveldb来进行数据存储



在miner/worker.go中挖出新块，然后发送通知（chan）,然后在

eth/handler.go中，进行通知给peers

```
// Mined broadcast loop
//广播新区块
func (pm *ProtocolManager) minedBroadcastLoop() {
	// automatically stops if unsubscribe
	for obj := range pm.minedBlockSub.Chan() {
		if ev, ok := obj.Data.(core.NewMinedBlockEvent); ok {
			pm.BroadcastBlock(ev.Block, true)  // First propagate block to peers
			pm.BroadcastBlock(ev.Block, false) // Only then announce to the rest
		}
	}
}
```

```
//广播交易
func (pm *ProtocolManager) txBroadcastLoop() {
	for {
		select {
		case event := <-pm.txsCh:
			pm.BroadcastTxs(event.Txs)

		// Err() channel will be closed when unsubscribing.
		case <-pm.txsSub.Err():
			return
		}
	}
}
```



peer的分类

```
static peer
```



源码阅读部分

eth/protocol.go

1、定义了使用的协议，比如获取区块头消息，获取区块消息等等。并定义了部分消息的结构体。

2、定义了交易池接口txPool



## 钥匙文件

每个账户都由一对钥匙定义，一个私钥和一个公钥。 账户以地址为索引，地址由公钥衍生而来，取公钥的最后 20个字节。每对私钥 /地址都编码在一个钥匙文件里。钥匙文件是JSON文本文件，可以用任何文本编辑器打开和浏览。钥匙文件的关键部分，账户私钥，通常用你创建帐户时设置的密码（也叫口令）进行加密(简单的说，就是用你的密码对私钥进行加密)。钥匙文件可以在以太坊节点数据目录的keystore子目录下找到。确保经常给钥匙文件备份！查看备份和恢复账号章节了解更多。创建钥匙和创建帐户是一样的。

1. 不必告诉任何人你的操作。
2. 不必和区块链同步。
3. 不必运行客户端。
4. 甚至不必连接到网络。

当然新账户不包含任何以太币。但它将会是你的，你大可放心，没有你的钥匙和密码，没有人能进入。

转换整个目录或任何以太坊节点之间的个人钥匙文件都是安全的。

所以务必要注意两个重点：

1、钥匙文件

2、你的密码（你的口令）

这两个东西缺一不可。

为了从账号发送交易，包括发送以太币，你必须同时有钥匙文件和密码。确保钥匙文件有个备份并牢记密码，尽可能安全地存储它们。这里没有逃亡路径，如果钥匙文件丢失或忘记密码，就会丢失所有的以太币。没有密码不可能进入账号，也没有忘记密码选项。所以一定不要忘记密码。

# 一、什么是难度

难度(Difficulty)一词来源于区块链技术的先驱比特币，用来度量挖出一个区块平均需要的运算次数。

难度(Difficulty)通过控制合格的解在空间中的数量来控制平均求解所需要尝试的次数，也就可以间接的控制产生一个区块需要的时间，这样就可以使区块以一个合理而稳定的速度产生。

当挖矿的人很多，单位时间能够尝试更多次时，难度就会增大，当挖矿的人减少，单位时间能够尝试的次数变少时，难度就降低。这样产生一个区块需要的时间就可以做到稳定。

什么是难度 

难度(Difficulty)一词来源于区块链技术的先驱比特币，用来度量挖出一个区块平均需要的运算次数。挖矿本质上就是在求解一个谜题，不同的电子币设置了不同的谜题。比如比特币使用SHA-256、莱特币使用Scrypt、以太坊使用Ethash。一个谜题的解的所有可能取值被称为解的空间，挖矿就是在这些可能的取值中寻找一个解。

这些谜题都有如下共同的特点：

1、没有比穷举法更有效的求解方法 

2、解在空间中均匀分布（这个很重要），从而使每一次穷举尝试找到一个解的概率基本一致

3、解的空间足够大，保证一定能够找到解 

假设现在有一种电子币，解所在的空间为0-99共100个数字，谜题为x<100。这个谜题非常简单，空间中的任何一个数字都能满足。如果想让谜题更难以求解该怎么做呢？把谜题改成x<50，现在空间中只有一半的数字能满足了，也就是说，现在的难度比原来大了。并且我们还能知道难度大了多少，原来求解平均要尝试1次，现在求解平均要尝试2次了，也就是说，x<50的难度是x<100的2/1=2倍。同理，如果谜题变成x<10，难度就是x<100的100/10=10倍。

关键方法CalcDifficulty在consensus/ethash/consensus.go中



以太坊的创始人把以太坊发展分为4个里程碑阶段：前沿（Frontier）、家园（Homestead）、大都会（Metropolis）和宁静（Serenity）。以太坊目前正处于第3阶段：大都会之拜占庭阶段。自从前沿阶段以来，“难度炸弹”（增加ETH开采难度的协议）就被编入了以太坊[区块链](https://www.jutuilian.com/)。

以太坊的最后一个里程碑（宁静阶段）将带来一个重大变化：以太坊的区块链共识算法将从工作量证明（PoW）变为权益证明（PoS）。也就是说，在以太坊网络能够从PoW切换到PoS之前，必须将矿工从PoW区块链切换到PoS区块链。



在PoW（proof of work，工作量证明）机制中，计算机进行算法的竞赛。首先解出答案并将新区块广播到网络的计算机获得新的加密货币奖励和区块中的交易费用。由于先算出答案的计算机将获得奖励，因此矿工有动力尽可能多地使用算力，以便获得区块奖励。然而，为了获得更多的哈希算力，矿工们需要花费更多的资源并支付更多的电力成本来运行采矿设备，这就意味着系统需要消耗大量的算力和电力。 



PoS(proof of stake,权益证明)：在这个机制里，起作用的是验证者而不是矿工，它的原理是：作为验证节点，首先你必须拥有一定数量的ETH，根据ETH的数量和时间会产生用于下注验证区块的权益。只有拥有权益的节点才能有效验证区块，当你验证的区块被打包进区块链，你将获得和所拥有的权益成正比的区块奖励。如果你验证恶意或错误的区块，那么你所下注的权益将被扣除。由于支持无效区块的行为将受到严惩，因此，相比于PoW机制，个人进行诚实行为的动机更强。

**无论是PoS或PoW，富人都会聚集更多财富。而且，随着每个新区块被创建出来，财富的差距将越拉越大**



PoS机制的实施将会吸引更多分布式节点的加入，为各种分布式应用的运行打下物理基础，以太坊将有希望成为去中心化领域的app商店，互联网的新时代也将到来



Frontier是2015年7月以太坊发行初期的试验阶段，那个时候的软件还不太成熟，但是可以进行基本的挖矿，学习，试验。系统运行之后，吸引了更多的人关注并参与到开发中来，以太坊作为一个应用平台，需要更多的人去开发自己的去中心化应用来实现以太坊本身的价值。随着人气渐旺，以太坊的价值也水涨船高。 Homestead是以太坊第一个正式的产品发行版本，于2016年3月发布。100%采用PoW挖矿，但是挖矿的难度除了因为算力增长而增加之外，还有一个额外的难度因子呈指数级增加，这就是难度炸弹（Difficulty Bomb）。由于PoS的运用将会降低挖矿的门槛，因为你不需要再去购买价格高昂的硬件矿机，只需要购买一定数量的ETH，将其作为保证金通过权益证明的方式验证交易有效性，即可拿到一定的奖励。因此，对矿工来说他们花高价购买的矿机将无用武之地，这势必会引起矿工的不满。为了防止PoW转PoS的过程中矿工联合起来抵制，从而分叉出两条以太坊区块链，难度炸弹被引入。难度炸弹指的是计算难度时除了根据出块时间和上一个区块难度进行调整外，加上了一个每十万个区块呈指数型增长的难度因子。 Homestead的下一阶段Metropolis又被分成了两个阶段：Byzantium和Constantinople。目前以太坊运行在Byzantium阶段。Constantinople的规划与开发预计将在今年晚些时候进行。 至于以太坊的最后一个阶段Serenity，即转成PoS的软件版本至少还要等一两年了。



区块难度不能低于以太坊的创世区块的难度，这是以太坊难度的下限。



挖矿的本质就是解决一个数学计算，谁先算出来谁就获得奖励（币），这个数学计算方式也很简单，就是一直不断的尝试碰撞结果。

矿机1秒内能计算的hash算法次数越多算力越大，挖的币越多。 

算力------每秒矿机的整体hash算法运算次数





我们在这里记录一下存储在数据库中的数据，我们约定一套规则，用来表述我们的意思,注意，这里只是抽象来看，不能看成是实际就这么存，这里只是为了体现kv的关系。

（hash,number） 这个表示说hash和number以及一些可唯一识别的数据拼接起来形成key

比如说(number,hash)=>total_difficulties和(number,hash)=>header，对=>total_difficulties可能是加了td形成key,而对=>header则可能是加了h形成key

（hash,number）=>header,则表示说在数据库总保存着这样的映射关系，就是根据hash和number我可以取到header,存的时候也是根据hash和number来保存header,在以太坊中的存储的一些数据。

- hash=>number
- number=>hash
- (number,hash)=>total_difficulties
- (number,hash)=>header
- "LastHeader"=>hash #"LastHeader"注意这里的表示实际的













我们列举一下core/headerchain.go做了哪些事情

1. ```go
   func (hc *HeaderChain) GetHeaderByNumber(number uint64) *types.Header
   
   func (hc *HeaderChain) HasHeader(hash common.Hash, number uint64) bool
   
   func (hc *HeaderChain) GetHeaderByHash(hash common.Hash) *types.Header
   
   func (hc *HeaderChain) GetHeader(hash common.Hash, number uint64) *types.Header
   
   func (hc *HeaderChain) WriteTd(hash common.Hash, number uint64, td *big.Int) error
   
   func (hc *HeaderChain) GetTdByHash(hash common.Hash) *big.Int 
   
   func (hc *HeaderChain) GetTd(hash common.Hash, number uint64) *big.Int
   
   func (hc *HeaderChain) GetAncestor(hash common.Hash, number, ancestor uint64, maxNonCanonical *uint64) (common.Hash, uint64)
   
   func (hc *HeaderChain) GetBlockHashesFromHash(hash common.Hash, max uint64) []common.Hash
   
   func (hc *HeaderChain) InsertHeaderChain(chain []*types.Header, writeHeader WhCallback, start time.Time) (int, error)
   
   func (hc *HeaderChain) ValidateHeaderChain(chain []*types.Header, checkFreq int) (int, error)
   
   
   func (hc *HeaderChain) WriteHeader(header *types.Header) (status WriteStatus, err error)
   
   func (hc *HeaderChain) GetBlockNumber(hash common.Hash) *uint64
   
   
   ```



布隆过滤器

布隆过滤器是一种基于Hash的高效查找结构，能够快速（常数时间内）回答“某个元素是否在一个集合内”的问题。布隆过滤器因为其高效性大量应用于网络和安全领域。

Hash可以将任意内容映射到一个固定长度的字符串，而且不同内容映射到相同串的概率很低。布隆过滤器采用了多个Hash函数来提高空间利用率。对同一个给定输入来说，多个Hash函数计算出多个地址，分别在位串的这些地址上标记为1，进行查找时，进行同样的计算过程，并查看对应元素，如果都为1，则说明较大概率是存在该输入。

































