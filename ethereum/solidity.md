对合约的状态变量使用public修饰的话，会自动生成一个针对该状态变量的一个get函数（与状态变量同名的一个函数---------编译器帮你实现的）。

mapping (address => uint) public balances   

mapping类型，它既不能获得映射的所有键的列表，也不能获得所有值的列表。

`msg.sender` 始终是当前（外部）函数调用的来源地址。这个一定要切记，在函数内，它一定是表示函数调用的

来源地址。（谁调用就表示谁的地址）。

### 特别注意

```
记住，如果你使用合约发送币给一个地址，当你在区块链浏览器上查看该地址时是看不到任何相关信息的。因为，实际上你发送币和更改余额的信息仅仅存储在特定合约的数据存储器中。通过使用事件，你可以非常简单地为你的新币创建一个“区块链浏览器”来追踪交易和余额。(这里的币指的是你自己发行的币----并不是以太币）
```



交易总是由发送人（创建者）签名。



在比特币中，要解决的一个主要难题，被称为“双花攻击 (double-spend attack)”：如果网络存在两笔交易，都想花光同一个账户的钱时（即所谓的冲突）会发生什么情况？交易互相冲突？

简单的回答是你不必在乎此问题。网络会为你自动选择一条交易序列，并打包到所谓的“区块”中，然后它们将在所有参与节点中执行和分发。如果两笔交易互相矛盾，那么最终被确认为后发生的交易将被拒绝，不会被包含到区块中。

作为“顺序选择机制”（也就是所谓的“挖矿”）的一部分，可能有时会发生块（blocks）被回滚的情况，但仅在链的“末端”。末端增加的块越多，其发生回滚的概率越小。因此你的交易被回滚甚至从区块链中抹除，这是可能的，但等待的时间越长，这种情况发生的概率就越小。



### 账户

以太坊中有两类账户（它们共用同一个地址空间）： **外部账户** 由公钥-私钥对（也就是人）控制； **合约账户** 由和账户一起存储的代码控制.

外部账户的地址是由公钥决定的，而合约账户的地址是在创建该合约时确定的（这个地址通过合约创建者的地址和从该地址发出过的交易数量计算得到的，也就是所谓的“nonce”）

无论帐户是否存储代码，这两类账户对 EVM 来说是一样的。

每个账户都有一个键值对形式的持久化存储。其中 key 和 value 的长度都是256位，我们称之为 **存储** 。

此外，每个账户有一个以太币余额（ **balance** ）（单位是“Wei”），余额会因为发送包含以太币的交易而改变。



### 交易

交易可以看作是从一个帐户发送到另一个帐户的消息（这里的账户，可能是相同的或特殊的零帐户，请参阅下文）。它能包含一个二进制数据（合约负载）和以太币。

如果目标账户含有代码，此代码会被执行，并以 payload 作为入参。

如果目标账户是零账户（账户地址为 `0` )，此交易将创建一个 **新合约** 。 如前文所述，合约的地址不是零地址，而

是通过合约创建者的地址和从该地址发出过的交易数量计算得到的（所谓的“nonce”）。 这个用来创建合约的交易

的 payload 会被转换为 EVM 字节码并执行。执行的输出将作为合约代码被永久存储。这意味着，为创建一个合

约，你不需要发送实际的合约代码，而是发送能够产生合约代码的代码。





```
     //定义
     modifier onlySeller() { // 修饰器
        require(
            msg.sender == seller,
            "Only seller can call this."
        );
        _;
    }
    //使用onlySeller
     function abort() public onlySeller { // Modifier usage
        // ...
    }
    
    event HighestBidIncreased(address bidder, uint amount); // 定义事件
    
    function bid() public payable {
        // ...
        emit HighestBidIncreased(msg.sender, msg.value); // 触发事件
    }
     struct Voter { // 结构
        uint weight;
        bool voted;
        address delegate;
        uint vote;
    }
    contract Purchase {
    enum State { Created, Locked, Inactive } // 枚举
    }
```



合约继承

```
//使用关键字is
contract ZombieFeeding is ZombieFactory {

}
```

在 Solidity 中，当你有多个文件并且想把一个文件导入另一个文件时，可以使用 `import` 语句：

```
import "./someothercontract.sol";

contract newContract is SomeOtherContract {

}
```

这样当我们在合约（contract）目录下有一个名为 `someothercontract.sol`的文件（ `./` 就是同一目录的意思），它就会被编译器导入。

在 Solidity 中，有两个地方可以存储变量 —— `storage` 或 `memory`。

**Storage** 变量是指永久存储在区块链中的变量。 **Memory** 变量则是临时的，当外部函数对某合约调用完成时，内存型变量即被移除。 你可以把它想象成存储在你电脑的硬盘或是RAM中数据的关系。

大多数时候你都用不到这些关键字，默认情况下 Solidity 会自动处理它们。 状态变量（在函数之外声明的变量）默认为“存储”形式，并永久写入区块链；而在函数内部声明的变量是“内存”型的，它们函数调用结束后消失。

然而也有一些情况下，你需要手动声明存储类型，主要用于处理函数内的 **结构体**和 **数组** 时：

```
contract SandwichFactory {
  struct Sandwich {
    string name;
    string status;
  }

  Sandwich[] sandwiches;

  function eatSandwich(uint _index) public {
    // Sandwich mySandwich = sandwiches[_index];

    // ^ 看上去很直接，不过 Solidity 将会给出警告
    // 告诉你应该明确在这里定义 `storage` 或者 `memory`。

    // 所以你应该明确定义 `storage`:
    Sandwich storage mySandwich = sandwiches[_index];
    // ...这样 `mySandwich` 是指向 `sandwiches[_index]`的指针
    // 在存储里，另外...
    mySandwich.status = "Eaten!";
    // ...这将永久把 `sandwiches[_index]` 变为区块链上的存储

    // 如果你只想要一个副本，可以使用`memory`:
    Sandwich memory anotherSandwich = sandwiches[_index + 1];
    // ...这样 `anotherSandwich` 就仅仅是一个内存里的副本了
    // 另外
    anotherSandwich.status = "Eaten!";
    // ...将仅仅修改临时变量，对 `sandwiches[_index + 1]` 没有任何影响
    // 不过你可以这样做:
    sandwiches[_index + 1] = anotherSandwich;
    // ...如果你想把副本的改动保存回区块链存储
  }
}
```

如果你还没有完全理解究竟应该使用哪一个，也不用担心 —— 在本教程中，我们将告诉你何时使用 `storage` 或是 `memory`，并且当你不得不使用到这些关键字的时候，Solidity 编译器也发警示提醒你的。

现在，只要知道在某些场合下也需要你显式地声明 `storage` 或 `memory`就够了！

solidity支持普通函数（非constructor）的重载



在 Solidity 中，有两个地方可以存储变量 —— `storage` 或 `memory`。

**Storage** 变量是指永久存储在区块链中的变量。 **Memory** 变量则是临时的，当外部函数对某合约调用完成时，内存型变量即被移除。 你可以把它想象成存储在你电脑的硬盘或是RAM中数据的关系。

大多数时候你都用不到这些关键字，默认情况下 Solidity 会自动处理它们。 状态变量（在函数之外声明的变量）默认为“存储”形式，并永久写入区块链；而在函数内部声明的变量是“内存”型的，它们函数调用结束后消失。

然而也有一些情况下，你需要手动声明存储类型，主要用于处理函数内的 **结构体** 和 **数组** 时：



描述函数可见性的修饰词:

public

private

internal

external



如果我们的合约需要和区块链上的其他的合约会话，则需先定义一个**interface** (接口)。