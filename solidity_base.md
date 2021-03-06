```
pragma solidity ^0.4.0;
```

Solidity中合约的含义就是一组代码（它的 *函数* )和数据（它的 *状态* ），它们位于以太坊区块链的一个特定地址上。

***注意***

* 一个函数和一个状态变量不能同名。
* `msg.sender` 始终是当前（外部）函数调用的来源地址。
* 记住，如果你使用合约发送币给一个地址，当你在区块链浏览器上查看该地址时是看不到任何相关信息的。因为，实际上你发送币和更改余额的信息仅仅存储在特定合约的数据存储器中。通过***使用事件***，你可以非常简单地为你的新币创建一个“区块链浏览器”来追踪交易和余额。
* 交易总是由发送人（创建者）签名。
* 如果目标账户含有代码，此代码会被执行，并以 payload(input或data) 作为入参。

Mappings 可以看作是一个 [哈希表](https://en.wikipedia.org/wiki/Hash_table) 它会执行虚拟初始化，以使所有可能存在的键都映射到一个字节表示为全零的值。 但是，这种类比并不太恰当，因为它既不能获得映射的所有键的列表，也不能获得所有值的列表。 因此，要么记住你添加到mapping中的数据（使用列表或更高级的数据类型会更好）



合约创建

* ***在内部，构造函数参数在合约代码之后通过 [ABI 编码](https://solidity-cn.readthedocs.io/zh/develop/abi-spec.html#abi) 传递***



## 可见性

* **external ** 外部函数作为合约接口的一部分，意味着我们可以从其他合约和交易中调用，也即意味着你如果在合约内部调用由external修饰的函数，是错误的做法,也就是它是提供给外部使用的，内部是不能调用的。
* ***public*** 函数是合约接口的一部分，可以在内部或通过消息调用。对于公共状态变量， 会自动生成一个 getter 函数

* **internal**  这些函数和状态变量只能是内部访问（即从当前合约内部或从它派生的合约访问），不使用 `this` 调用。（这个修饰符跟java的protected是一样的）

* **private**  函数和状态变量仅在当前定义它们的合约中使用，并且不能被派生合约使用。

***可见性标识符的定义位置，对于状态变量来说是在类型后面，对于函数是在参数列表和返回关键字中间。***

```
/*类型后面*/
uint public data; 

 /*在参数列表和返回关键字中间。*/
 function f(uint a) private pure returns (uint b) { return a + 1; }
```

编译器自动为所有 **public** 状态变量创建 getter 函数,即是同名函数（但你不能显示写同名函数）



函数修饰器

```
// 这个合约只定义一个修饰器，但并未使用： 它将会在派生合约中用到。
    // 修饰器所修饰的函数体会被插入到特殊符号 _; 的位置。
    // 这意味着如果是 owner 调用这个函数，则函数会被执行，否则会抛出异常。
    modifier onlyOwner {
        require(msg.sender == owner);
        _;
    }
```

修饰器modifier 是合约的可继承属性， 并可能被派生合约覆盖。



```
关键字 `payable` 非常重要，否则函数会自动拒绝所有发送给它的以太币。
```

`constant` 是 `view` 的别名。

### Pure 函数

函数可以声明为 `pure` ，在这种情况下，承诺不读取或修改状态。

### Fallback 函数

合约可以有一个未命名的函数。这个函数不能有参数也不能有返回值。 如果在一个到合约的调用中，没有其他函数与给定的函数标识符匹配（或没有提供调用数据），那么这个函数（fallback 函数）会被执行。

为了接收以太币，fallback 函数必须标记为 `payable`。 如果不存在这样的函数，则合约不能通过常规交易接收以太币。

***请确保您在部署合约之前彻底测试您的 fallback 函数，以确保执行成本低于 2300 个 gas。***

一个没有定义 fallback 函数的合约，直接接收以太币（没有函数调用，即使用 `send` 或 `transfer`）会抛出一个异常， 并返还以太币（在 Solidity v0.4.0 之前行为会有所不同）。所以如果你想让你的合约接收以太币，必须实现 fallback 函数。

### 函数重载

合约可以具有多个不同参数的同名函数（除了构造函数）

事件

最多三个参数可以接收 `indexed` 属性，从而使它们可以被搜索

支持多重继承。



## 接口

接口类似于抽象合约，但是它们不能实现任何函数。还有进一步的限制：

1. 无法继承其他合约或接口。
2. 无法定义构造函数。
3. 无法定义变量。
4. 无法定义结构体
5. 无法定义枚举。

安全

智能合约中不要使用随机数。

为了避免重入，你可以使用下面撰写的“检查-生效-交互”（Checks-Effects-Interactions）模式