**PBFT（拜占庭容错)**

他的核心思想是：对于每一个收到命令的将军，都要去询问其他人，他们收到的命令是什么。

采用PBFT方法，本质上就是利用通信次数换取信用。每个命令的执行都需要节点间两两交互去核验消息，通信代价是非常高的。通常采用PBFT算法，节点间的通信复杂度是节点数的平方级的。



拜占庭容错是一种共识算法，即如何使远距离通信的人们对一个提案达成一致意见。

与普通的共识算法（例如，majority wins，即超过一半人赞成即有效）不同的是，**PBFT可以容忍投票的人中产生叛徒或者不响应**。

举个合适的例子，我是一个愚昧的国王，没有自己的判断力，我不知道应该对敌国进攻还是投降好。我有一些大臣，我希望听从他们的意见做出决定，但是他们现在都离我很远，我只能通过飞鸽传书的方式告知他们目前的问题，得到他们的选择。然而，可能出现大臣叛变，故意提出相反的观点（错误的节点），也可能出现鸽子在传输过程中飞错了，我没有得到该大臣的选择（网络堵塞）。PBFT可以保证如果我有3f+1的大臣的话，即使其中有f个大臣叛变或者没有响应，我依然可以得出共识的正确结果。

**为什么有f个节点未响应或出错时，为了保证系统的正常，我需要总共有3f+1个节点进行投票。**

同样用国王的例子，假设除了我（国王）一共有n个大臣，我知道其中有f个大臣是叛徒或者未响应，所以我一定要能从n-f个大臣的回应中进行判断（如果上述f个大臣都是未响应）。然而由于是飞鸽传书（异步传输），所以当我陆续收到n-f个传来的消息后，我并不知道之后是否还会有新的消息传来。因为如果f个有问题的大臣都是未响应，那么我将不会收到新的消息，如果其中有大臣是叛徒，我之后还会收到消息，但作为国王的现在不知道是哪种情况，却需要立刻作出进攻还是投降的判断。

最坏的情况下，剩下的f个大臣都是好人，只是鸽子飞得慢我还没收到消息，也就是说我收到消息的n-f的大臣中有f个大臣都是叛徒，即f个叛徒和n-2f个好人。由于多数者胜，所以只有当n-2f>f的情况下，作为国王的我会做出正确的决定，即n>3f，n最小需要取3f+1。





​     基于拜占庭将军问题，一致性的确保主要分为这三个阶段：预准备（pre-prepare）、准备(prepare)和确认(commit)。流程如下图所示：

![](/Users/daminyang/github/github.com/teachain/study/ethereum/PBFT.png)



 

 其中C为发送请求端，0123为服务端，3为宕机的服务端，具体步骤如下：

1. Request：请求端C发送请求到任意一节点，这里是0

2. Pre-Prepare：服务端0收到C的请求后进行广播，扩散至123

3. Prepare：123,收到后记录并再次广播，1->023，2->013，3因为宕机无法广播

4. Commit：0123节点在Prepare阶段，若收到超过一定数量的相同请求，则进入Commit阶段，广播Commit请求

 5.Reply：0123节点在Commit阶段，若收到超过一定数量的相同请求，则对C进行反馈

 

 根据上述流程，在 N ≥ 3F + 1 的情況下一致性是可能解決，N为总计算机数，F为有问题的计算机总数。

 ![](/Users/daminyang/github/github.com/teachain/study/ethereum/pbft_2.png)

​     由此可以看出，拜占庭容错能够容纳将近1/3的错误节点误差，IBM创建的Hyperledger就是使用了该算法作为共识算法。

