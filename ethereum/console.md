因为在初始化的时候，console就先执行了一下代码

```
var eth = web3.eth; 

var personal = web3.personal;

var ethash = web3.ethash; 

var personal = web3.personal;

var txpool = web3.txpool; 

var net = web3.net; 

var rpc = web3.rpc; 

var admin = web3.admin; 

var debug = web3.debug; 

var eth = web3.eth; 

var miner = web3.miner;
```

所以你在控制台里可以直接使用eth.accounts获取本节点所有的内容。

所有有响应的rpct调用会经过这个方法的

```
func (c *Client) handleResponse(msg *jsonrpcMessage) {...}
```



https://github.com/ethereum/web3.js

console里定义了一个send的方法，这个方法是注册给js用的，也就是给jeth用

```

// init retrieves the available APIs from the remote RPC provider and initializes
// the console's JavaScript namespaces based on the exposed modules.
func (c *Console) init(preload []string) error {

	// Initialize the JavaScript <-> Go RPC bridge
	bridge := newBridge(c.client, c.prompter, c.printer)

	//在往js里注入一个对像
	c.jsre.Set("jeth", struct{}{})

	//从js里获取一个jeth对象
	jethObj, _ := c.jsre.Get("jeth")

	//为jeth添加一个方法 send，这样子，js就可以调用到go的代码bridge.Send了
	jethObj.Object().Set("send", bridge.Send)

	//为jeth添加一个方法 sendAsync，这样子，js就可以调用到go的代码bridge.Send了
	jethObj.Object().Set("sendAsync", bridge.Send)

	//从js里获取一个console对象
	consoleObj, _ := c.jsre.Get("console")

    //修改一下log默认的方法，让其调用go的方法c.consoleOutput
	consoleObj.Object().Set("log", c.consoleOutput)

	//修改一下error默认的方法，让其调用go的方法c.consoleOutput
	consoleObj.Object().Set("error", c.consoleOutput)

	//编译并执行bignumber.js代码
	// Load all the internal utility JavaScript libraries
	if err := c.jsre.Compile("bignumber.js", jsre.BigNumber_JS); err != nil {
		return fmt.Errorf("bignumber.js: %v", err)
	}
	//编译并执行web3.js代码
	if err := c.jsre.Compile("web3.js", jsre.Web3_JS); err != nil {
		return fmt.Errorf("web3.js: %v", err)
	}
	//执行代码，确认web3是否可用
	if _, err := c.jsre.Run("var Web3 = require('web3');"); err != nil {
		return fmt.Errorf("web3 require: %v", err)
	}
	//执行代码，确认web3和jeth是否可用,这里就用bridge.Send替换掉了web3j对象里的provider的send
	//方法，也就是当它发送rpc请求时，直接被go的代码接管过来了。原来通过http或websocket来建立的连接，
	//现在直接换成了net.Pipe的net.Conn
	if _, err := c.jsre.Run("var web3 = new Web3(jeth);"); err != nil {
		return fmt.Errorf("web3 provider: %v", err)
	}
	//...省略的代码
}
```

