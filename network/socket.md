###字节序###

<font color="red">下文中提到的函数，想看原型，请打开/usr/include目录来查看，看源码才更清晰。(Linux下)</font>

```
#include <netinet/in.h>
#include <bits/socket.h>
#include <sys/un.h>
#include <arpa/inet.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <syslog.h> //日志系统
```


* 大端字节序（网络字节序）
* 小端字节序 (主机字节序)


###大端字节序（网络字节序）###
一个整数的高位字节（23~32bit）存储在内存的低地址，低位字节(0~7bit)存储在内存的高地址处。（网大，也就是我们假设，从左往右，地址编址是从小到大的，数字从左往右我们知道高为到低位的，如果是这样的字节序，我们就称为大端字节序。---整数的高位存在地址的低地址）

###小端字节序（主机字节序）###
整数的高位字节存储在内存的高地址，而低位字节则存储在内存的低地址处。（现代PC大多数采用小端字节序，因此小端字节序又被称为主机字节序。）

Linux提供了如下4个函数来完成主机字节序和网络字节序之间的转换：

```
 #include <netinet/in.h>
 unsigned long int htonl(unsigned long int hostlong);
 unsigned short int htons(unsigned short int hostshort);
 unsigned long int ntohl(unsigned long int netlong);
 unsigned short int ntohs(unsigned short int netshort);
```

长整型函数通常用来转换IP地址，短整型函数用来转换端口号（当然不限于此。任何格式化的数据通过网络传输时，都应该使用这些函数来转换字节序）

###通用socket地址###

结构体sockaddr表示的是socket的地址,定义在bits/socket.h中,也就是需要使用这个结构体时，必须

```
 #include <bits/socket.h>
```

它的源码里结构体sockaddr是这样子定义的。

```
struct sockaddr
{
    sa_family_t sa_family;
    
    char sa_data[14];
}
```

<font color="red">地址族类型通常与协议族类型对应</font>

协议族

* PF_UNIX
* PF_INET
* PF_INET6

地址族

* AF_UNIX
* AF_INET
* AF_INET6


宏PF_ * 和 AF_ *都定义在bits/socket.h头文件中，且后者与前者有完全相同的值，所以二者通常混用。

sa_data成员用于存放socket地址值，但是，不同的协议族的地址值具有不同的含义和长度。

* PF_UNIX   108字节
* PF_INET   6字节 
* PF_INET6  26字节

从上面可以看出，14字节的sa_data根本无法完全容纳多数协议族的地址，因此,linux定义了下面这个新的通用的socket地址结构体：

```
struct sockaddr_storage
{
   	sa_family_t sa_family;
   
   	unsigned long int __ss_align;
   
   	char __ss_padding[128-sizeof(__ss_align)];
}
```

上面这个两个socket地址结构体(sockaddr和sockaddr_storage)显然很不好用，所以Linux为各个协议族提供了专门的socket地址结构体。

```
  #include <sys/un.h>
```


UNIX本地域协议族使用如下专用socket地址结构体：

```
   struct sockaddr_un
   {
        sa_family_t sin_family;/*地址族：AF_UNIX*/
        char sun_path[108];/*文件路径名*/
   }
```

TCP/IP协议族有sockaddr_in和sockadd_in6两个专用socket地址结构体，它们分别用于IPv4和IPv6:

```
struct sockaddr_in
{
    sa_family_t sin_family; /*地址族：AF_INET*/
    u_int16_t sin_port;   /*端口号，要用网络字节序表示*/
    struct in_addr sin_addr;/*IPv4地址结构体*/
}

struct in_addr
{
   u_int32_t s_addr; /*IPv4地址，要用网络字节序表示*/
}


struct sockaddr_in6
{
   sa_family_t sin6_family;/*地址族:AF_INET6*/
   u_int16_t sin6_port;/*端口号，要用网络字节序表示*/
   u_int32_t sin6_flowinfo;/*流信息，应设置为0*/
   struct in6_addr sin6_addr;/*IPv6结构体*/
   u_int32_t sin6_scope_id;  /*scope ID,尚处于实验阶段*/
}

struct in6_addr
{
   unsigned char sa_addr[16];/*IPv6地址，要用网络字节序表示*/
}
   
```

