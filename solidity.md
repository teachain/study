智能合约的部署和调用都需要花费gas

solidity一旦出现异常，所有的执行都将会被回撤，这主要是为了保证合约执行的原子性，以避免中间状态出现的数据不一致。

智能合约部署的时候，会得到一个contract address,例如

```
0x08970fed061e7747cd9a38d680a601510cb659fb
```

并且消耗若干gas

当调用智能合约的时候，to将会是智能合约的地址

```
0x08970fed061e7747cd9a38d680a601510cb659fb
```

这个操作也是需要消耗若干gas的。

decoded input 是合约调用的输入参数。

decoded output是合约调用的函数返回数据。

- 合约类似面向对象语言中的类。
- 支持继承

每个合约中可包含

1. 状态变量(State Variables)

2. 函数(Functions)

3. 函数修饰符（Function Modifiers）

4. 事件（Events）

5. 结构类型(Structs Types)

6. 枚举类型(Enum Types)

   

   #### 状态变量

    变量值会永久存储在合约的存储空间

   ### 函数（Functions）

   智能合约中的一个可执行单元。函数上增加payable标识，即可接收ether，并会把ether存在当前合约



# 整数

整数所占位数可以指定从`uint8/int8`到`uint256/int256`，以8为步长单位递增的不同长度。`uint`和`int`默认表示`uint256`和`int256`。