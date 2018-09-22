首先先确认你的terminal环境下，是否已经可以使用geth和puppeth命令，如果没有，自行百度谷歌。

在当前的用户工作目录（我的工作目录是/Users/daminyang/workspace ）下,创建目录ethereum，

接着在ethereum目录下创建三个目录,分别命名为node1,node2,node3。使用的命令如下：

打开terminal:

```
cd /Users/daminyang/workspace
mkdir ethereum
cd ethereum
mkdir node1
mkdir node2
mkdir node3
```

为了方便后面的测试，先创建三个账户

```
geth --datadir ./node1/data account new
geth --datadir ./node2/data account new
geth --datadir ./node3/data account new
```

这三个账户分别分布在不同的节点上。

我们使用puppeth工具来生成创世区块所需的json文件

```
puppeth(命令行输入)
+-----------------------------------------------------------+
| Welcome to puppeth, your Ethereum private network manager |
|                                                           |
| This tool lets you create a new Ethereum network down to  |
| the genesis block, bootnodes, miners and ethstats servers |
| without the hassle that it would normally entail.         |
|                                                           |
| Puppeth uses SSH to dial in to remote servers, and builds |
| its network components out of Docker containers using the |
| docker-compose toolset.                                   |
+-----------------------------------------------------------+

Please specify a network name to administer (no spaces or hyphens, please)
> simple_private_chain(命令行输入，名字由你自己取)

Sweet, you can set this via --network=simple_private_chain next time!

INFO [08-09|16:10:40.771] Administering Ethereum network           name=simple_private_chain
WARN [08-09|16:10:40.772] No previous configurations found         path=/Users/daminyang/.puppeth/simple_private_chain

What would you like to do? (default = stats)
 1. Show network stats
 2. Configure new genesis
 3. Track new remote server
 4. Deploy network components
> 2(命令行输入)
Which consensus engine to use? (default = clique)
 1. Ethash - proof-of-work
 2. Clique - proof-of-authority
> 1

Which accounts should be pre-funded? (advisable at least one)
> 0x7c63404e0dab5a90ff97f5525a00596f0768bd65 (这里输入一个你需要预先分配以太币的账户地址)
> 0x(回车)
Specify your chain/network ID if you want an explicit one (default = random)
> 665577(指定你想要的网络id)
INFO [08-09|16:15:34.541] Configured new genesis block 

What would you like to do? (default = stats)
 1. Show network stats
 2. Manage existing genesis
 3. Track new remote server
 4. Deploy network components
> 2

 1. Modify existing fork rules
 2. Export genesis configuration
 3. Remove genesis configuration
> 2

Which file to save the genesis into? (default = simple_private_chain.json)
> (回车)
INFO [08-09|16:15:46.845] Exported existing genesis block 

What would you like to do? (default = stats)
 1. Show network stats
 2. Manage existing genesis
 3. Track new remote server
 4. Deploy network components
> (ctrl+c,退出)

```

我们使用以下命令来初始化三个节点

```
geth --datadir ./node1/data init simple_private_chain.json
geth --datadir ./node2/data init simple_private_chain.json
geth --datadir ./node3/data init simple_private_chain.json
```

启动第一个节点

```
geth --networkid 665577 --nodiscover --datadir  ./node1/data --rpc --rpcapi net,eth,web3,personal  --rpcport 8001 -port 5001 console
miner.start() #开始挖矿，必须，不然下面的转账就没有办法进行了。也就是整个网络里必须至少有一个节点是在挖矿的
```

我们保持这个terminal的window,我们新开一个terminal的window

```
cd /Users/daminyang/workspace/ethereum
geth attach ipc:./node1/data/geth.ipc #这里我们使用ipc连接上了节点1(node1)
net.peerCount #查看连接的节点数量
admin.peers #查看连接的节点列表
eth.accounts #查看此节点下所有的账户
eth.getBalance(eth.accounts[0]) #查看第一个账户的余额
personal.newAccount("account") #新建一个账户,密码为account
eth.getBalance(eth.accounts[1]) #查看一下余额
personal.unlockAccount(eth.accounts[0] ,"123456") #解锁账户，我们原先设置的账户的密码统一为123456
eth.sendTransaction({from : eth.accounts[0], to : eth.accounts[1] , value : web3.toWei(1,"ether")}) #从第一个账户转给第二个账户账户1eth
admin.nodeInfo.enode #查看节点信息
```

从查看节点信息我们得到

"enode://267e1667def6553f662a1c2bd4ca602c0463e4fe91511728afb3520273c3594b60f8c9f6a5c8044d6c2a827f4bbaf9287253b8834f8dc846d7cee24bf9d3e4e5@[::]:5001?discport=0"

我们再新开一个terminal的window

```
cd /Users/daminyang/workspace/ethereum
geth --networkid 665577 --nodiscover --datadir  ./node2/data --rpc --rpcapi net,eth,web3,personal  --rpcport 8002 -port 5002 console #启动第一个节点
```

我们再新开一个terminal的window（我们将[::]替换成127.0.0.1）

```
cd /Users/daminyang/workspace/ethereum
geth attach ipc:./node2/data/geth.ipc #这里我们使用ipc连接上了节点1(node2)
net.peerCount #查看连接的节点数量
admin.peers #查看连接的节点列表
eth.accounts #查看此节点下所有的账户
admin.addPeer("enode://267e1667def6553f662a1c2bd4ca602c0463e4fe91511728afb3520273c3594b60f8c9f6a5c8044d6c2a827f4bbaf9287253b8834f8dc846d7cee24bf9d3e4e5@127.0.0.1:5001?discport=0")#添加节点
net.peerCount #查看连接的节点数量(再看一下咯)
admin.peers #查看连接的节点列表(再看一下咯)
```

同启动2一样，启动3

然后就可以在各个节点上进行转账什么的了

```

personal.unlockAccount(eth.accounts[0] ,"123456")

eth.sendTransaction({from : eth.accounts[0], to : '0x5d861a4f6412187d64ac4cb73ad87c0dbd9fa970' , value : web3.toWei(1,"ether")})
```

我们可以在任意一个节点上查询到账户的余额，但是我们必须是在账户创建的节点上解锁账户（因为账号是本地保存的，是在节点上的）

```
eth.coinbase #查看当前挖矿的账户
miner.setEtherbase("你需要切换的挖矿账户地址")#当前挖矿的账户
```

使用mist进行链接私有链

```
 /Applications/Mist.app/Contents/MacOS/Mist --rpc /Users/daminyang/workspace/ethereum/node1/data/geth.ipc
```

