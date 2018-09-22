##TCP##

搞过网络的人，一般都会用到抓包分析工具，在windows下一般就是wireshark，在linux下面一般系统自带tcpdump.
tcpdump是强大的socket抓包工具，可分析网络通信情况，mac下使用必须用root权限执行该工具。
sudo tcpdump， 输入root密码即可使用。

面向连接、字节流和可靠传输。

使用TCP协议通信的双方必须先建立连接，然后才能开始数据的读写。双方都必须为该连接分配必要的内核资源，以管理连接的状态和连接上数据的传输。TCP连接是全双工的，即双方的数据读写可以通过一个连接进行。完成数据交换之后，通信双方都必须断开连接以释放系统资源。

<font color="red">TCP协议的这种连接是一对一的，所以基于广播和多播（目标是多个主机地址）的应用程序不能使用TCP服务。而无连接协议UDP则非常适合于广播和多播。</font>

###发送###
当发送端应用程序连接执行多次写操作时，TCP模块先将这些数据放入TCP发送缓冲区中。当TCP模板真正开始发送数据时，发送缓冲区中这些等待发送的数据可能被封装成一个或多个TCP报文段发出。因此，TCP模块发送出的TCP报文段的个数和应用程序执行的写操作次数之间并有固定的数量关系。

###接收###
当接收端收到一个或多个TCP报文段后，TCP模块将它们携带的应用程序数据按照TCP报文段的序号依次放入TCP接收缓冲区中，并通知应用程序读取数据。应用程序执行的读操作次数和TCP模块接收到的TCP报文段个数之间也没有固定的数量关系。

也就是说应用程序对数据的发送和接收是没有边界限制的。

TCP传输是可靠的。首先，TCP协议采用发送应答机制，即发送端发送的每个TCP报文段都必须得到接收方的应答，才认为这个TCP报文段传输成功。其次，TCP协议采用超时重传机制，发送端在发送出一个TCP报文段之后启动定时器，如果在定时时间内未收到应答，它将重发改报文段。最后，因为TCP报文段是以IP数据包发送的，而IP数据包报到达接收端可能乱序，重复，所以TCP协议还会对接收到的TCP报文重排，整理，再交付给应用层。

TCP头部中的16位窗口大小(window size),是TCP流量控制的一个手段。这里说的窗口，指的是接收通告窗口（Receiver Window,RWND）.它告诉对方本端的TCP接收缓冲区还能容纳多少字节的数据，这样对方就可以控制发送数据的速度。

TCP头部最长是60字节，其中包括20字节的固定部分，还是最多40字节的选项字段。

TCP的一端都有自己的序号seq,A端发送到B端，那么B端回复给A端的ack就是A端的seq+1,A端发送给B端的ack就是B端的seq+1,伪代码

```
A--->B    

B---->A   ack=seqA+1

A---->B   ack=seqB+1 

```
从上面可以看出确认号（ack）都是对方的seq+1

TCP连接是全双工的，所以它允许两个方向的数据传输被独立关闭，换言之，通信的一端可以发送结束报文段给对方，告诉它本端已经完成了数据的发送，但允许继续接收来自对方的数据，直到对方也发送结束报文段以关闭连接。TCP连接的这种状态称为半关闭状态(half close).

服务器和客户端应用程序判断对方是否已经关闭连接的方法是:read系统调用返回0(收到结束报文段)。

<font color="red">socket网络编程接口通过shutdown函数提供了对半关闭的支持。</font>






