##Golang基础

###接口值###
   
   从概念上讲一个接口的值，接口值，由两个部分组成，一个具体的类型和那个类型的值。它们被称为接口的动态类型和动态值（可类比java的接口的实现类和实现类实例）。
  
####注意####
 在GO语言中，变量总是被一个定义明确的值初始化，即使接口类型也不例外。对于一个接口的零值就是它的类型和值得部分都是nil,所以遇到接口类型变量的时候，一定要时刻注意接口类型变量的动态类型和动态值（一定是两部分type和value），也就是
 
 * 动态类型 （使用%T来输出动态类型的名称）
 * 动态值

golang中对信号的处理主要使用os/signal包中的两个方法：

* 一个是notify方法，用来监听收到的信号。
* 一个是stop方法，用来取消监听

func Notify(c chan<-os.Signal,sig...os.Signal)

* 第一个参数表示接收信号的管道
* 第二个及后面的参数表示设置要监听的信号，如果不设置表示监听所有的信号

空接口interface{}没有任何方法签名，也就意味着任何类型都实现了空接口。


###交叉编译###

在版本1.6环境下

首先明确的是我们必须找到一个叫make.bash的文件
我们通过运行这样一些命令来看看

```
go env
```
这样我们就得到了go的一些环境变量参数

那么我们就cd到GOROOT目录，在我的mac上这个值是GOROOT=/usr/local/go,所以输入以下命令

```
cd /usr/local/go/src
ls -l
```
这样我们就看到了make.bash这个文件了

以下我们尝试在mac上交叉编译出windows下的exe文件

* cd到make.bash所在的目录下,然后输入以下命令
 
 ```
  sudo GOOS=windows GOARCH=amd64 ./make.bash
 ```

* 上面的命令执行完毕以后，我们先cd到我们需要编译的文件所在目录，假设我们需要编译的文件叫config_generator.go,编译成config_generator.exe,我们输入以下命令

```
GOOS=windows GOARCH=amd64 go build -x -o config_generator.exe config_generator.go
```

这样我们就可以得到config_generator.exe可执行文件了


补充：在sudo GOOS=windows GOARCH=amd64 ./make.bash 命令执行以后我们得到了一个错误

```
##### Building Go bootstrap tool.
cmd/dist
ERROR: Cannot find /Users/yangdamin/go1.4/bin/go.
Set $GOROOT_BOOTSTRAP to a working Go tree >= Go 1.4.
```

这个时候我到官网网站上找到了go1.4.3.darwin-amd64.tar.gz这个压缩包，然后我解压到当前用户目录下，路径是/Users/yangdamin,最后的文件路径是/Users/yangdamin/go1.4,go1.4是我自己命名的(这个路径是固定的，就是提示路径是什么样子的，你就定义成什么样子)，然后修改了./.bash_profile
增加了一行

```
   export GOROOT_BOOTSTRAP=/Users/yangdamin/go1.4
```
之后再次执行了sudo GOOS=windows GOARCH=amd64 ./make.bash 这个命令

这个命令执行完毕没有错误以后，继续上面提到的步骤，最后我们得到config_generator.exe 这个可执行文件。


交叉编译

　　从go1.5开始，将源码编译成非本地系统的程序，不再需要先生成目的系统的编译工具了，直接加参数编译即
　　
　　
```
GOOS=windows GOARCH=amd64 go build
```



###并发编程###

主函数返回时，所有的goroutine都会被直接打断，程序退出，除了从主函数退出或直接终止程序之外，没有其它的编程方法能够让一个goroutine来打断另一个goroutine的执行，如果需要请求一个goroutine结束，必须通过channel来实现，给需要结束的goroutine发送消息，让它自行结束，别的goroutine是不能直接去结束另外一个goroutine的。

服务器编程的套路中，基本来说就是，每一个连接都会起一个goroutine,这样一个goroutine就可以处理一个用户请求，和用户交互，直接连接断开。

创建一个可以发送数据类型T的chan,使用ch:=make(chan T)即可

