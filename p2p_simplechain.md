```
一个参数
--netrestrict 如果设置了该参数，除了种子节点，目前通过查询邻居节点时，返回的节点只有这个参数中指定的ip才可以进入桶中

目前发现一个问题，主动ping过来的节点我们并没有做这个限制
```

https://eips.ethereum.org/EIPS/eip-778

一些组件：

* 一个私钥
* 一个网络白名单
* 引导节点



一些限制

* udp的端口必须大于1024
* 每个桶中相同的前缀的ip地址最多不超过2个
* 整个地址表中相同的前缀的ip地址最多不超过10个
* 也就是说我们一旦验证它是存活的，我们就把它放在桶的列表的最前面，我们定时从表中随机选取一个桶，然后从桶的列表的尾部取一个节点去验证，一旦验证它是存活的，我们就把它放在桶的列表的最前面，它若是挂了，那么我们就从替换列表中选取一个节点，替换掉该桶的尾部的那个节点。



将node加入桶的过程。

* 首先根据id判断，如果是自身的话，肯定不能加入桶
* 根据id查找它可以加入哪个桶
* 根据它可以加入哪个桶，查询这个对应的桶是否已经存在该节点，如果已经存在，那么就直接返回（已在，不需要再加入）
* 如果不在桶中，那么就看一下桶是否满了，如果已经满了，那么我们就暂且添加到列表中（如果已经存在，也直接返回）
* 是否是受限ip







```
可以把基于每类的IP网络进一步分成更小的网络，每个子网由路由器界定并分配一个新的子网网络地址,子网地址是借用基于每类的网络地址的主机部分创建的。划分子网后，通过使用掩码，把子网隐藏起来，使得从外部看网络没有变化，这就是子网掩码。
```

## IPv4

IP地址由四段组成，每一个段是一个字节，一个字节共8位，最大值是255

IP地址由两部分组成，即网络地址和主机地址。网络地址表示其属于互联网的哪一个网络，主机地址表示其属于该网络中的哪一台主机。二者是主从关系。

#### 规定

* A类地址的第一个字节的第一位为0，使用前一个字节来表示网络地址，后三个字节来表示主机地址
* B类地址的第一个字节的第一位为1，第二位为0，使用前两个字节来表示网络地址，后两个字节来表示主机地址。
* C类地址的第一个字节的第一位为1，第二位为1，第三位为0，使用前三个字节来表示网络地址，后一个字节来表示主机地址。
* D类地址，1110开头(多播组号)
* E类地址，111110（留待后用）

一般我们也就是能够使用A、B、C类地址,全0和全1的都保留不用。

看ip地址的时候我们要注意段和位的说法，段是指用用点分十进制来表示的时候，例如1.0.0.0,它表示的是位

00000001 00000000 00000000 00000000

我们要注意的是32位，4个字节。

A类：(1.0.0.0-126.0.0.0)（默认子网掩码：255.0.0.0或 0xFF000000）第一个字节为网络号，后三个字节为主机号。该类IP地址的最前面为“0”，所以地址的网络号取值于1~126之间。一般用于大型网络。

B类：(128.0.0.0-191.255.0.0)（默认子网掩码：255.255.0.0或0xFFFF0000）前两个字节为网络号，后两个字节为主机号。该类IP地址的最前面为“10”，所以地址的网络号取值于128~191之间。一般用于中等规模网络。

C类：(192.0.0.0-223.255.255.0)（子网掩码：255.255.255.0或 0xFFFFFF00）前三个字节为网络号，最后一个字节为主机号。该类IP地址的最前面为“110”，所以地址的网络号取值于192~223之间。一般用于小型网络。

D类：是多播地址。该类IP地址的最前面为“1110”，所以地址的网络号取值于224~239之间。一般用于多路广播用户[1] 。

E类：是保留地址。该类IP地址的最前面为“1111”，所以地址的网络号取值于240~255之间。

在IP地址3种主要类型里，各保留了3个区域作为私有地址，其地址范围如下： 
A类地址：10.0.0.0～10.255.255.255 
B类地址：172.16.0.0～172.31.255.255 
C类地址：192.168.0.0～192.168.255.255

### 子网掩码

当我们对一个网络进行子网划分时，基本上就是将它分成小的网络

**子网掩码是一个32位的2进制数，其对应网络地址的所有位置都为1，对应于主机地址的所有位置都为0。**

将IP地址和子网掩码都换算成二进制,然后进行**与运算**,结果就是**网络地址**.

```




```



节点发现的数据包设计

