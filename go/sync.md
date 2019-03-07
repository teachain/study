##sync包源码阅读##

我们来看一下最简单的一个类型sync.Once的源码（我把注释去掉了）

```
package sync

import (
	"sync/atomic"
)

type Once struct {
	m    Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

```
就正如它的命名一样，once,也就是只执行一次，意思是说针对<font color="red">每一个sync.Once,它的唯一的的方法，Do(f func())，无论你调用它多少次，它都只会执行一次。</font>

从Do的源码里我们看到了它用了一个uint32的变量来控制的，而且他依赖于原子操作和借用锁（Mutext）里实现的。

* 先是原子读取done这个变量，判断这个变量如果已经等于1，则直接返回了
* 如果变量不为1，那么先取得锁，然后调用参数指定的函数，函数执行完之后，原子地将done设置为1

从源码可以看出这样的一个类型，我们借助于Mutex和atomic都是可以分分钟实现出来的。



    golang的goroutine机制有点像线程池：
        一、go 内部有三个对象： P对象(processor) 代表上下文（或者可以认为是cpu），M(work thread)代表工作线程，G对象（goroutine）.
        二、正常情况下一个cpu对象启一个工作线程对象，线程去检查并执行goroutine对象。碰到goroutine对象阻塞的时候，会启动一个新的工作线程，以充分利用cpu资源。所以有时候线程对象会比处理器对象多很多。