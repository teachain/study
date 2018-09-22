在console下的命令

```
admin.nodeInfo #查看节点信息

eth.coinbase #查看当前矿工账号

miner.setEtherbase(eth.accounts[1]) #修改矿工账号

eth.sendTransaction({from:eth.accounts[0],to:eth.accounts[1],value:web3.toWei(1,"ether")}) #转账

txpool.status  #查看交易状态

eth.blockNumber #查看当前最大区块号

eth.getBlock(number) #查看区块高度为number的区块的数据

personal.unlockAccount(eth.accounts[0],'密码') #解锁账户，要操作资金，必须进行这个操作

miner.start(threads) #指定开多个线程开始挖矿

miner.stop() #停止挖矿


```

