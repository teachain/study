**PoA共识机制**

所谓授权证明PoA(Proof of Authority），就是由一组授权节点来负责新区块的产生和区块验证。

以太坊源码中带有Clique共识算法即为一种PoA共识算法。

在PoA中，验证者（validator）是整个共识机制的关键。验证者不需要昂贵的显卡，也不需要足够的资产，但他必须具有已知的，并且已获得验证的身份。验证者通过放置这个身份来获得担保网络的权利，从而换取区块奖励。若是验证者在整个过程中有恶意行为，或与其他验证者勾结。那通过链上管理可以移除和替换恶意行为者。现有的法律反欺诈保障会被用于整个网络的参与者免受验证者的恶意行为。



**PoA共识机制的特点**

PoA是依靠预设好的授权节点(signers)，负责产生block.可以由已授权的signer选举(投票超过50%)加入新的signer。即使存在恶意signer,他最多只能攻击连续块(数量是 `(SIGNER_COUNT / 2) + 1)` 中的1个,期间可以由其他signer投票踢出该恶意signer。可指定产生block的时间。



节点：普通的以太坊节点，没有区块生成的权利。

矿工：具有区块生成权利的以太坊节点

委员会：所有矿工的集合

在POA共识算法红，区块数据与Pow有些区别，主要体现在header结构中

- coinbase 被提名为矿工的节点地址
- nonce 提名分类，添加或者删除
- 在Epoch时间点，存储当前委员会集合Singners



**每一笔交易都要由网络中的 每一个节点 进行处理**。在以太坊区块链上进行的每一项操作都必须由网络中的每一个节点并行处理。区块链的[设计](https://link.jianshu.com?t=http%3A%2F%2Fwww.aibbt.com%2Fa%2Ftag%2F%25e8%25ae%25be%25e8%25ae%25a1%2F)就是这样的，这是使得公有链具有权威性的一部分。一个节点不需要依赖其他节点来告诉它区块链的当前状态是什么，它自己会搞清楚。

**这给以太坊交易吞吐量带来了根本性的限制**：它不能高于我们对单个节点所要求的交易吞吐量。