#####注意#####
<font color="red">一个为nil的channel,在其上进行发送和接收操作都会阻塞。所以channel务必使用make函数进行创建并初始化。channel的零值是nil,channel是一个引用类型，我们复制一个channel或用于函数参数传递时，我们只是拷贝了一个channel引用，因此调用者和被调用者将引用同一个channel。</font>

发送和接收两个操作都是用<-运算符，在发送语句中，<-运算符分割channel和要发送的值。在接收语句中，<-运算符写在channel对象之前。

使用内置的close函数就可以关闭一个channel,close(ch)

对一个已经close过的channel进行发送操作将会导致panic异常
，对一个已经被close过的channel进行接收操作依然可以接收到之前已经成功发送的数据，如果channel中已经没有数据的话，将产生一个零值的数据。


当使用make函数进行channel的创建的时候，如果不指定channel的容量或将容量指定为0，那么这个channel将是无缓冲的，若容量大于0，则这个channel是带缓冲的。


当通过一个无缓冲channel发送数据时，接收者收到数据发生在唤醒发送者之前。也就是说对一个无缓冲的channel进行发送操作之后，这个发送者就进入阻塞状态，等待接收者接收数据并唤醒发送者（happens before）。


并不一定每一个channel都需要关闭，只有当需要告诉接受者goroutine,所有的数据已经全部发送时才需要关闭channel,不管一个channel是否被关闭，当它没有被引用时将会被go语言的垃圾自动回收器回收。（跟打开一个文件不一样）


对一个channel最多只能执行一次close操作，多次执行close操作，将会引发panic。

当然对一个为nil的channel进行任何的操作都会引发panic。


利用<-箭头的方向来确定是发送还是接收，发送 chan<- 也就是从数据的流向来判断，箭头指向chan，说明要把数据放进channel里（形象的说法）。<-chan 表明说要从channel里拿数据出来。


<font color="red">永远只在发送者所在的goroutine中调用close函数。</font>

对带缓存channel进行发送操作就是想内部缓存队列的尾部插入元素，接收操作则是从队列的头部删除元素，如果内部缓存队列是满的，那么发送操作将阻塞直到因另一个goroutine执行接收操作而释放了新的队列空间，相反，如果channel是空的，接收操作将阻塞直到有另一个goroutine执行发送操作而向队列插入元素。


<font color="red">无缓冲channel保证了每个发送操作和接收操作的同步。</font>


<font color="red">一个没有任何case的select语句写作select{},会永远地等待下去</font>



如果多个case同时就绪时（意思是可执行，不会被阻塞），select会随机地选择一个执行，这样来保证每一个channel都有平等的被select的机会。

Go默认是用一个CPU核心的，除非手动设置runtime.GOMAXPROCS。

methon中的receiver应该是用值还是指针，这个时候，请你把receiver当作method的第一个参数来看，你就知道是用值还是指针了。


如果匿名字段实现了一个method，那么包含这个匿名字段的struct也能调用该method。

interface类型

interface类型定义了一组方法，如果某个对象实现了某个接口的所有方法，则此对象就实现了此接口。


如果我们定义了一个interface的变量，那么这个变量里面可以存实现这个interface的任意类型的对象。

Go通过interface实现了duck-typing:即"当看到一只鸟走起来像鸭子、游泳起来像鸭子、叫起来也像鸭子，那么这只鸟就可以被称为鸭子"。

一个函数把interface{}作为参数，那么他可以接受任意类型的值作为参数，如果一个函数返回interface{},那么也就可以返回任意类型的值。
目前有两种方法可以确定interface变量中保存了哪个类型的对象

* Comma-ok断言 
Go语言里面有一个语法，可以直接判断是否是该类型的变量： value, ok = element.(T)，这里value就是变量的值，ok是一个bool类型，element是interface变量，T是断言的类型。如果element里面确实存储了T类型的数值，那么ok返回true，否则返回false。
* switch测试

```
 switch value := element.(type) {
            case int:
                fmt.Printf("list[%d] is an int and its value is %d\n", index, value)
            case string:
                fmt.Printf("list[%d] is a string and its value is %s\n", index, value)
            case Person:
                fmt.Printf("list[%d] is a Person and its value is %s\n", index, value)
            default:
                fmt.Println("list[%d] is of a different type", index)
        }
```

