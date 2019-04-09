# 1. 以太坊中PoA产生的背景

如果你想用以太坊搭建一个联盟/私有链, 并要求该链交易成本更低甚至没有, 交易延时更低,并发更高, 还拥有完全的控制权(意味着被攻击概率更低). 目前以太坊采用PoW或后续的casper能否满足要求?

- 首先, pow存在**51%攻击问题**, 恶意挖矿者超过全网算力的51%后基本上就能完全控制整个网络. 由于链无法被更改, 已上链的数据也无法更改, 但恶意挖矿者也可以做一些DoS攻击阻止合法交易上链,考虑到具有相同创世块的旷工都能加入你的网络, 潜在的安全隐患会长期存在.
- 其次, PoW**大量的电力资源消耗**也是需要作为后续成本考虑. PoS可以解决部分Pow问题, 比如节约电力,在一定程度上保护了51％的攻击(恶意旷工会被惩罚), 但从控制权和安全考虑还有欠缺, 因为PoS还是允许任何符合条件的旷工加入。

在已经运行的测试网络Ropsten中, 由于pow设定的难度较低,恶意旷工滥用较低的PoW难度并将gaslimit扩大到90亿（正常是470万），发送大量的交易瘫痪了整个网络。而在此之前，攻击者也尝试了多次非常长的重组(reorgs)，导致不同客户端之间的分叉，甚至不同的版本。

这些攻击的根本原因是PoW网络的安全性依赖于背后的算力。而从零开始重新启动一个新的testnet将不会解决任何问题，因为攻击者可以一次又一次地进行相同的攻击。 Parity团队决定采取紧急解决方案，回滚大量的块，并设置不允许gaslimit超过某一阈值的软分叉规则。

虽然Parity的解决方案可能在短期内有效, 但是这不是优雅的：Ethereum本身应该具有动态gaslimit限制; 也不可移植：其他客户端需要自己实现新的软分叉逻辑; 并与同步模式不兼容, 也不支持轻客户端; 尽管并不完美，但是Parity的解决方案仍然可行。 一个更长期的替代解决方案是使用PoA共识,相对简单并容易实现.

# 2. PoA的特点

