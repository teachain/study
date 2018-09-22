# DNS

DNS（Domain Name System，域名系统），是一种用于将域名解析为IP的服务器系统，当你上网时输入一个网址，它之所以能够找到该网址指向的服务器地址，都是靠域名系统来进行解析的。

CNAME是别名，意思是这个域名还有另外一个名字，两者指向同一个IP。 A指的是Address，即IP地址。 NS指的是服务器主机名

## DNS查询工具dig

```
dig，其实是一个缩写，即Domain Information Groper
当直接使用dig命令，不加任何参数和选项时，dig会向默认的上连DNS服务器查询“.”（根域）的NS记录。
```

### dig的基本的命令格式是：

```
dig @dnsserver name querytype   
#比如使用google的8.8.8.8的dns服务器查百度的a记录,则命令为
# dig @8.8.8.8 www.baidu.com A
```

如果你设置的dnsserver是一个域名，那么dig会首先通过默认的上连DNS服务器去查询对应的IP地址，然后再以设置的dnsserver为上连DNS服务器。如果你没有设置@dnsserver，那么dig就会依次使用/etc/resolv.conf里的地址作为上连DNS服务器。而对于querytype，你可以设置A/AAAA/PTR/MX/ANY等值，默认是查询A记录。



通过dig命令也可以查看DNS服务器的主从关系。

```
dig -t soa www.baidu.com
SOA是start of authority的简称，提供了DNS主服务器的相关信息，在soa之后我们可以看到7个参数，依次是： 
1、DNS主服务器名 
2、管理员的E-mail，这里是baidu.dns.master@baidu.com，由于@在数据库文件里有特殊作用，所以这里是用.代替的。 
3、更新序号。表示数据库文件的新旧，一般是用时间来表示，这里1703230011表示的是2017年3月23日进行了一次更新，当天更新编号0011. 
4、更新频率。 表示每5秒，slave服务器就要向master服务器索取更新信息。 
5、失败重试时间，当某些原因导致Slave服务器无法向master服务器索取信息时，会隔5秒就重试一次。 
6、失效时间。如果一直重试失败，当重试时间累积达到86400秒时，不再向主服务器索取信息。 
7、缓存时间。默认的TTL缓存时间。
```