<font color="red">element.(type)语法不能在switch外的任何逻辑里面使用，如果你要在switch外面判断一个类型就使用comma-ok。</font>

runtime.Gosched()表示让CPU把时间片让给别人,下次某个时候继续恢复执行该goroutine。

<font color="red">默认情况下，调度器仅使用单线程，也就是说只实现了并发。想要发挥多核处理器的并行，需要在我们的程序中显式调用 runtime.GOMAXPROCS(n) 告诉调度器同时使用多个线程。GOMAXPROCS 设置了同时运行逻辑代码的系统线程的最大数量，并返回之前的设置。如果n < 1，不会改变当前设置。以后Go的新版本中调度得到改进后，这将被移除</font>

for i := range c能够不断的读取channel里面的数据，直到该channel被显式的关闭。

select默认是阻塞的，只有当监听的channel中有发送或接收可以进行时才会运行，当多个channel都准备好的时候，select是随机的选择一个执行的。

GET请求消息体为空，POST请求带有消息体

GET提交的数据会放在URL之后，以?分割URL和传输数据，参数之间以&相连

POST方法是把提交的数据放在HTTP包的body中

HTTP/1.1协议中定义了5类状态码， 状态码由三位数字组成，第一个数字定义了响应的类别

* 1XX 提示信息 - 表示请求已被成功接收，继续处理
* 2XX 成功 - 表示请求已被成功接收，理解，接受
* 3XX 重定向 - 要完成请求必须进行更进一步的处理
* 4XX 客户端错误 - 请求有语法错误或请求无法实现
* 5XX 服务器端错误 - 服务器未能实现合法的请求


通过下面的代码我们可以看到整个的http处理过程：

```
func (srv *Server) Serve(l net.Listener) error {
    defer l.Close()
    var tempDelay time.Duration // how long to sleep on accept failure
    for {
        rw, e := l.Accept()
        if e != nil {
            if ne, ok := e.(net.Error); ok && ne.Temporary() {
                if tempDelay == 0 {
                    tempDelay = 5 * time.Millisecond
                } else {
                    tempDelay *= 2
                }
                if max := 1 * time.Second; tempDelay > max {
                    tempDelay = max
                }
                log.Printf("http: Accept error: %v; retrying in %v", e, tempDelay)
                time.Sleep(tempDelay)
                continue
            }
            return e
        }
        tempDelay = 0
        c, err := srv.newConn(rw)
        if err != nil {
            continue
        }
        go c.serve()
    }
}
```

监控之后如何接收客户端的请求呢？上面代码执行监控端口之后，调用了srv.Serve(net.Listener)函数，这个函数就是处理接收客户端的请求信息。这个函数里面起了一个for{}，首先通过Listener接收请求，其次创建一个Conn，最后单独开了一个goroutine，把这个请求的数据当做参数扔给这个conn去服务：go c.serve()。这个就是高并发体现了，<font color="red">用户的每一次请求都是在一个新的goroutine去服务，相互不影响。</font>

默认情况下，Handler里面是不会自动解析form的，必须显式的调用r.ParseForm()后，你才能对这个表单数据进行操作。

r.Form里面包含了所有请求的参数，比如URL中query-string、POST的数据、PUT的数据，所有当你在URL的query-string字段和POST冲突时，会保存成一个slice，里面存储了多个值，Go官方文档中说在接下来的版本里面将会把POST、GET这些数据分离开来。
 
注意:<font color="red">
Request本身也提供了FormValue()函数来获取用户提交的参数。如r.Form["username"]也可写成r.FormValue("username")。调用r.FormValue时会自动调用r.ParseForm，所以不必提前调用。r.FormValue只会返回同名参数中的第一个，若参数不存在则返回空字符串。
</font>

###如何有效的防止用户多次递交相同的表单###
解决方案是在表单中添加一个带有唯一值的隐藏字段。在验证表单时，先检查带有该惟一值的表单是否已经递交过了。如果是，拒绝再次递交；如果不是，则处理表单进行逻辑处理。另外，如果是采用了Ajax模式递交表单的话，当表单递交后，通过javascript来禁用表单的递交按钮。

