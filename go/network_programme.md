TCP/Ip协议族是一个四层协议系统，自底而上分别是数据链路层，网络层，传输层和应用层。每一层完成不同的功能，且通过若干协议来实现，上层协议使用下层协议提供的服务。

数据链路层两个常用的协议是ARP协议（Address Resolve Protocol,地址解析协议）和RARP(Reverse Address Resolve Protocol,逆地址解析协议)。它们实现了IP地址和机器物理地址（mac地址）之间的相互转换。

网络层使用IP地址寻址一台机器，而数据链路层使用物理地址寻址一台机器，因为网络层必须先将目标机器的IP地址转化成物理地址，才能使用数据链路层提供的服务，这就是ARP协议的用途。

###网络层###

网络层实现数据包的选路和转发。WAN(广域网)通常使用众多分级的路由器来连接分散的主机或LAN（局域网）,因此，通信的两台知己一般不是直接相连的，而是通过多个中间节点（路由器）连接的，网络层的任务就是选择这些中间节点，以确定两台主机自检的通信路径。同时，网络层对上层协议隐藏了网络拓扑连接的细节，使得在传输层和网络应用程序看来，通信的双方是直接相连的。

网络层最核心的协议是IP协议（Internet Protocol,因特网协议），IP协议根据数据包的目的IP地址来决定如何投递它。如果数据包不能直接发送给目标主机，那么IP协议就为它寻找一个合适的下一跳路由器，并将数据包交付给该路由器来转发。多次重复着一个过程，数据包最终到达目标主机，或者由于发送失败而被丢弃。可见IP协议使用逐跳（hop by hop）的方式确定通信路径。

网络层另外一个重要的协议是ICMP协议（Internet control Message Protocol,因特网控制报文协议），它是IP协议的重要补充，主要用于检测网络连接。

ICMP一共使用32位来表示一个报文，8位用来表示类型，8位用来表示代码，16位用做校验和。ICMP报文分为两大类：

* 差错报文
* 查询报文。

差错报文主要用来回应网络错误，比如目标不可到达（类型值为3）和重定向（类型值为5）。

查询报文用来查询网络信息，比如我们平常用的ping命令就是使用ICMP报文来查看目标是否可到达的（类型值为8）。

<font color="red">注意：ICMP协议并非严格意义上的网络层协议，因为它使用处于同一层的IP协议提供的服务。（一般来说，上层协议使用下层协议提供的服务）。</font>



###传输层###

传输层为两台主机上的应用程序提供端到端（end to end）的通信，与网络层使用的逐跳通信方式不同，传输层只关心通信的起始端和目的端，而不在乎数据包的中转过程。

传输层协议主要有三个:

* TCP协议
* UDP协议
* SCTP协议

TCP协议(Transmission Control Protocol,传输控制协议)为应用层提供可靠的、面向连接的和基于流的服务。TCP协议使用超时重传、数据确认等方式来确保数据包被正确地发送至目的端，因此TCP服务是可靠的。使用TCP协议通信的双方必须先建立TCP连接，并在内核中为该连接维持一些必要的数据结构，比如连接的状态，读写缓冲区，以及诸多定时器等，当通信结束时，双方必须关闭连接以释放这些内核数据。TCP服务是基于流的，基于流的数据没有边界（长度）限制，它源源不断地从通信的一端流入另一端，发送端可以逐个字节地向数据流中写入数据，接收端也可以逐个字节地将它们读出。

UDP协议（User Datagram Protocol,用户数据包协议）为应用层提供不可靠，无连接和基于数据包的服务。应用程序要自己处理数据确认，超时重传等逻辑，每个UDP数据包都有一个长度，接收端必须以该长度为最小单位将其所有内容一次性读出，否则数据将截断。


SCTP协议（Stream Control Transimission Protocol,流控制传输协议），他是为了在因特网上传输电话信号而设计的。

###应用层###

* telnet协议是一种远程登录协议
* OSPF协议（开放最短路径优先）
* DNS（Domain Name Service,域名服务）协议，提供域名到ip地址的转换

<font color="red">当发送端应用程序使用send(或者write)函数向一个TCP连接写入数据时，内核中的TCP模块首先将这些数据复制到与该连接对应的TCP内核发送缓冲区，然后TCP模块调用IP模块提供的服务。</font>

经过数据链路层封装的数据成为帧（frame）,帧的最大传输单元（Max Transmit Unit,MTU），即帧最多能携带多少上层协议数据，通常受到网络类型的限制，<font color="red">以太网的MTU是1500字节.</font>

<font color="red">帧才是最终在物理网络上传输的字节序列</font>

<font color="red">应用层向数据链层走叫"封装"，从数据链路层向应用层走叫"分用"</font>

<font color="red">可以通过/etc/services文件来查看常用的协议和它们使用的端口号</font>

###ARP协议工作原理###

ARP协议能实现任意网络层地址到任意物理地址的转换。

其工作原理是：主机向自己所在的网络广播一个ARP请求，该请求包含目标机器的网络地址。此网络上的其他机器都将收到这个请求，但只有被请求的目标机器会回应一个ARP应答，其中包含自己的物理地址。

<font color="red">/etc/resolv.conf存放DNS服务器的IP地址，这个文件是由网络管理程序写入的。</font>

host -t A www.baidu.com 可以得到域名的别名和ip地址，host是一个命令。


我们来看一个具体的例子，上层协议需要下层协议提供的服务。以太网帧的MTU是1500字节,因此它携带的IP数据包的数据部分最多是1480字节（IP头部要占用20字节）。我们现在考虑用ip数据包封装一个长度为1481字节的ICMP报文，则这个ip数据包就需要分片（已经大于1480），那么这个ip数据包的分片就需要在网络层来处理，也就是分片好了以后再封装成以太网帧，然后调用数据链路层的服务来发送IP数据包。

