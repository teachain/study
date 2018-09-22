# truffle

可以使用两种客户端geth和ethereumjs-testrpc，任选其一

如果需要使用ethereumjs-testrpc，那么就必须先安装

```
npm install -g ganache-cli # npm install -g ethereumjs-testrpc（这是老版本）
```

安装，在terminal中输入以下命令即可安装truffle

```
npm install -g truffle  #安装truffle
truffle version         #查看truffle的版本，也可确认是否成功安装truffle
```

常用的命令

```
truffle init            #初始化项目
truffle compile         #编译合约
truffle migrate         #部署合约，这里就表示已经将只能合约部署到了以太坊的区块链
truffle migrate --reset #来强制重编译并发布所有合约
truffle test            #测试合约
truffle unbox pet-shop  #使用pet-shop模板来创建一个Dapp项目

```



执行truffle compile命令后，在contracts目录下的智能合约文件均被编译成json文件。

这些json文件就是truffle用来部署合约的编译文件。

truffle.js是truffle的配置文件，启动好以太坊本地结点以后，我们需要让truffle去识别它并使用它，这就需要在

truffle.js中配置相关属性：

```
module.exports = {
  networks: {
    development: {
      host: "127.0.0.1",
      port: 8545,
      network_id: "*" // Match any network id
    }
  }
};
```

查看链接https://truffleframework.com/docs/truffle/reference/configuration即可知道详细的配置。



```
 Error encountered, bailing. Network state unknown. Review successful transactions manually.


 启动节点后,节点中默认的账户是被锁定的,无法执行部署合约的操作,需要调用下面的命令解锁账户:
        > personal.unlockAccount(eth.accounts[0],"password", 1000*60*20)
        
Error: authentication needed: password or unlock
Error: exceeds block gas limit
Error: intrinsic gas too low
```



```
//把json文件中的abi拷贝到这里来（压缩一下）
var abi = [{"inputs": [{"name": "_greeting","type": "string"}],"payable": false,"stateMutability": "nonpayable","type": "constructor"},{"constant": true,"inputs": [],"name": "greet","outputs": [{"name": "","type": "string"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": false,"inputs": [{"name": "_newgreeting","type": "string"}],"name": "setGreeting","outputs": [],"payable": false,"stateMutability": "nonpayable","type": "function"},{"constant": false,"inputs": [],"name": "kill","outputs": [],"payable": false,"stateMutability": "nonpayable","type": "function"}];
//实例化合约(abi和你自己部署好的合约地址---0xb52bb3ce336f71a14345c78e5b2f8e63685e3f92)
var HelloWorld = eth.contract(abi).at('0xb52bb3ce336f71a14345c78e5b2f8e63685e3f92')
//调用合约的方法
HelloWorld.greet()


```