我们上传文件主要三步处理：

* 表单中增加enctype="multipart/form-data"
* 服务端调用r.ParseMultipartForm,把上传的文件存储在内存和临时文件中
* 使用r.FormFile获取文件句柄，然后对文件进行存储等处理。


###客户端上传文件###

```
package main

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "os"
)

func postFile(filename string, targetUrl string) error {
    bodyBuf := &bytes.Buffer{}
    bodyWriter := multipart.NewWriter(bodyBuf)

    //关键的一步操作
    fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
    if err != nil {
        fmt.Println("error writing to buffer")
        return err
    }

    //打开文件句柄操作
    fh, err := os.Open(filename)
    if err != nil {
        fmt.Println("error opening file")
        return err
    }
    defer fh.Close()

    //iocopy
    _, err = io.Copy(fileWriter, fh)
    if err != nil {
        return err
    }

    contentType := bodyWriter.FormDataContentType()
    bodyWriter.Close()

    resp, err := http.Post(targetUrl, contentType, bodyBuf)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    resp_body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    fmt.Println(resp.Status)
    fmt.Println(string(resp_body))
    return nil
}

// sample usage
func main() {
    target_url := "http://localhost:9090/upload"
    filename := "./astaxie.pdf"
    postFile(filename, target_url)
}
```

上面的例子详细展示了客户端如何向服务器上传一个文件的例子，客户端通过multipart.Write把文件的文本流写入一个缓存中，然后调用http的Post方法把缓存传到服务器

注意：<font color="red">如果你还有其他普通字段例如username之类的需要同时写入，那么可以调用multipart的WriteField方法写很多其他类似的字段。</font>


###服务器端处理文件上传###

```
http.HandleFunc("/upload", upload)

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method) //获取请求的方法
    if r.Method == "GET" {
        crutime := time.Now().Unix()
        h := md5.New()
        io.WriteString(h, strconv.FormatInt(crutime, 10))
        token := fmt.Sprintf("%x", h.Sum(nil))

        t, _ := template.ParseFiles("upload.gtpl")
        t.Execute(w, token)
    } else {
        r.ParseMultipartForm(32 << 20)
        file, handler, err := r.FormFile("uploadfile")
        if err != nil {
            fmt.Println(err)
            return
        }
        defer file.Close()
        fmt.Fprintf(w, "%v", handler.Header)
        f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)  // 此处假设当前目录下已存在test目录
        if err != nil {
            fmt.Println(err)
            return
        }
        defer f.Close()
        io.Copy(f, file)
    }
}
```

通过上面的代码可以看到，处理文件上传我们需要调用r.ParseMultipartForm，里面的参数表示maxMemory，调用ParseMultipartForm之后，上传的文件存储在maxMemory大小的内存里面，如果文件大小超过了maxMemory，那么剩下的部分将存储在系统的临时文件中。我们可以通过r.FormFile获取上面的文件句柄，然后实例中使用了io.Copy来存储文件

注意:<font color="red">
获取其他非文件字段信息的时候就不需要调用r.ParseForm，因为在需要的时候Go自动会去调用。而且ParseMultipartForm调用一次之后，后面再次调用不会再有效果.
</font>

<font color="red">注意：请详细阅读 database/sql 源码</font>


比较好的mysql驱动

```
https://github.com/go-sql-driver/mysql 支持database/sql，全部采用go写
```

###cookie###

客户端发起一个请求

服务器端接受请求，并响应，响应的过程中将cookie发送给客户端

客户端接收到服务器端的响应，并把cookie保存在本地

客户端再发起一个请求，这个请求会携带cookie发送给服务器

Cookie是由浏览器维持的，存储在客户端的一小段文本信息，伴随着用户请求和页面在Web服务器和浏览器之间传递。用户每次访问站点时，Web应用程序都可以读取cookie包含的信息。浏览器设置里面有cookie隐私数据选项，打开它，可以看到很多已访问网站的cookies

