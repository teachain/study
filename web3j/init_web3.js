//导入web3模块
var Web3 = require('web3');

//浏览器模式
if (window.ethereum) {
    window.web3 = new Web3(ethereum);
    try {
        //这一句主要是为了防止metamask开启隐私模式，调用获取账户地址获取不到的问题
        //也就是说代码里如果没有这一句，调用web3.eth.getAccounts()获取不到用户
        //基本上就是这个问题，解决方案：要么自己手动关闭metamask的的隐私模式，要么通过以下代码
        // Request account access if needed
        await ethereum.enable();
        // Acccounts now exposed
        web3.eth.sendTransaction({ /* ... */ });
    } catch (error) {
        // User denied account access...
    }
}
// Legacy dapp browsers...
else if (window.web3) {
    window.web3 = new Web3(web3.currentProvider);
    // Acccounts always exposed
    web3.eth.sendTransaction({ /* ... */ });
}
// Non-dapp browsers...
else {
    console.log('Non-Ethereum browser detected. You should consider trying MetaMask!');
}