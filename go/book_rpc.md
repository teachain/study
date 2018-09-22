##RPC##

RPC模型也有两种，PRC via TCP或者RPC via HTTP。RPC via TCP相较于后者传输的数据更少（少了HTTP包的那些部分），所以效率更高，但是却实际上不适合在Internet这样存在着危险的环境中使用。所以RPC via TCP适合在内网的两个进程间交互使用。

###golang内置的rpc###

golang 内置了三种rpc,你可以根据你项目的需求选用其中一种

####1、TCP##
我们先从基本的tcp的rpc开始讲起，下面是服务器端

```
import (
	"fmt"
	"net"
	"net/rpc"
)

//1、定义服务
type EchoServer bool

//2、定义服务的方法method
//参数和返回可以是基本类型或是自定义类型都可以，但必须是指针，并且返回一个error
func (s *EchoServer) Echo(req *string, res *string) error {
	//要注意指针操作，如果res = req，你会发现客户端得不到结果
	*res = *req
	return nil
}

func main() {
	//3、new服务
	echo := new(EchoServer)
	//4、注册服务
	if err := rpc.Register(echo); err != nil {
		fmt.Println(err)
		return
	}
	addr := ":8080"
	var listener *net.TCPListener
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		//5、用rpc处理连接
		go rpc.ServeConn(conn)
	}
}
```
下面是tcp的rpc的客户端实现

```
package main

import (
	"fmt"
	"net/rpc"
)

type EchoServer bool

func main() {
	//1、定义服务
	client, err := rpc.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	args := "this is tcp request"
	var reply string
	//调用服务方法
	err = client.Call("EchoServer.Echo", &args, &reply)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("tcp rpc reply=%s\n", reply)
}
```

###2、HTTP###

server端

```
package main

import (
	"fmt"
	"net/http"
	"net/rpc"
)

//1、定义服务
type Service struct {
}

//2、定义服务的方法
func (this *Service) Echo(req *string, res *string) error {
	fmt.Println("Echo...")
	*res = *req
	return nil
}

func main() {
	var err error
	//3、创建服务
	s := new(Service)
	//4、注册服务
	err = rpc.Register(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	//5、接管http请求
	rpc.HandleHTTP()
	http.ListenAndServe(":8080", nil)
}
```

client端

```
package main

import (
	"fmt"
	"net/rpc"
)

//1、定义服务
type Service struct {
}

func main() {
	//1、连接到服务器，走http协议
	client, err := rpc.DialHTTP("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	args := "This is http request"
	var reply string
	//2、调用远程服务
	err = client.Call("Service.Echo", &args, &reply)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("http rpc reply=%s\n", reply)
}
```

###3、json###

server端

```
  package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

//1、定义服务
type Service struct {
}

//2、定义服务的方法
func (this *Service) Echo(req *string, res *string) error {
	fmt.Println("Echo...")
	*res = *req
	return nil
}

func main() {
	//3、创建服务
	s := new(Service)
	//4、注册服务
	rpc.Register(s)
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		//5、将连接交给jsonrpc处理
		go jsonrpc.ServeConn(conn)
	}
}
  
```

client端

```
package main

import (
	"fmt"
	"net/rpc/jsonrpc"
)

type Service struct {
}

func main() {
	//连接服务器
	client, err := jsonrpc.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	args := "This is json request"
	var reply string
	//调用远程服务
	err = client.Call("Service.Echo", &args, &reply)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("json rpc reply =%s\n", reply)
}

```

从以上三种实现方式来看，在服务器一端，在golang的内置rpc中编程，其实步骤是很清晰的，它们是

* 1、定义服务
* 2、定义服务的方法
* 3、new服务
* 4、注册服务
* 5、将rpc和链接建立关系


其中上述的三种方式的rpc，前4个步骤是完全一样的，不同的是第5个步骤将rpc和链接建立关系

* 在tcp中是通过rpc.ServeConn(conn)

* 在http中是通过rpc.HandleHTTP()

* 在jsonrpc中是通过jsonrpc.ServeConn(conn)

其中我们可以看出，在代码的实现上，tcp和jsonrpc只有一行代码是不一样的

在客户端一端，代码几乎就是一样的

* 建立连接（使用不同的包的Dial,然后得到不同的client）

* 调用远程服务


