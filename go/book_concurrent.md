##Go并发编程##

无论什么时候，只要有两个goroutine并发访问同一个变量，且至少其中的一个是写操作的时候，就是发生数据竞争。

一个好的经验法则是根本就没有什么所谓的良性数据竞争，所以我们一定要避免数据竞争。

避免数据竞争一般有三种方法：

* 不要去写变量（实际上就是把数据竞争产生的条件去掉了）

* 避免从多个goroutine访问变量（实际上就是把数据竞争产生的条件去掉了） 

* 允许多个goroutine去访问变量，但是在同一时刻最多只有一个goroutine在访问，这种方法被称为"<font color="red">互斥</font>"


在Go中，互斥可以多种实现方式，可以用sync.Mutex互斥锁来实现，也可以用一个容量为1的channel来保证最多只有一个goroutine在同一时刻访问一个共享变量。
比如

```
  var locker=make(chan struct{},1)

  locker<-struct{}
  
  //临界区
  
  <-locker
  
```


几乎类似于sync.Mutex。


sync.RWMutex是多读单写锁，其允许多个只读操作并行执行，但写操作会完全互斥，有就是说读的时候，多个goroutine并行读但不能写，

读的时候不能写，写的时候不能读

所有并发的问题都可以用一致的，简单的既定的模式来规避。所以可能的话，将变量限制在goroutine内部，如果是多个goroutine都需要访问的变量，使用互斥条件来访问。


<font color="red">Go的调度器使用了一个叫做GOMAXPROCS的变量来决定会有多少个操作系统的线程同时执行Go的代码。其默认的值是运行时机器上的CPU的核心数，所以在一个有8个核心的机器上时，调度器一次会在8个OS线程上去调度Go代码。</font>


###条件变量###

在Go语言中，sync.Cond类型代表了条件变量。创建一个条件变量只有一个方法<font color="red">（必须使用该方法创建）</font>

```
	func NewCond(l Locker) * Cond

```

<font color="red">条件变量总是要与互斥量组合使用</font>

从这个函数的参数Locker的定义来看

```

	type Locker interface {
        Lock()
        Unlock()
   }
   
```

因为互斥锁和读写锁都实现了该接口，所以条件变量的创建的函数的参数具体类型可以是互斥锁sync.Mutext也可以是sync.RWMutex

类型*sync.Cond的方法集合中有3个方法

* <font color="red">Wait方法（等待通知)</font>
* <font color="red">Signal方法（单发通知）</font>
* <font color="red">Broadcast方法（广播通知）</font>


调用Wait方法会自动地对该条件变量关联的那个锁（NewCond的参数）进行解锁（<font color="red">也就是调用条件变量的Wait方法之前，必定需要先对相关的锁进行锁定，否则调用Wait方法就会引发恐慌--对一个没有加锁的互斥锁进行解锁操作，会引起恐慌。不用程序员自己去调用UnLock方法,Wait方法里包含了</font>），并且使得调用方所在的Goroutine被阻塞（通俗点说，调用Wait方法后的代码先暂时不执行了，等着被唤醒吧）。

一旦该方法（Wait方法）收到通知，就会试图再次锁定该锁（<font color="red">Wait方法里包含了这个功能，不用程序员自己去调用Lock方法</font>），再次锁定该锁有两种可能：

* 锁定成功，它就会唤醒那个被它阻塞的Goroutine(Wait方法后面的代码得以继续执行)
* 锁定失败，它会等待下一个通知，被它阻塞的Goroutine继续被阻塞。（继续等）


<font color="red">注意：

* 调用Wait方法之前，必须锁定与此条件变量相关联的锁，当然解锁操作也是配套的。

* 调用Signal和Broadcast方法都不需要锁定和解锁相关联的锁（与Wait方法不同）

</font>

惯用法

```
	c.L.Lock()
	for !condition() {
	    c.Wait()
	}
	... make use of condition ...
	c.L.Unlock()
	
```















