工作量证明的目的是使区块的创建变得困难，从而阻止女巫攻击者恶意重新生成区块链。因为SHA256是完全不可预测的伪随机函数，创建有效区块的唯一方法就是简单地不断试错，不断地增加随机数的数值，查看新的哈希数值是否小于目标数值。

以太坊的虚拟机的代码研读和预编译合约的学习。

这个package主要是共识算法

consensus/consensus.go主要是包含了以下接口

1. ChainReader
2. Engine
3. Pow  （注意到Engine作为内嵌的接口接入在该接口中）

consensus/ethash主要是以太坊的pow共识算法

consensus/ethash/ethash.go和consensus/ethash/consensus.go是结构体Ethash的定义和方法的实现。

Ethash是接口Engine的实现，当然也是接口Pow实现。



```

```