<font color="red">所有专用地址类型的变量在实际使用时都需要转换为通用socket地址类型sockaddr(强制转换即可)，因为所有socket编程接口使用的地址参数的类型都是sockaddr。</font>

通常，人们习惯用可读性好的字符串来表示IP地址：

* 用点分十进制字符串表示IPv4地址
* 用十六进制字符串表示IPv6地址

<font color="red">但编程中我们需要先把它们转换为整数（二进制）方能使用。</font>
需要

```
#include <arpa/inet.h>
```

函数原型是:

```
in_addr_t  inet_addr(const char* strptr);

int inet_aton(const char* cp,struct in_addr inp);

char* inet_ntoa(struct in_addr in);

```

这是用于IPv4地址的转换函数

* inet_addr函数将用点分十进制字符串表示的IPv4地址转换为用网络字节序整数表示的IPv4地址，它失败时返回INADDR_NONE.
* inet_aton函数完成和inet_addr同样的功能，但是将转化结果存在于参数inp指向的地址结构中，它成功是返回1，失败则返回0.
* inet_ntoa函数将用网络字节序整数表示的IPv4地址转化为用点分十进制字符串表示的IPv4地址。<font color="red">但需要注意的是，该函数内部用一个静态变量存储转换结果，函数的返回值指向该静态内存，因此，inet_ntoa是不可重入的。</font>如果你想要保存转换的结果，只能用拷贝的方式，即用另外一块内存把内容先拷贝出来。

###通用的IP地址转换函数(<font color="red">重点记忆</font>)###
同时使用于IPv4地址和IPv6地址

```
/*af表示地址族,可以是AF_INET或AF_INET6*/

int inet_pton(int af,const char* src,void* dst);

/*参数cnt指定目标存储单元的大小*/
const char* inet_ntop(int af,const void* src,char dst,socklen_t cnt);

#define INET_ADDRSTRLEN  16（暂时不明白为什么）
#define INET6_ADDRSTRLEN  46（暂时不明白为什么）
/*这两个宏是在<netinet/in.h>中定义的*/

```

###1、创建socket###

```
int socket(int domain,int type,int protocol);
/*domain ,可选值为PF_INET,PF_INET6等等*/
/*type可选值SOCK_STREAM和SOCK_UGRAM等等*/
/*几乎在所有情况下，我们都应该它设置为0，表示使用默认协议*/
```

<font color="red">socket系统调用成功时返回一个socket文件描述符，失败则返回-1并设置errno</font>

###2、命名socket###

<font color="red">将一个socket与socket地址绑定称为给socket命名。</font>在服务器程序中，我们通常要命名socket,因为只有命名后客户端才能知道该如何连接它。客户端则通常不需要命名socket,而是采用匿名方式，即使用操作系统自动分配的socket地址。命名socket的系统调用是bind

```
int bind(int sockfd,const struct sockaddr* my_addr,socklen_t addrlen);

```


<font color="red">bind系统调用成功是返回0，失败则返回-1,并设置errno</font>.

其中两种常见的errno是EACCES和EADDRINUSE,它们的含义分别是:

* <font color="red"> EACCES,被绑定的地址是受保护的地址，仅超级用户能够访问比如普通用户将socket绑定到知名服务端口(端口号为0~1023)上时，bind将返回EACCES错误。</font>
* <font color="red">EADDRINUSE,被绑定的地址正在使用中，比如将socket绑定到一个处于TIME_WAIT状态的socket地址（又抑或是同一个程序起两个实例的时候，也会有这个问题--使用同一个端口时）。</font>

###3、监听socket###

```
int listen(int sockfd,int backlog);

```
sockfd参数指定被监听的socket,backlog参数提示内核监听队列的最大长度，监听队列的长度如果超过backlog，服务器将不在受理新的客户连接，客户端也将受到ECONNREUSED错误信息。backlog参数的典型值是5。

<font color="red">listen系统调用成功时返回0.失败则返回-1并设置errno。</font>


###4、接收连接###
从listen监听队列中接受一个连接:

```
int accept(int sockfd,struct sockaddr * addr,socklen_t *addrlen);

```

