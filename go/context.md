##Context##

1、不要将Context存储在结构类型中，而是将Context显式地传递给需要它的每个函数，并且将Context作为第一个参数。不要传递一个值为nil的Context给函数。尽管这个函数允许你这么做，你也不要将一个值为nil的Context传递给函数，如果你不确定要使用那一个Context,那么就使用context.TODO这个Context传递给函数。

2、<font color="red">Context是线程安全的，多个线程可以同时访问它。</font>

总结一下就是

* 不要将 Contexts 放入结构体，相反 context 应该作为第一个参数传入，命名为 ctx 。 func DoSomething（ctx context.Context，arg Arg）error { // ... use ctx ... }

* 即使函数允许，也不要传入 nil 的 Context。如果不知道用哪种 Context，可以使用 context.TODO() 。
* 使用context的Value相关方法只应该用于在程序和接口中传递的和请求相关的元数据，不要用它来传递一些可选的参数
* 相同的 Context 可以传递给在不同的 goroutine ；Context 是并发安全的。


Context 的调用应该是链式的，通过 WithCancel ， WithDeadline ， WithTimeout 或 WithValue 派生出新的 Context。当父 Context 被取消时，其派生的所有 Context 都将取消。

通过 context.WithXXX 都将返回新的 Context 和 CancelFunc。调用 CancelFunc 将取消子代，移除父代对子代的引用，并且停止所有定时器。未能调用 CancelFunc 将泄漏子代，直到父代被取消或定时器触发

Context是一个interface{}

Done()，返回一个channel。当times out或者调用cancel方法时，将会close掉。
Err()，返回一个错误。该context为什么被取消掉。
Deadline()，返回截止时间和ok。
Value()，返回值。

在context包内部已经为我们实现好了两个空的Context，可以通过调用Background()和TODO()方法获取。一般的将它们作为Context的根，往下派生。

上下文，自然就包括上文和下文，上文传递给下文