```
hash+sign+data
hash=Keccak256(sign+data)
sign=signMethod(data,privateKey)
data[0]表示消息的类型
data[1:]是经过了rlp编码的数据
```

 #### 消息类型：

* p_pingV4=1
* p_pongV4= 2
* p_findnodeV4 = 3
* p_neighborsV4=4
* p_enrRequestV4=5
* p_enrResponseV4=6





pong 消息（server-client）

* 将client的真实地址（对外地址）返回给client
* 将client的ping消息中的hash返回给client
* 将超时时间返回给client
* 将server的序列号返回给client



测试框架

```
https://github.com/stretchr/testify/
```



目前来看

discv5使用的是内存数据库库，并没有存储文件。也就是没有进行持久化。

然后我们需要注意的地方是

```
//启用节点发现功能
	if !srv.NoDiscovery {
		//如果启用了v5节点发现功能
		if srv.DiscoveryV5 {
		  //要特别注意这里
			unhandled = make(chan discover.ReadPacket, 100)
			//conn 在sconn中主要是用来发送消息的，读取消息是从unhandled中读取的。
			sconn = &sharedUDPConn{conn, unhandled}
		}
		cfg := discover.Config{
			PrivateKey:  srv.PrivateKey,
			NetRestrict: srv.NetRestrict,
			Bootnodes:   srv.BootstrapNodes,
			Unhandled:   unhandled,
			Log:         srv.log,
		}
		//节点发现v4
		ntab, err := discover.ListenUDP(conn, srv.localnode, cfg)
		if err != nil {
			return err
		}
		srv.ntab = ntab
		//决定选择那些节点来连接有这个方法决定的
		srv.discmix.AddSource(ntab.RandomNodes())
		srv.staticNodeResolver = ntab
	}

	// Discovery V5
	if srv.DiscoveryV5 {
		var ntab *discv5.Network
		var err error
		if sconn != nil {
			//从nodeDBPath为"",可得知此时它使用的是内存数据库
			ntab, err = discv5.ListenUDP(srv.PrivateKey, sconn, "", srv.NetRestrict)
		} else {
			//从nodeDBPath为"",可得知此时它使用的是内存数据库
			ntab, err = discv5.ListenUDP(srv.PrivateKey, conn, "", srv.NetRestrict)
		}
		if err != nil {
			return err
		}
		if err := ntab.SetFallbackNodes(srv.BootstrapNodesV5); err != nil {
			return err
		}
		srv.DiscV5 = ntab
	}
```



然后我们再把注意力放在这

```
// sharedUDPConn implements a shared connection. Write sends messages to the underlying connection while read returns
// messages that were found unprocessable and sent to the unhandled channel by the primary listener.
type sharedUDPConn struct {
	*net.UDPConn
	unhandled chan discover.ReadPacket
}

// ReadFromUDP implements discv5.conn
func (s *sharedUDPConn) ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error) {
	packet, ok := <-s.unhandled
	if !ok {
		return 0, nil, errors.New("connection was closed")
	}
	l := len(packet.Data)
	if l > len(b) {
		l = len(b)
	}
	copy(b[:l], packet.Data[:l])
	return l, packet.Addr, nil
}
```



从上面我们可知，这里的ReadFromUDP其实不是原生的，是包装出来的

其实它的数据是从

```
// 处理udp数据包
// readLoop runs in its own goroutine. it handles incoming UDP packets.
func (t *UDPv4) readLoop(unhandled chan<- ReadPacket) {
	defer t.wg.Done()
	if unhandled != nil {
		defer close(unhandled)
	}
	buf := make([]byte, maxPacketSize)
	for {
		log.Error("UDPv4 readLoop")
		nbytes, from, err := t.conn.ReadFromUDP(buf)
		if netutil.IsTemporaryError(err) {
			// Ignore temporary read errors.
			t.log.Debug("Temporary UDP read error", "err", err)
			continue
		} else if err != nil {
			// Shut down the loop for permament errors.
			if err != io.EOF {
				t.log.Debug("UDP read error", "err", err)
			}
			return
		}
		//从这里我们看到，这是同步在处理，并不是异步处理
		if t.handlePacket(from, buf[:nbytes]) != nil && unhandled != nil {
			select {
			//统一处理我们不能处理的错误,这里我们把数据包交给discv5来处理，有可能它能处理。
			case unhandled <- ReadPacket{buf[:nbytes], from}:
			default:
			}
		}
	}
}
```





小结：

在discv4的时候，一个数据包是这样子的

```
hash+sign+msgType+msg
headSize=hashSize+signSize (目前的配置就是97=32+65)
```

