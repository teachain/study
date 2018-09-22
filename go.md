一些有用的库

```
github.com/rjeczalik/notify #文件夹文件变化监听库
github.com/chain/chain  #金融项目的chain
github.com/hashicorp/golang-lru #Least Recently Used，即最近最少使用算法。
```

go语言自动补全代码，需要添加gocode的程序。

```
go get github.com/nsf/gocode
```

排空channel的技巧

```
done:
	// Empty work channel
	for {
		select {
		case <-self.workCh:
		default:
			break done
		}
	}
```

只要记住箭头方向，指向chan（关键字）表示往chan里写，不指向chan表示从chan里读取。

```
//只能向chan里写数据
func send(c chan<- int) {
    for i := 0; i < 10; i++ {
        c <- i
    }
}
//只能取channel中的数据
func recv(c <-chan int) {
    for i := range c {
        fmt.Println(i)
    }
}
```



#### Timer

Timer是一次性的时间触发事件。这一点与Ticker不同，Ticker是按一定时间间隔持续触发时间事件。这一点要非常明确，Timer是一次性触发事件，要重用则必须进行Reset.

Timer三种创建姿势:

```
t:= time.NewTimer(d) //这个创建的timer是看得见的,可调用Reset来重用，也可以调用Stop进行取消。
t:= time.AfterFunc(d, f)//这里的t只能调用Stop方法了，这个创建的timer是看得见的，但是它不能重用
c:= time.After(d) //func After(d Duration) <-chan Time,这个创建的timer是隐形的，你看不到
```

一定要反复强调，它是一次性触发事件。



### 反射

**反射有两个问题，在使用前需要三思：**

1. 大量的使用反射会损失一定性能
2. Clear is better than clever. Reflection is never clear.

**Go的类型设计上有一些基本原则，理解这些基本原则会有助于你理解反射的本质：**

1. 变量包括 <type, value> 两部分。理解这一点你就知道为什么`nil != nil`了。
2. type包括 `static type`和`concrete type`. 简单来说 `static type`是你在编码是看见的类型，`concrete type`是runtime系统看见的类型。
3. 类型断言能否成功，取决于变量的`concrete type`，而不是`static type`. 因此，一个 reader变量如果它的concrete type也实现了write方法的话，它也可以被类型断言为writer.
4. Go中的反射依靠`interface{}`作为桥梁，因此遵循原则3. 例如，反射包.Kind方法返回的是`concrete type`, 而不是`static type`.

- Kind() Kind
  Kind返回该接口的具体分类
  Kind代表Type类型值表示的具体分类。零值表示非法分类。

  ```
  const (
      Invalid Kind = iota
      Bool
      Int
      Int8
      Int16
      Int32
      Int64
      Uint
      Uint8
      Uint16
      Uint32
      Uint64
      Uintptr
      Float32
      Float64
      Complex64
      Complex128
      Array
      Chan
      Func
      Interface
      Map
      Ptr
      Slice
      String
      Struct
      UnsafePointer
  )
  ```



## sync.WaitGroup作为参数传入函数中，一定要传指针