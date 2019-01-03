### eth协议

1. GetBlockHeadersMsg 获取区块头消息
2. BlockHeadersMsg  区块头消息
3. GetBlockBodiesMsg 获取区块体消息
4. BlockBodiesMsg  区块体消息
5. GetNodeDataMsg（>=eth63）获取节点数据消息
6. NodeDataMsg（>=eth63）节点数据消息
7. GetReceiptsMsg 获取收据消息
8. ReceiptsMsg  收据消息
9. NewBlockHashesMsg 区块哈希消息（区块哈希+高度）---Fetcher处理
10. NewBlockMsg 新区块消息 （区块+总难度）---Fetcher处理
11. TxMsg 交易消息

从上面的消息来看，除了区块消息和节点数据消息以外，其余消息都是支持批量获取的。



TCP连接建立之后，先进行握手协议

```
type statusData struct {
   ProtocolVersion uint32  //当前协议的版本
   NetworkId       uint64  //所在的网络Id
   TD              *big.Int //节点的区块的总难度
   CurrentBlock    common.Hash //当前最新区块的哈希
   GenesisBlock    common.Hash  //创世区块的哈希
}
```

握手的时候是这样子做的，本节点发送一个statusData消息给远端，然后远端也回复一个statusData给本节点。

然后我们就可以根据ProtocolVersion，NetworkId，GenesisBlock来确认我们是不是在同一个以太坊网络中。

也即是说我们要重点判断这些字段是否和本地的相等，不相等就关闭该连接。（因为我们要我们的小团体，不允许

别的异常分子来到我们的网络中）。一旦确认这些（ProtocolVersion，NetworkId，GenesisBlock）是相同的，

我们就认为他和我们是一伙的，我们可以继续进行消息的交互了。这样，我们就知道了所有建立了TCP连接的节点

的总难度和当前最新区块的哈希。



收到交易的处理：

交易来自于节点p，先做一个标记，记录节点p已经知晓该批交易。

将交易加入交易池中。pm.acceptTxs==1时，才可以接收交易。



收到新区块的处理：

区块来自于节点p，先做一个标记，记录节点p已经知晓该区块。

接着把该区块消息放入fetcher的队列中。

然后拿当前获取的新区块得到上一区块总难度和握手时得到的节点p的总难度对比，

如果是比较新的总难度，我们记录一下节点p的新的头部哈希和总难度，然后拿本节点的总难度和

节点p的总难度对比，如果发现节点p的节点的总难度比本地的总难度要大，那么就与节点p进行区块同步。



关于同步：

Downloader在系统中只生成了一个对象（实例）。

Downloader的synchronising字段来控制是否在同步，也就是不会说有多个节点同一个时间和本节点进行同步。

downloader进行New的时候要特别注意下文的注释

```
// New creates a new downloader to fetch hashes and blocks from remote peers.
func New(mode SyncMode, stateDb ethdb.Database, mux *event.TypeMux, chain BlockChain, lightchain LightChain, dropPeer peerDropFn) *Downloader {
    //注意
	//特别注意这里呀，当lightchain为nil的情况下，它把chain赋值给了lightchain
	if lightchain == nil {
		lightchain = chain
	}
...
```

程序一启动就会启动一个goroutine来轮询，根据续断器和新节点通道来触发，然后从已经建立连接的节点中选择难度最大的那个节点进行同步。



在eth包中的protocolManager中调用fetcher中的方法，我们看到它就处理4条消息

```
pm.fetcher.Notify(p.id, block.Hash, block.Number, time.Now(), p.RequestOneHeader, p.RequestBodies) //NewBlockHashesMsg

pm.fetcher.Enqueue(p.id, request.Block) //NewBlockMsg

headers = pm.fetcher.FilterHeaders(p.id, headers, time.Now()) //BlockHeadersMsg

pm.fetcher.FilterBodies(p.id, transactions, uncles, time.Now())//BlockBodiesMsg

```



当从网络中收到一个完整的区块时的处理（NewBlockMsg）：

1、Fetcher直接将它放入通道中

从优先队列中取出

1、如果它是一个未来区块（比下一个区块的高度要大），那么把它放入队列中。

2、如果它既然不是一个叔块（在当前区块链高度往前maxUncleDist之内），或者它已经存在于本地数据库中，那么我们直接忽略该区块

3、不是上述两种情况，那么它可能是下一个区块或者是叔块，那么我们就进入下一步处理

4、用该区块的ParentHash来获取区块，看它的上一个区块是否存在，如果不存在，那么我们就退出

5、如果它的上一个区块存在，那么我们快速检验该区块头，合法我们就将该区块广播出去。不合法的话，我们不广播区块。然后我们把区块插入链中，并广播该区块的哈希（我们虽然不广播整个区块，但我们还是广播它的哈希—尽管它的头部校验失败了）

6、把区块插入链 这步详情待写。



谁告诉你有新区块的哈希，你就向谁要完整的区块头和区块体。

Fetcher干了什么事情？

一句话概括：收到新区块消息的时候，放入自己的队列中，然后插入到区块链，收到哈希的时候，去获取区块头，去获取区块体，然后把区块头和区块体组合起来，然后放入自己的队列中，然后插入到区块链。

进行同步的时候，并没有完整的区块消息到来，所以NewBlockMsg都是由Fetcher来处理，而区块头（BlockHeadersMsg）和区块体（BlockBodiesMsg），都经过Fetcher的判断以后，才决定是由Fetcher处理还是交由downloader来处理。



downloader也是处理4条消息

```
pm.downloader.DeliverHeaders(p.id, headers) //BlockHeadersMsg,如果不是fetcher发出的

pm.downloader.DeliverBodies(p.id, transactions, uncles) //BlockBodiesMsg,如果不是fetcher发出的

pm.downloader.DeliverNodeData(p.id, data) //NodeDataMsg

pm.downloader.DeliverReceipts(p.id, receipts) //ReceiptsMsg
```



在fast或light下，得到区块头以后是直接调用InsertHeaderChain插入

在full或fast下,得到区块头以后是直接调用d.queue.Schedule(chunk, origin)

```
blockTaskPool和blockTaskQueue
```