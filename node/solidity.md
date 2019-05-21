# 合约结构

- 状态变量
- 函数
- 事件Event
- 结构体
- 枚举类型



#### 状态变量

状态变量是永久地存储在合约存储中的值。

它的可见性有三:

- public
- internal
- private

默认为internal



#### 函数

函数是合约中代码的可执行单元。

它的可见性有四：

- external
- public
- internal
- private

默认为public

函数修饰器

```
function (<parameter types>) {internal|external} [pure|constant|view|payable] [returns (<return types>)]


pragma solidity >=0.4.22 <0.7.0;

contract MyPurchase {
    address public seller;

    modifier onlySeller() { // 修饰器
        require(
            msg.sender == seller,
            "Only seller can call this."
        );
        _;
    }

    function abort() public onlySeller { // 修饰器用法
        // ...
    }
}
```



#### 事件Event

```
pragma solidity >=0.4.21 <0.7.0;
contract TinyAuction {
    event HighestBidIncreased(address bidder, uint amount); // 定义事件

    function bid() public payable {
        // ...
        emit HighestBidIncreased(msg.sender, msg.value); // 触发事件
    }
}
```



#### 结构体

```
pragma solidity >=0.4.0 <0.7.0;

contract TinyBallot {
    struct Voter { // 结构体
        uint weight;
        bool voted;
        address delegate;
        uint vote;
    }
}
```



#### 枚举

```
pragma solidity >=0.4.0 <0.7.0;

contract Upchain {
    enum State { Created, Locked, InValid } // 枚举
}
```



#### 整型

int / uint ：分别表示有符号和无符号的不同位数的整型变量。 支持关键字 uint8 到 uint256 （无符号，从 8 位到 256 位）以及 int8 到 int256，以 8 位为步长递增。 uint 和 int 分别是 uint256 和 int256 的别名（特别关注）。