sockfd参数是执行过listen系统调用的监听socket,addr参数用来获取被接受连接的远端socket地址，该socket地址的长度由addrlen参数指出。


<font color="red">accept成功时返回一个新的连接socket,该socket唯一地标识了被接受的这个连接，服务器可通过读写该socket来与被接受连接对应的客户端通信。accept失败时返回-1并设置errno。
accept只是从监听队列中取出连接，而不论连接处于何种状态，更不关心任何网络状况的变化。</font>


###5、发起连接###

```
int connect(int sockfd,struct sockaddr *serv_addr,socklen_t addrlen);

```

<font color="red">connect成功时返回0，一旦成功建立连接，sockfd就唯一标识了这个连接，客户端就可以通过读写sockfd来与服务器通信。connect失败则返回-1并设置errno.</font>

常见的errno是ECONNREFUSED和ETIMEDOUT:

* ECONNREFUSED，目标端口不存在，连接被拒绝。
* ETIMEDOUT，连接超时。

###6、关闭连接###

关闭一个连接实际上就是关闭该连接对应的socket,这可以通过如下关闭普通文件描述符的系统调用来完成：

```
int close(int fd);

```
fd参数是待关闭的socket,不过，close系统调用并非总是立即关闭一个连接，而是将fd的引用计数减1，只有当fd的引用计数为0时 ,才真正关闭连接。多进程程序中，一次fork系统调用默认将使父进程中打开的socket的引用计数加1，因此我们必须在父进程和子进程中都对该socket执行close调用才能将连接关闭。（close是同时关闭读和写的）

<font color="red">如果无论如何都要立即终止连接（而不是讲socket的引用计数减1），可以使用如下的shutdown系统调用（相对于close来说，它是专门为网络编程设计的）</font>

```
int shutdown(int sockfd,int howto);

```
howto的可选值为：

* SHUT_RD 关闭sockfd上读的这一半，应用程序不能再针对socket文件描述符执行读操作，并且该socket接收缓冲区中的数据都被丢失。
* SHUT_WR  关闭sockfd上写的这一半，sockfd的发送缓冲区中的数据会在真正关闭连接之前全部发送出去，应用程序不可再对socket文件描述符执行写操作。这种情况下，连接处于半关闭状态。 
* SHUT_RDWR  同时关闭sockfd上的读和写

<font color="red">shutdown成功时返回0，失败则返回-1并设置errno。</font>

###TCP数据读写###
对文件的读写操作read和write同样适用于socket,但是socket编程接口提供了几个专门用于socket数据读写的系统调用，它们增加了对数据读写的控制。

```
ssize_t recv(int sockfd,void* buf,size_t len,int flags);

ssize_t send(int sockfd,const void* buf,size_t len,int flags);

```

recv读取sockfd上的数据，buf和len参数分别指定读缓冲区（应用缓冲区--非tcp读缓冲区）的位置和大小，flags通常设置为0即可。recv成功时返回实际读取到的数据的长度，它可能小于我们期望的程度len,因此我们可能要多次调用recv，才能读取到完整的数据。recv可能返回0，这以为着通信对方已经关闭了连接了。recv出错时返回-1并设置errno.


send往sockfd上写入数据，buf和len参数分别指定写缓冲区（应用写缓冲区--非tcp写缓冲区）的位置和大小。send成功时返回实际写入的数据的长度，失败则返回-1并设置errno。

flags的可选值

* MSG_CONFIRM 指示数据链路层协议持续监听对方的回应，直到得到答复，它仅能用于SOCK_DGRAM和SOCK_RAW类型的socket。
* MSG_DONTROUTE 不查看路由表，直接将数据发送给本地局域网络内的主机。这表示发送者确切地直到目标主机就在本地网络上。
* MSG_DONTWAIT 对socket的此次操作将是非阻塞的。
* MSG_MORE 告诉内核应用程序还有更多数据要发送，内核将超市等待新数据写入tcp发送缓冲区后一并发送，这样可防止TCP发送过多小的报文段，从而提高传输效率。
* MSG_WAITALL 读操作仅在读取到指定数量的字节后才返回。
* MSG_OOB    发送或接收紧急数据。
* MSG_NOSIGNAL 往读端关闭的管道或者socket链接中写数据时不引发SIGPIPE信号。

