cookie是有时间限制的，根据生命期不同分成两种：会话cookie和持久cookie；

如果不设置过期时间，则表示这个cookie生命周期为从创建到浏览器关闭止，只要关闭浏览器窗口，cookie就消失了。这种生命期为浏览会话期的cookie被称为会话cookie。会话cookie一般不保存在硬盘上而是保存在内存里。

如果设置了过期时间(setMaxAge(606024))，浏览器就会把cookie保存到硬盘上，关闭后再次打开浏览器，这些cookie依然有效直到超过设定的过期时间。存储在硬盘上的cookie可以在不同的浏览器进程间共享，比如两个IE窗口。而对于保存在内存的cookie，不同的浏览器有不同的处理方式。 

Go设置cookie

Go语言中通过net/http包中的SetCookie来设置
```
http.SetCookie(w ResponseWriter, cookie *Cookie)
```

<font color="red">一定要注意，这里是服务器端向客户端写入cookie</font>,也就是将cookie保存在客户端中

Go读取cookie

```
   cookie, _ := r.Cookie("username")//读取一个
   
   for _, cookie := range r.Cookies() {//读取多个
   
   //r 为http.Request
   
```


###session###

session，简而言之就是在服务器上保存用户操作的历史信息。服务器使用session id来标识session，session id由服务器负责产生，保证随机性与唯一性，相当于一个随机密钥，避免在握手或传输中暴露用户真实密码。但该方式下，仍然需要将发送请求的客户端与session进行对应，所以可以借助cookie机制来获取客户端的标识（即session id），也可以通过GET方式将id提交给服务器。

cookie和session都是因为http是无状态的，两次请求之间是没有任何关联的而采取的方式
都是用一些唯一的值来回传输来确认是同一个用户

比如当使用session的时候，注意，在不能使用cookie的情况下，我们就必须使用某一个变量来保存sessionId,并且把这个值发送到客户端，客户端再发起请求的时候，需要将这个变量提交给服务器，服务器将这个变量在下发给客户端的过程中是保存在服务器端的，当提交上来的时候就跟这个sessionId做一下关联(也就是查询有没有这个sessionId,如果有，则使用它来取得用户的信息)，这个sessionId里映射很多的信息的，这个由服务器端自己来做，比如可以将一个用户的全部信息关联在sessionId上，map[sessionId]UserInfo

session，中文经常翻译为会话，其本来的含义是指有始有终的一系列动作/消息，比如打电话是从拿起电话拨号到挂断电话这中间的一系列过程可以称之为一个session。然而当session一词与网络协议相关联时，它又往往隐含了“面向连接”和/或“保持状态”这样两个含义。

session在Web开发环境下的语义又有了新的扩展，它的含义是指一类用来在客户端与服务器端之间保持状态的解决方案。有时候Session也用来指这种解决方案的存储结构。

session机制是一种服务器端的机制，服务器使用一种类似于散列表的结构(也可能就是使用散列表)来保存信息。

但程序需要为某个客户端的请求创建一个session的时候，服务器首先检查这个客户端的请求里是否包含了一个session标识－称为session id，如果已经包含一个session id则说明以前已经为此客户创建过session，服务器就按照session id把这个session检索出来使用(如果检索不到，可能会新建一个，这种情况可能出现在服务端已经删除了该用户对应的session对象，但用户人为地在请求的URL后面附加上一个JSESSION的参数)。如果客户请求不包含session id，则为此客户创建一个session并且同时生成一个与此session相关联的session id，这个session id将在本次响应中返回给客户端保存。

session机制本身并不复杂，然而其实现和配置上的灵活性却使得具体情况复杂多变。这也要求我们不能把仅仅某一次的经验或者某一个浏览器，服务器的经验当作普遍适用的。

session和cookie的目的相同，都是为了克服http协议无状态的缺陷，但完成的方法不同。session通过cookie，在客户端保存session id，而将用户的其他会话消息保存在服务端的session对象中，与此相对的，cookie需要将所有信息都保存在客户端。因此cookie存在着一定的安全隐患，例如本地cookie中保存的用户名密码被破译，或cookie被其他网站收集

