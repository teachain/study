用geth console进入geth的console以后进可以使用以下的命令进行操作。



```
//创建新账号
personal.newAccount() //方式一：回车后需要自己输入两次密码，才会生成账号
personal.newAccount("123456") //方式二：回车后就直接生成好了账号

//查看节点信息
admin.nodeInfo

//开始挖矿
miner.start() //方式一：不指定使用的线程数
miner.start(线程数) //方式二：指定使用的线程数

//停止挖矿
miner.stop()

//查看当前矿工账号
eth.coinbase

//修改矿工账号
miner.setEtherbase("需要设置的矿工地址")

//查看账户信息
eth.accounts[0]

//查看账户余额
eth.getBalance(eth.accounts[0]) //方式一
web3.fromWei(eth.getBalance(eth.accounts[0]),"ether")//方式二
 
//解锁账号(使用账户资金前都需要先解锁账号)
personal.unlockAccount(eth.accounts[0]) //方式一：下一步输入密码
personal.unlockAccount(eth.accounts[0],"123456") //方式二：指定密码

//转账
eth.sendTransaction({from:eth.accounts[0],to: "0x27bb95f3cad3910189f17984b7c5e9fdf5eb5cda",value:web3.toWei(3,"ether")})

//查看交易池状态
txpool.status

//查看最新的区块高度
eth.blockNumber

//通过区块号查看区块
eth.getBlock(1)

```



