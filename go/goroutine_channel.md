##并发##

Go语言红的并发程序可以用两种手段来实现：

* goroutine和channel  (顺序通信进程 communicating sequential processes)

* 多线程共享内存


在Go语言中，每一个并发的执行单元叫做一个goroutine。

当一个程序启动时，其主函数即在一个单独的goroutine中运行，我们叫它main goroutine。

新的goroutine会用go语言来创建。在语法上，go语句是一个普通的函数或方法调用前加上关键字go。go语句会使其语句中的函数在一个新创建的goroutine中运行。而go语句本身会迅速完成。（<font color="red">go语句本身执行完毕，goroutine的函数并不是马上就执行了，只是告诉系统，创建一个goroutine,加入到调度队列中，有空就调度一下，让我执行一下呗。是这个意思。</font>）看以下实例：

```
f() //call f();wait for it to return

go f() // create a new goroutine that calls f();don't wait

```
<font color="red">主函数返回时，所有的goroutine都会直接被打断，程序退出。除了从主函数退出或者直接终止程序之外，没有其它的编程方法能够让goroutine A 来打断goroutine B的执行。（意思是没有直接地调用函数或方法来打断，但可以通过一种通信机制来实现一个goroutine打断另外一个goroutine的目的，通过channel通信，让goroutine自行结束执行。）</font>

使用内置的make函数创建channel（可以创建有缓冲和无缓冲的channel）

```
ch:=make(chan T)

```
T为具体的数据类型，比如int ,string等等

<font color="red">一个channel对应一个由make创建的底层数据结构的引用。当我们复制一个channel或用于函数参数传递时，我们只是拷贝了一个channel引用，因此调用者和被调用者将引用同一个channel对象。channel的零值是nil。使用内置的close函数就可以关闭一个channel。</font>


一个基于<font color="red">无缓存channels</font>的发送操作将导致发送者所在的goroutine阻塞，直到另一个goroutine在相同的channels上执行接收操作，当发送的值通过channels成功传输之后，两个goroutine可以继续执行后面的语句。反之，如果接收操作先发生，那么接受者goroutine也将阻塞，直到有另一个goroutine在相同的channels上执行发送操作。基于无缓存channels的发送和接收操作将导致两个goroutine做一次同步操作。因为这个原因，无缓存channels有时候也被称为同步channels。

<font color="red">核心：当通过一个无缓存channels发送数据时，接收者收到数据发生在唤醒发送者goroutine之前。</font>也就是发送者在接受者没有收到数据之前阻塞睡眠了，等到接受者接收到数据之后，接收者就唤醒发送者，说你可以继续了。

当一个channel作为一个<font color="red">函数参数</font>时，它一般总是被专门用于<font color="red">只发送</font>或者<font color="red">只接收</font>。


<font color="red">用一个channel的关闭事件广播开来，结束其他goroutine的执行（这些goroutine里都有select 这个channel）。</font>

<font color="red">无论任何时候，只要有两个goroutine并发访问同一变量，且至少其中的一个是写操作的时候就会发生数据竞争。</font>

根本就没有什么所谓的良性数据竞争，我们一定要避免数据竞争。


* 我们可以用一个容量只有1的channel来保证最多只有一个goroutine在同一时刻访问一个共享变量。一个只能为1和0的信号量叫做二元信号量。

```
  //当然sema是全局的
  var sema=make(chan struct{},1)
  sema<-struct{}{} //相当于互斥锁的lock
  //受保护的代码区
  <-sema //相当于互斥锁的unlock
  
```
  
<font color="red">所有并发的问题都可以用一致的、简单的既定的模式来规避。所以可能的话，将变量限定在goroutine内部；如果是多个goroutine都需要访问的变量，使用互斥条件来访问。</font>


##生产者消费者模式##

并发编程最常见的就是生产者消费者模式,该模式通过平衡生产线程和消费线程的工作能力来提高程序的整体处理数据的速度。简单地说，就是生产者生产一些数据，然后放到buffer中，同时消费者从buffer中来取这些数据。这样就让生产消费变成了异步的两个过程。当buffer中没有数据时，消费者就进入等待过程；而当buffer中数据已满时，生产者则需要等待buffer中数据被取出后再写入。


###导入包###

<font color="red">导入包的重命名只影响当前的源文件。</font>


###互斥锁###