而在discv5一个数据包是这样子的

```
versionPrefix+sign+msgType+msg
headSize=versionPrefixSize+signSize(目前的配置就是87=22+65)
```

也就是说默认，它采取的方式是从conn中读取出一个数据包(udp包都是整包的，它是有边界的，与tcp区别)以后，先使用discv4的解码方式来解码，如果能够解码，那么就由discv4来解决，否则就把数据放入unhandled通道中，将数据交由discv5去处理。



我们从les模块中看它的调用

```
if srvr.DiscV5 != nil {
		for _, topic := range s.lesTopics {
			topic := topic
			go func() {
				logger := log.New("topic", topic)
				logger.Info("Starting topic registration")
				defer logger.Info("Terminated topic registration")

				srvr.DiscV5.RegisterTopic(topic, s.closeCh)
			}()
		}
	}
```

注册多个主题时，就是这么操作的。当



```
addTopic
removeRegisterTopic
nextRegisterLookup
registerLookupDone
ticketRegistered
```



我们看主题查找的模块调用

```
// discoverNodes wraps SearchTopic, converting result nodes to enode.Node.
func (pool *serverPool) discoverNodes() {
	ch := make(chan *discv5.Node)
	go func() {
		pool.server.DiscV5.SearchTopic(pool.topic, pool.discSetPeriod, ch, pool.discLookups)
		close(ch)
	}()
	for n := range ch {
		pubkey, err := decodePubkey64(n.ID[:])
		if err != nil {
			continue
		}
		pool.discNodes <- enode.NewV4(pubkey, n.IP, int(n.TCP), int(n.UDP))
	}
}



if pool.server.DiscV5 != nil {
		pool.discSetPeriod = make(chan time.Duration, 1)
		pool.discNodes = make(chan *enode.Node, 100)
		pool.discLookups = make(chan bool, 100)
		go pool.discoverNodes()
	}
	
	
	
//主题查找
func (net *Network) SearchTopic(topic Topic, setPeriod <-chan time.Duration, found chan<- *Node, lookup chan<- bool) {
	for {
		select {
		case <-net.closed:
			return
		case delay, ok := <-setPeriod:
			select {
			case net.topicSearchReq <- topicSearchReq{topic: topic, found: found, lookup: lookup, delay: delay}:
			case <-net.closed:
				return
			}
			if !ok {
				return
			}
		}
	}
}
```



要特别注意这些链接内容的更改

```
https://github.com/ethereum/devp2p/blob/master/discv5/discv5.md

https://github.com/ethereum/devp2p/blob/master/discv5/discv5-wire.md

https://github.com/ethereum/devp2p/blob/master/discv5/discv5-theory.md #看理论

https://github.com/ethereum/devp2p/blob/master/discv5/discv5-rationale.md
```





### 门票

广告应在队列中保持固定的时间，即`target-ad-lifetime`。为了维持这一保证，新注册受到限制，注册人必须等待一定时间才能被接纳。当节点尝试放置广告时，它会收到一张“票”，告知他们必须等待多长时间才能被接受。当等待时间过去时，由注册者节点保留票证并将其出示给广告媒体。

等待时间常数为：

```
target-ad-lifetime = 15min
```

根据以下规则确定为任何注册尝试分配的等待时间：

- 当表格已满时，将根据整个表格中最早的广告的生存时间来分配等待时间，即注册者必须等待表格位置可用。
- 当主题队列已满时，等待时间取决于队列中最早的广告的生命周期。`target-ad-lifetime - oldest-ad-lifetime`在这种情况下，分配的时间为。
- 否则，广告可能会立即放置。

票证是不透明的对象，存储由发行节点确定的任意信息。尽管编码和票证验证的详细信息取决于实现，但票证必须包含足够的信息以验证：

- 尝试使用票证的节点是请求票证的节点。
- 该票证仅对单个主题有效。
- 该票只能在注册窗口内使用。
- 该票不能多次使用。

理解：就是自己（本节点）要发布主题，那么远端节点我们就把它看成广告媒体，那么本节点就看成注册人，也就是本节点要等待一定的时间以后才能到远端节点去注册。比如本节点是A，远端节点是B,A想要发布主题，也即是说A想要在B出放置广告，那么发送ping的过程中，B会返回给A一张票据，并告知A要等待多长时间，B才会接受A的广告放置请求，把广告放置在B处。

发行票证后，保留票证的节点必须等待注册窗口打开。注册窗口的长度为10秒。通过注册窗口后，票证将失效。

等待时间  注册窗口（10秒）   票据失效