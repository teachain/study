  [网络地址转换](http://baike.baidu.com/view/875777.htm)(NAT,Network Address Translation)属接入广域网(WAN)技术，是一种将私有（保留）地址转化为合法IP地址的转换技术，它被广泛应用于各种类型Internet接入方式和各种类型的网络中。原因很简单，NAT不仅完美地解决了lP地址不足的问题，而且还能够有效地避免来自网络外部的攻击，隐藏并保护网络内部的计算机。

NAT（Network Address Translation，网络地址转换）是将IP 数据包头中的IP 地址转换为另一个IP 地址的过程。在实际应用中，NAT 主要用于实现私有网络访问公共网络的功能。这种通过使用少量的公有IP 地址代表较多的私有IP 地址的方式，将有助于减缓可用IP地址空间的枯竭 

NAT工作原理

```
NAT的基本工作原理是:当私有网主机和公共网主机通信的IP包经过NAT网关时，将IP包中的源IP或目的IP在私有IP和NAT的公共IP之间进行转换。NAT网关有2个网络端口，其中公共网络端口的IP地址是统一分配的公共IP，为202.204.65.2；私有网络端口的IP地址是保留地址，为192.168.1.1。私有网中的主机 192.168.1.2向公共网中的主机166.111.80.200发送了1个IP包（Des=166.111.80.200,Src=192.168.1.2）。当IP包经过NAT网关时，NAT会将IP包的源IP转换为NAT的公共 IP并转发到公共网，此时IP包（Des=166.111.80.200，Src=202.204.65.2）中已经不含任何私有网IP的信息。由于IP 包的源IP已经被转换成NAT的公共IP，响应的IP包（Des=202.204.65.2,Src=166.111.80.200）将被发送到NAT。 这时，NAT会将IP包的目的IP转换成私有网中主机的IP，然后将IP包（Des=192.168.1.2，Src=166.111.80.200）转 发到私有网。对于通信双方而言，这种地址的转换过程是完全透明的。
```

NAPT技术

```
由于NAT实现是私有IP和NAT的公共IP之间的转换，那么，私有网中同时与公共网进行通信的主机数量就受到NAT的公共IP地址数量的限制。为了克服 这种限制，NAT被进一步扩展到在进行IP地址转换的同时进行Port的转换，这就是网络地址端口转换NAPT（Network Address Port Translation）技术。
    NAPT与NAT的区别在于，NAPT不仅转换IP包中的IP地址，还对IP包中TCP和UDP的Port进行转换。这使得多台私有网主机利用1个NAT公共IP就可以同时和公共网进行通信。（NAPT多了对TCP和UDP的端口号的转换）

    私有网主机192.168.1.2要访问公共网中的 Http服务器166.111.80.200。首先，要建立TCP连接，假设分配的TCP Port是1010，发送了1个IP包（Des=166.111.80.200:80,Src=192.168.1.2:1010）,当IP包经过NAT 网关时，NAT会将IP包的源IP转换为NAT的公共IP，同时将源Port转换为NAT动态分配的1个Port。然后，转发到公共网，此时IP包 （Des=166.111.80.200：80，Src=202.204.65.2:2010）已经不含任何私有网IP和Port的信息。由于IP包的源 IP和Port已经被转换成NAT的公共IP和Port，响应的IP包 （Des=202.204.65.2:,Src=2010166.111.80.200:80）将被发送到NAT。这时NAT会将IP包的目的IP转换成 私有网主机的IP，同时将目的Port转换为私有网主机的Port，然后将IP包 （Des=192.168.1.2:1010，Src=166.111.80.200:80）转发到私网。对于通信双方而言，这种IP地址和Port的转 换是完全透明的。
```

upnp和nat-pmp差不多，就是在路由器和内部机器提供一个中间服务，内部机器请求upnp将其使用到的端口跟某个外网端口绑定，这样当路由器收到外网请求时先去upnp里查找是否此外网端口已经被upnp映射，如果被映射则将数据转发到内部机器对应的端口。

**napt是路由器肯定带的功能**，其产生的nat映射表有多种类型，但都**有时效**，也就是超过一段时间原来的nat映射就无效，然后新建新的nat映射。**nat映射必须先由内部机器向外部网络发起请求才会产生。**

upnp是把映射关系长期保存下来，**外部机器可以主动向内部机器请求网络连接。** **所以首先要路由器开启upnp功能（一般由用户去路由器设置里手动开启upnp）**，然后内部机器的程序要自己实现upnp客户端功能：主动查找upnp服务，主动增加映射、删除映射等。

**客户端无法控制natp的映射，可以主动控制upnp映射。**



https://blog.csdn.net/chenlycly/article/details/52357088



DNAT Destination Network Address Translation 目的网络地址转换换，

 SNAT Source Network Address Translation 源网络地址转换，

其作用是将ip数据包的源地址转换成另外一个地址，可能有人觉得奇怪，好好的为什么要进行ip地址转换啊，为了弄懂这个问题，我们要看一下局域网用户上公网的原理，假设内网主机A（192.168.2.8）要和外网主机B（61.132.62.131）通信，A向B发出IP数据包，如果没有SNAT对A主机进行源地址转换，A与B主机的通讯会不正常中断，因为当路由器将内网的数据包发到公网IP后，公网IP会给你的私网IP回数据包，这时，公网IP根本就无法知道你的私网IP应该如何走了。所以问它上一级路由器，当然这是肯定的，因为从公网上根本就无法看到私网IP，因此你无法给他通信。为了实现数据包的正确发送及返回，网关必须将A的址转换为一个合法的公网地址，同时为了以后B主机能将数据包发送给A，这个合法的公网地址必须是网关的外网地址，如果是其它公网地址的话，B会把数据包发送到其它网关，而不是A主机所在的网关，A将收不到B发过来的数据包，所以内网主机要上公网就必须要有合法的公网地址，而得到这个地址的方法就是让网关进行SNAT(源地址转换），将内网地址转换成公网址(一般是网关的外部地址），所以大家经常会看到为了让内网用户上公网，我们必须在routeros的firewall中设置snat，俗称IP地址欺骗或伪装（masquerade)  

区分SNAT和DNAT 从定义来讲它们一个是源地址转换，一个是目标地址转换。

都是地址转换的功能，将私有地址转换为公网地址。 要区分这两个功能可以简单的由连接发起者是谁来区分： 内部地址要访问公网上的服务时（如web访问），内部地址会主动发起连接，由路由器或者防火墙上的网关对内部地址做个地址转换，将内部地址的私有IP转换为公网的公有IP，网关的这个地址转换称为SNAT，主要用于内部共享IP访问外部。 当内部需要提供对外服务时（如对外发布web网站），外部地址发起主动连接，由路由器或者防火墙上的网关接收这个连接，然后将连接转换到内部，此过程是由带有公网IP的网关替代内部服务来接收外部的连接，然后在内部做地址转换，此转换称为DNAT，主要用于内部服务对外发布。 在配置防火墙或者路由acl策略时要注意这两个NAT一定不能混淆。