互斥锁是传统的并发程序对共享资源是进行访问控制的主要手段。它由标准库代码包sync中的Mutex结构体类型代表。<font color="red">sync.Mutex类型（确切地说，*sync.Mutex类型）只有两个公开方法--Lock和Unlock</font>。顾名思义，前者被用于锁定当前的互斥量，而后来则用来对当前的互斥量进行解锁。
<font color="red">类型sync.Mutex的零值表示了未被锁定的互斥量。也就是说它是一个开箱即用的工具。Lock和Unlock一定要成对出现。</font>

建议:

* 把对同一个互斥锁的成对的锁定和解锁操作放在同一个层次的代码块中。
* 在同一个函数或方法中对某个互斥量进行锁定和解锁。
* 把互斥量作为某一个结构体类型中的字段，以便在该类型的多个方法中使用它。
* 应该使代表互斥锁的变量的访问权限尽量地低。这样才能尽量避免它在不相关的流程中被吴用。从而导致程序不正确的行为。

###读写锁###

在Go语言中，读写锁由结构体类型sync.RWMutex代表。与互斥锁类型，sync.RWMutex而理性的零值就已经是立即可用的读写锁了。它包含了两对方法Lock和unLock,RLock和RUnlock,Lock和unLock代表了对写操作的锁定和解锁。RLock和RUnlock表示了对读操作的锁定和解锁。
也即是说写的时候，不能读，读的时候不能写，但是一个goroutine在读，另外一个goroutine也可以读，但不能写，也就是允许单写多读。


##chan##

<font color="red">针对通道的操作本身是同步的。</font>

在同一时刻，仅有一个Goroutine能向一个通道发送元素值，同时也仅有一个Goroutine能从它那里接收元素值。在通道中，各个元素值都是严格按照被发送至此的先后顺序排列的，最早被发送至通道的元素值会被最先接收。也就是先进先出。通道中的元素值都具有原子性。它们是不可分割的。通道中的每个元素值都只可能被某一个Goroutine接收。已被接收的元素值会立刻被从通道中删除。

```
   chan1:=make(chan int,1)
   
   chan2:=make(chan int,0)
   
```

上面的代码中，chan1是带缓冲的，也就是当chan1中还不存在元素的时候，可以往里面写入一个元素成功，并且不阻塞，然后到（在未读出之前）打算写入第二个元素的时候，会阻塞，直到有一个Goroutine将元素读出。<font color="red">（可以把它当锁来使用）</font>

而chan2却不同，当chan1中还不存在元素的时候，可以往里面写入一个元素会被阻塞，直到有一个Goroutine将元素读出，才唤醒当前的chan2所在的Goroutine。

如下示例

```
package main

import (
	"fmt"
)

func main() {

	fmt.Println("go start......")

	chan1 := make(chan int, 1)

	chan1 <- 1

	fmt.Println("prepare to receive.....")

	<-chan1
}

以下是它的输出
go start......
prepare to receive.....

```

```

package main

import (
	"fmt"
)

func main() {

	fmt.Println("go start......")

	chan1 := make(chan int, 0)

	chan1 <- 1

	fmt.Println("prepare to receive.....")

	<-chan1
}

以下是它的运行结果
go start......
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
......

```


<font color="red">试图从一个未被初始化的通道类型（即值为nil的通道类型的变量）那里接收元素值或向通道中发送元素值会造成当前Goroutine的永久阻塞。也就是说chan必须进行make。</font>

<font color="red">试图向一个已被关闭的通道发送元素值，那么会立即引发一个运行时恐慌。</font>

在select代码块中进行发送操作。


<font color="red">最先被阻塞的那个Goroutine会被最先被唤醒。在我们向通道发送一个值之后，该通道将会得到该值的一个副本，而非该值本身。当这个副本形成之后，我们对那个原来的值的任何修改都不会影响到通道中相应的副本。</font>

在select语句中，每个分支依然以关键字case开始，跟在每个case后面的只能是针对某个通道的发送语句或接收语句。在开始执行select语句的时候，所有跟在case关键字右边的发送语句或接收语句中的通道表达式和元素表达式都会先被求值，求值的顺序是自上而下，从左到右的。无论它们所在的case是否有可能被选择都会是这样。


当我们将for和select联合使用时，要注意的是如果有break,一定要注意，当在select中使用break的时候，只会退出当前的select,而没有退出for,所以要在select块以外再次使用break来退出for。

<font color="red">不正确的对channel进行close操作是引起运行时恐慌的源头</font>

1）、对一个已经关闭的channel再次close,会引起运行时恐慌。

2）、对一个值为nil的channal进行close操作，会引起运行时恐慌。

3）、向一个已经关闭的channel发送元素值，会引起运行时恐慌。


































