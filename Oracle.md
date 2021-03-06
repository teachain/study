## 预言机（Oracle）是什么？

**区块链外信息写入区块链内的机制，一般被称为预言机 (oracle mechanism) 。**

**预言机的功能就是将外界信息写入到区块链内，完成区块链与现实世界的数据互通。它允许确定的智能合约对不确定的外部世界作出反应，是智能合约与外部进行数据交互的唯一途径，也是区块链与现实世界进行数据交互的接口。**

###区块链为什么需要预言机？

**区块链是一个确定性的、封闭的系统环境，目前区块链只能获取到链内的数据，而不能获取到链外真实世界的数据，区块链与现实世界是割裂的。**

一般智能合约的执行需要触发条件，**当智能合约的触发条件是外部信息时（链外），就必须需要预言机来提供数据服务，通过预言机将现实世界的数据输入到区块链上**，因为智能合约不支持对外请求。

具体原因是这样的。区块链是确定性的环境，它不允许不确定的事情或因素，**智能合约不管何时何地运行都必须是一致的结果，所以虚拟机（VM）不能让智能合约有 network call（网络调用）**，不然结果就是不确定的。

我们通过一个例子来说明一下。

假设现在我被关进了一个小黑屋里（不要多想，只是例子🌝），我对外面的世界发生了什么一无所知，不知道外面是否有人，即使呼叫也没有人回应，而我知道外界信息的方式，只有外面的人在门口把他看到的听到的都告诉我，我才可以得知。

例子虽然不太恰当，但智能合约就像这个例子中的我一样，它无论何时何地，都无法主动向外寻求信息，只能外部把消息或数据给到里面。**而预言机就是这个在外面输送消息和数据的人。**

好像这么看来，智能合约并不是很智能呀，是的，智能合约其实是完成的不智能的事情，即写好了条件和结果，当给它条件的时候，就可以触发，但也不会马上执行，还需要合约相关的人进行私钥签署才可以执行。

所以，网上很多文章其实都有水分，比如智能合约某个时间或者触发某个条件就可以自动执行之类的，只能说这样的句子在逻辑上可能是有问题的。关于预言机的很多文章也有水分，描述的并不准确。

**好了，上面就是区块链为什么需要预言机，因为智能合约无法主动去获取链外的数据，只能被动接受数据。**