目前Go标准包没有为session提供任何支持，我们需要自己动手来实现go版本的session管理和创建。

session的基本原理是由服务器为每个会话维护一份信息数据，客户端和服务端依靠一个全局唯一的标识来访问这份数据，以达到交互的目的。当用户访问Web应用时，服务端程序会随需要创建session，这个过程可以概括为三个步骤：

* 生成全局唯一标识符（sessionid）
* 开辟数据存储空间。一般会在内存中创建相应的数据结构，但这种情况下，系统一旦掉电，所有的会话数据就会丢失，如果是电子商务类网站，这将造成严重的后果。所以为了解决这类问题，你可以将会话数据写到文件里或存储在数据库中，当然这样会增加I/O开销，但是它可以实现某种程度的session持久化，也更有利于session的共享
* 将session的全局唯一标示符发送给客户端。

以上三个步骤中，最关键的是如何发送这个session的唯一标识这一步上。考虑到HTTP协议的定义，数据无非可以放到请求行、头域或Body里，所以一般来说会有两种常用的方式：cookie和URL重写。

* Cookie 服务端通过设置Set-cookie头就可以将session的标识符传送到客户端，而客户端此后的每一次请求都会带上这个标识符，另外一般包含session信息的cookie会将失效时间设置为0(会话cookie)，即浏览器进程有效时间。至于浏览器怎么处理这个0，每个浏览器都有自己的方案，但差别都不会太大(一般体现在新建浏览器窗口的时候)；
* URL重写 所谓URL重写，就是在返回给用户的页面里的所有的URL后面追加session标识符，这样用户在收到响应之后，无论点击响应页面里的哪个链接或提交表单，都会自动带上session标识符，从而就实现了会话的保持。虽然这种做法比较麻烦，但是，如果客户端禁用了cookie的话，此种方案将会是首选。


##目录操作##

文件操作的大多数函数都是在os包里面

* func Mkdir(name string, perm FileMode) error

	创建名称为name的目录，权限设置是perm，例如0777

* func MkdirAll(path string, perm FileMode) error

  根据path创建多级子目录，例如astaxie/test1/test2。
  
* func Remove(name string) error

 删除名称为name的目录，当目录下有文件或者其他目录是会出错
 
* func RemoveAll(path string) error
  
  根据path删除多级子目录，如果path是单个名称，那么该目录下的子目录全部删除。
  
##文件操作##

新建文件可以通过如下两个方法

* func Create(name string) (file *File, err Error)

	根据提供的文件名创建新的文件，返回一个文件对象，默认权限是0666的文件，
	返回的文件对象是可读写的。

* func NewFile(fd uintptr, name string) *File

	根据文件描述符创建相应的文件，返回一个文件对象

通过如下两个方法来打开文件：

* func Open(name string) (file *File, err Error)

	该方法打开一个名称为name的文件，但是是只读方式，内部实现其实调用了OpenFile。

* func OpenFile(name string, flag int, perm uint32) (file *File, err Error)

	打开名称为name的文件，flag是打开的方式，只读、读写等，perm是权限

写文件

写文件函数：

* func (file *File) Write(b []byte) (n int, err Error)

	写入byte类型的信息到文件

* func (file *File) WriteAt(b []byte, off int64) (n int, err Error)

	在指定位置开始写入byte类型的信息

* func (file *File) WriteString(s string) (ret int, err Error)

	写入string信息到文件
	
<font color="red">Go语言里面删除文件和删除文件夹是同一个函数</font>

##单元测试##

```
文件名必须是*_test.go的类型，*代表要测试的文件名，函数名必须以Test开头,也即是类似如下函数命名

func Testxxx(t *testing.T)

xxx为要测试的函数

在需要测试的目录下输入go test命令，将对当前目录下的所有*_test.go文件进行编译并自动运行测试。

*_test.go和*.go文件需同包和同目录

```

<font color="red">不检查close()的返回值是一个常见但非常严重的编程错误,所以务必要检查close函数的返回值。</font>



