- PoA是依靠预设好的授权节点(signers)，负责产生block.
- 可以由已授权的signer选举(投票超过50%)加入新的signer。
- 即使存在恶意signer,他最多只能攻击连续块(数量是 `(SIGNER_COUNT / 2) + 1)` 中的1个,期间可以由其他signer投票踢出该恶意signer。
- 可指定产生block的时间。

# 3. PoA需要解决的问题

1. 如何控制挖矿频率,即出块时间
2. 如何验证某个块的有效性
3. 如何动态调整授权签名者(signers)列表,并全网动态同步
4. 如何在signers之间分配挖矿的负载或者叫做挖矿的机会

对应的解决办法如下:

1. 协议规定采用固定的block出块时间, 区块头中的时间戳间隔为15s
2. 先看看block同步的方法,从中来分析PoA中验证block的解决办法

有两种同步blockchain的方法

1. 经典方法是从创世块开始挨个执行所有交易。 这是经过验证的，但是在Ethereum的复杂网络中，计算量非常大。
2. 另一个是仅下载区块头并验证其有效性，之后可以从网络下载任意的近期状态对最近的区块头进行检查。

显然第二种方法更好. 由于PoA方案的块可以仅由可信任的签名者来创建, 因此，客户端看到的每个块（或块头）可以与可信任签名者列表进行匹配。 要验证该块的有效性就必须得到该块对应的签名者列表, 如果签名者在该列表中带包该块有效. 这里的挑战是如何维护并及时更改的授权签名者列表？ 存储在智能合约中?不可行, 因为在快速轻量级同步期间无法访问状态。

因此, **授权签名者列表必须完全包含在块头中** 。那么需要改变块头的结构, 引入新的字段来满足投票机制吗? 这也不可行：改变这样的核心数据结构将是开发者的噩梦。

所以授权签名者名单必须完全符合当前的数据模型, 不能改变区块头中的字段，而是 **复用当前可用的字段: Extra字段. **

**Extra** 是可变长数组, 对它的修改是 `非侵入` 操作, 比如RLP,hash操作都支持可变长数据. Extra中包含所有签名者列表和当前节点的签名者对该区块头的签名数据(可以恢复出来签名者的地址).

1. 更新一个动态的签名者列表的方法是复用区块头中的 **Coinbase和Nonce字段** ，以创建投票方案：

- 常规的块中这两个字段置为0
- 如果签名者希望对授权签名者列表进行更改，则将： 
  -  **Coinbase** 设置为被投票的签名者
  - 将 **Nonce** 设置为 **0** 或 **0xff ... f** 投票,代表 **添加或移除** 
  - 任何同步的客户端都可以在块处理过程中“统计”投票，并通过投票结果来维护授权签名者列表。

为了避免一个无限的时间来统计投票，我们设置一个投票窗口, 为一个epoch,长度是30000个block。每个epoch的起始清空所有历史的投票, 并作为签名者列表的检查点. 这允许客户端仅基于检查点哈希进行同步，而不必重播在链路上完成的所有投票。

1. 目前的方案是在所有signer之间轮询出块, 并通过算法保证同一个signer只能签名 `(SIGNER_COUNT / 2) + 1)` 个block中第一个.

综上, PoA的工作流程如下:

1. 在创世块中指定一组初始授权的signers, **所有地址** 保存在创世块Extra字段中
2. 启动挖矿后, 该组signers开始对生成的block进行 **签名并广播.** 
3.  **签名结果** 保存在区块头的Extra字段中
4. Extra中更新当前高度已授权的 **所有signers的地址** ,因为有新加入或踢出的signer
5. 每一高度都有一个signer处于IN-TURN状态, 其他signer处于OUT-OF-TURN状态, IN-TURN的signer签名的block会 **立即广播** , OUT-OF-TURN的signer签名的block会 **延时** 一点随机时间后再广播, 保证IN-TURN的签名block有更高的优先级上链
6. 如果需要加入一个新的signer, signer通过API接口发起一个proposal, 该proposal通过复用区块头 **Coinbase(新signer地址)和Nonce("0xffffffffffffffff")** 字段广播给其他节点. 所有已授权的signers对该新的signer进行"加入"投票, 如果赞成票超过signers总数的50%, 表示同意加入
7. 如果需要踢出一个旧的signer, 所有已授权的signers对该旧的signer进行"踢出"投票, 如果赞成票超过signers总数的50%, 表示同意踢出

## signer对区块头进行签名

1. Extra的长度至少65字节以上(签名结果是65字节,即R, S, V, V是0或1)
2. 对blockHeader中所有字段除了Extra的 **后65字节** 外进行 **RLP编码** 
3. 对编码后的数据进行 `Keccak256` **hash** 
4. 签名后的数据(65字节)保存到Extra的 **后65字节** 中

## 授权策略

以下建议的策略将减少网络流量和分叉

- 如果签名者被允许签署一个块（在授权列表中，但最近没有签名）。 
  - 计算下一个块的最优签名时间（父块时间+ BLOCK_PERIOD）。
  - 如果签名人是in-turn，立即进行签名和广播block。
  - 如果签名者是out-of-turn，延迟 `rand(SIGNER_COUNT * 500ms)` 后再签名并广播

## 级联投票

当移除一个授权的签名者时,可能会导致其他移除前的投票成立. 例: ABCD4个signer, AB加入E,此时不成立(没有超过50%), 如果ABC移除D, 会自动导致加入E的投票成立(2/3的投票比例)

## 投票策略

因为blockchain可能会小范围重组(small reorgs), 常规的投票机制(cast-and-forget, 投票和忘记)可能不是最佳的，因为包含单个投票的block可能不会在最终的链上,会因为已有最新的block而被抛弃。

一个简单但有效的办法是对signers配置"提议(proposal)".例如 "add 0x...", "drop 0x...", 有多个并发的提议时, 签名代码"随机"选择一个提议注入到该签名者签名的block中,这样多个并发的提议和重组(reorgs)都可以保存在链上.

该列表可能在一定数量的block/epoch 之后过期，提案通过并不意味着它不会被重新调用，因此在提议通过时不应立即丢弃。

- 加入和踢除新的signer的投票都是立即生效的,参与下一次投票计数
- 加入和踢除都需要 **超过当前signer总数的50%** 的signer进行投票
- 可以踢除自己(也需要超过50%投票)
- 可以并行投票(A,B交叉对C,D进行投票), 只要最终投票数操作50%
- 进入一个新的epoch, 所有之前的pending投票都作废, 重新开始统计投票

## 投票场景举例

- ABCD, AB先分别踢除CD, C踢除D, 结果是剩下ABC
- ABCD, AB先分别踢除CD, C踢除D, B又投给C留下的票, 结果是剩下ABC
- ABCD, AB先分别踢除CD, C踢除D, 即使C投给自己留下的票, 结果是剩下AB
- ABCDE, ABC先分别加入F(成功,ABCDEF), BCDE踢除F(成功,ABCDE), DE加入F(失败,ABCDE), BCD踢除A(成功, BCDE), B加入F(此时BDE加入F,满足超过50%投票), 结果是剩下BCDEF

# 4. PoA中的攻击及防御

- 恶意签名者(Malicious signer). 恶意用户被添加到签名者列表中，或签名者密钥/机器遭到入侵. 解决方案是，N个授权签名人的列表，任一签名者只能对每K个block签名其中的1个。这样尽量减少损害，其余的矿工可以投票踢出恶意用户。
- 审查签名者(Censoring signer). 如果一个签名者（或一组签名者）试图检查block中其他signer的提议(特别是投票踢出他们), 为了解决这个问题，我们将签名者的允许的挖矿频率限制在1/(N/2)。如果他不想被踢出出去, 就必须控制超过50%的signers.
- "垃圾邮件"签名者(Spamming signer). 这些signer在每个他们签名的block中都注入一个新的投票提议.由于节点需要统计所有投票以创建授权签名者列表, 久而久之之后会产生大量垃圾的无用的投票, 导致系统运行变慢.通过epoch的机制,每次进入新的epoch都会丢弃旧的投票
- 并发块(Concurrent blocks). 如果授权签名者的数量为N，我们允许每个签名者签名是1/K，那么在任何时候，至少N-K个签名者都可以成功签名一个block。为了避免这些block竞争( **分叉** )，每个签名者生成一个新block时都会加一点随机延时。这确保了很难发生分叉。



# PoA共识引擎算法实现分析

## clique中一些概念和定义

-  **EPOCH_LENGTH** : epoch长度是30000个block, 每次进入新的epoch,前面的投票都被清空,重新开始记录,这里的投票是指加入或移除signer

-  **BLOCK_PERIOD** : 出块时间, 默认是15s

-  **UNCLE_HASH** : 总是 `Keccak256(RLP([]))` ,因为没有uncle

-  **SIGNER_COUNT** : 每个block都有一个signers的数量

- SIGNER_LIMIT

   : 等于 

  ```
  (SIGNER_COUNT / 2) + 1
  ```

   . 每个singer只能签名连续SIGNER_LIMIT个block中的1个 

  - 比如有5个signer:ABCDE, 对4个block进行签名, 不允许签名者为ABAC, 因为A在连续3个block中签名了2次

-  **NONCE_AUTH** : 表示投票类型是加入新的signer; 值= `0xffffffffffffffff` 

-  **NONCE_DROP** : 表示投票类型是踢除旧的的signer; 值= `0x0000000000000000` 

-  **EXTRA_VANITY** : 代表block头中Extra字段中的保留字段长度: 32字节

-  **EXTRA_SEAL** : 代表block头中Extra字段中的存储签名数据的长度: 65字节

-  **IN-TURN/OUT-OF-TURN** : 每个block都有一个in-turn的signer, 其他signers是out-of-turn, in-turn的signer的权重大一些, 出块的时间会快一点, 这样可以保证该高度的block被in-turn的signer挖到的概率很大.

clique中最重要的两个数据结构:

- 共识引擎的结构:

```
    type Clique struct {
        config *params.CliqueConfig // 系统配置参数
        db ethdb.Database // 数据库: 用于存取检查点快照
        recents *lru.ARCCache //保存最近block的快照, 加速reorgs
        signatures *lru.ARCCache //保存最近block的签名, 加速挖矿
        proposals map[common.Address]bool //当前signer提出的proposals列表
        signer common.Address // signer地址
        signFn SignerFn // 签名函数
        lock sync.RWMutex // 读写锁
    }
```

- snapshot的结构:

```
    type Snapshot struct {
        config *params.CliqueConfig // 系统配置参数
        sigcache *lru.ARCCache // 保存最近block的签名缓存,加速ecrecover
        Number uint64 // 创建快照时的block号
        Hash common.Hash // 创建快照时的block hash
        Signers map[common.Address]struct{} // 此刻的授权的signers
        Recents map[uint64]common.Address // 最近的一组signers, key=blockNumber
        Votes []*Vote // 按时间顺序排列的投票列表
        Tally map[common.Address]Tally // 当前的投票计数，以避免重新计算
    }
```

除了这两个结构, 对block头的部分字段进行了复用定义, ethereum的block头定义:

```
    type Header struct {
        ParentHash common.Hash 
        UncleHash common.Hash 
        Coinbase common.Address 
        Root common.Hash 
        TxHash common.Hash 
        ReceiptHash common.Hash 
        Bloom Bloom 
        Difficulty *big.Int 
        Number *big.Int 
        GasLimit *big.Int 
        GasUsed *big.Int 
        Time *big.Int 
        Extra []byte 
        MixDigest common.Hash 
        Nonce BlockNonce 
    }
```

- 创世块中的Extra字段包括: 

  - 32字节的前缀(extraVanity)
  - 所有signer的地址
  - 65字节的后缀(extraSeal): 保存signer的签名

- 其他block的Extra字段只包括extraVanity和extraSeal

- Time字段表示产生block的时间间隔是:blockPeriod(15s)

- Nonce字段表示进行一个投票: 添加( nonceAuthVote: `0xffffffffffffffff` )或者移除( nonceDropVote: `0x0000000000000000` )一个signer

- Coinbase字段存放 

  被投票

   的地址 

  - 举个栗子: signerA的一个投票:加入signerB, 那么Coinbase存放B的地址

- Difficulty字段的值: 1-是 **本block的签名者** (in turn), 2- **非本block的签名者** (out of turn)

下面对比较重要的函数详细分析实现流程

## Snapshot.apply(headers)

创建一个新的授权signers的快照, 将从上一个snapshot开始的区块头中的proposals更新到最新的snapshot上

1. 对入参headers进行完整性检查: 因为可能传入多个区块头, **block号必须连续** 
2. 遍历所有的header, 如果block号刚好处于epoch的起始(number%Epoch == 0),将snapshot中的Votes和Tally复位( **丢弃历史全部数据** )
3. 对于每一个header,从签名中恢复得到 **signer** 
4. 如果该signer在snap.Recents中, 说明 **最近已经有过签名** , 不允许再次签名, 直接 **返回** 结束
5.  **记录** 该signer是该block的签名者: `snap.Recents[number] = signer` 
6. 统计header.Coinbase的投票数,如果 **超过signers总数的50%** 
7. 执行加入或移除操作
8. 删除snap.Recents中的一个signer记录: key=number- (uint64(len(snap.Signers)/2 + 1)), 表示释放该signer,下次可以对block进行签名了
9. 清空被移除的Coinbase的投票
10. 移除snap.Votes中该Conibase的所有投票记录
11. 移除snap.Tally中该Conibase的所有投票数记录

## 共识引擎clique的初始化

在 `Ethereum.StartMining` 中,如果Ethereum.engine配置为clique.Clique, 根据当前节点的矿工地址(默认是acounts[0]), 配置clique的 **签名者** : `clique.Authorize(eb, wallet.SignHash)` ,其中 **签名函数** 是SignHash,对给定的hash进行签名.

## 获取给定时间点的一个快照 Clique.snapshot

- 先查找Clique.recents中是否有缓存, 有的话就返回该snapshot

- 在查找持久化存储中是否有缓存, 有的话就返回该snapshot

- 如果是创世块 

  1. 从Extra中取出所有的signers
  2. `newSnapshot(Clique.config, Clique.signatures, 0, genesis.Hash(), signers)`

  - signatures是最近的签名快照
  - signers是所有的初始signers

  1. 把snapshot加入到Clique.recents中, 并持久化到db中

- 其他普通块 

  - 沿着父块hash一直往回找是否有snapshot, 如果没找到就记录该区块头
  - 如果找到最近的snapshot, 将前面记录的headers 都 `applay` 到该snapshot上
  - 保存该最新的snapshot到缓存Clique.recents中, 并持久化到db中

## Clique.Prepare(chain , header)

Prepare是共识引擎接口之一. 该函数配置header中共识相关的参数(Cionbase, Difficulty, Extra, MixDigest, Time)

- 对于非epoch的block( `number % Epoch != 0` ):

1. 得到Clique.proposals中的投票数据(例:A加入C, B踢除D)
2. 根据snapshot的signers分析投票数否有效(例: C原先没有在signers中, 加入投票有效, D原先在signers中,踢除投票有效)
3. 从被投票的地址列表(C,D)中, **随机选择一个地址** ,作为该header的Coinbase,设置Nonce为加入( `0xffffffffffffffff` )或者踢除( `0x0000000000000000` )
4.  `Clique.signer` 如果是本轮的签名者(in-turn), 设置header.Difficulty = diffInTurn(1), 否则就是diffNoTurn(2)
5. 配置header.Extra的数据为[ `extraVanity` + `snap中的全部signers` + `extraSeal` ]
6. MixDigest需要配置为nil
7. 配置时间戳:Time为父块的时间+15s

## 重点: Clique.Seal(chain, block , stop)

Seal也是共识引擎接口之一. 该函数用clique.signer对block的进行签名. 在pow]算法中, 该函数进行hash运算来解"难题".

- 如果signer没有在snapshot的signers中,不允许对block进行签名
- 如果不是本block的签名者,延时一定的时间(随机)后再签名, 如果是本block的签名者, 立即签名.
- 签名结果放在Extra的extraSeal的65字节中

## Clique.VerifySeal(chain, header)

VerifySeal也是共识引擎接口之一.

1. 从header的签名中恢复账户地址,改地址要求在snapshot的signers中
2. 检查header中的Difficulty是否匹配(in turn或out of turn)

## Clique.Finalize

Finalize也是共识引擎接口之一. 该函数生成一个block, 没有叔块处理,也没有奖励机制

1.  `header.Root` : 状态根保持原状
2.  `header.UncleHash` : 为nil
3.  `types.NewBlock(header, txs, nil, receipts)` : 封装并返回最终的block

## API.Propose(addr, auth)

添加一个proposal: 调用者对addr的投票, auth表示加入还是踢出

## API.Discard(addr)

删除一个proposal

