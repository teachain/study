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









# go modules

想要了解go mod 怎么使用，我们的第一反应应该是在命令行里输入以下内容：

```
go mod help
```

我们来尝试使用go mod来进行包管理，首先建立一个项目的根目录

```
mkdir test
```

接着使用命令初始化一个项目

```
go mod init test
```

执行完这个命令之后，会在test目录下生成一个go.mod文件，打开这个文件

```
vim go.mod
```

我们会发现它的内容是

```
module test
```

所以我们得到一个结论，如果我们不想执行go mod init命令，我们其实可以直接在项目的根目录下新建一个

go.mod文件，然后写上类似module test这样的内容即可。

注意点：

1. 依然使用原来的import方式，也就是说如果在test目录下新建了一个目录core,core里定义了一个函数func 

   Println(str string){fmt.Println(str)},那么在test目录下的main.go文件里调用core中的Println函数

   那么需要按原来使用gopath的时候一样的方式导入，import "test/core"

2. 当项目需要使用第三方的时候，例如 import "github.com/gin-gonic/gin" 这个仍然按gopath的方式进行，这个时候我们并不需要事先把github.com/gin-gonic/gin进行下载下来，当我们执行了go build,go run,go install的时候，go mod会帮助我们把响应的依赖下载下来，也就是当我们使用go build，go install ,go test以及go list时，go会自动得更新go.mod文件，将依赖关系写入其中,并且生成一个go.sum文件。

3. 当采用了go mod的方式来管理包的时候，我们不想再次去下载依赖，我想直接把这个项目以及依赖都发给我的同事，让他能够直接编译。那么我们可以在项目的根目录下执行 go mod vendor 命令，go mod就会在根目录下生成vendor目录，项目的依赖就全部在这个目录下了。

4. modules和传统的GOPATH不同，不需要包含例如src，bin这样的子目录，一个源代码目录甚至是空目录都可以作为module，只要其中包含有go.mod文件。

5. go.mod文件里只会记录项目里直接进行import的哪些项目依赖和版本，至于这么依赖项目中的依赖项目是不写入这个文件的。

6. 当在一个不存在go.mod文件源码根目录下，存在不为空的vendor目录的时候，执行go mod init 项目名 这个命令的时候，go mod会将vendor中所有的依赖都会写入到go.mod文件中，当然都是包括版本号的。

7. 当从vendor中写入go.mod的时候，可能会引起错误，这个时候最彻底的解决方案就是直接删掉vendor目录，让go mod来分析依赖包并下载依赖包。

8. 实际的操作是go mod把依赖项以带版本号的方式存放依赖项的源码的，统一放在了gopath/pkg/mod目录下。

#### 更新依赖的版本

```
go mod edit -require="github.com/chromedp/chromedp@v0.1.0" #依赖项目的链接@版本号
go mod tidy   #add missing and remove unused modules
```



##### 



  