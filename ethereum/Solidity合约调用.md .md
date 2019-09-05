在 Solidity 中，call 函数簇可以实现跨合约的函数调用功能，其中包括 call、delegatecall 和 callcode 三种方式。

* call

  ```
  <address>.call(...) returns (bool)
  ```

* delegatecall

  ```
  <address>.delegatecall(...) returns (bool)
  ```

* callcode

  ```
  <address>.callcode(...) returns (bool)
  ```

  

三种调用方式的异同点

- call: 最常用的调用方式，调用后内置变量 msg 的值会修改为调用者，执行环境为被调用者的运行环境(合约的 storage)。
- delegatecall: 调用后内置变量 msg 的值不会修改为调用者，但执行环境为调用者的运行环境。
- callcode: 调用后内置变量 msg 的值会修改为调用者，但执行环境为调用者的运行环境。





### fallback函数（没名没参没返回）

回退函数是合约里的特殊函数，没有名字，没有参数，没有返回值。当调用的函数找不到时，就会调用默认的fallback函数，这是个真真正正的三无产品。

```
 function() external{
       //todo
 }
```

无名字，无参数，无返回值，必须用external修饰。