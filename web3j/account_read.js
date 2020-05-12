//导入web3模块
var Web3 = require('web3');
//实例化web3
var web3 = new Web3(new Web3.providers.HttpProvider("http://localhost:8545"));
//方法一
//调用获取账号地址接口
//新版本metamask必须以异步的方式(加回调函数的方式)获取
console.log(web3.version);
web3.eth.getAccounts(function(error, result) {
	if (error) {
		console.log(error);
		return;
	}
	console.log("now has accounts:" + JSON.stringify(result));
});

//方法二（其实是同一个方法，这个只是秀一下肌肉而已）
async function getAccount() {
	var accounts = await web3.eth.getAccounts();
	console.log("async get accounts:" + JSON.stringify(accounts));
}
getAccount();