#### **DNS是什么**

https://blog.csdn.net/ceshi986745/article/details/51787406?utm_source=blogxgwz4

DNS（Domain Name System——域名系统）的作用非常简单，就是根据提供的域名，来查询相应的IP地址。这个系统的目的也很简单，避免人们去记忆数字形式的IP地址，因为字母比数字要容易记些。就好像你的电话簿，你不可能记住里面所有的电话号码，但你可以很容易的知道一个人的名字。

#### **查询过程**

在linux下，你可以通过命令dig来显示DNS的查询过程。

```
dig baidu.com
```

#### 查询结果

```
; <<>> DiG 9.10.6 <<>> baidu.com
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 64975
;; flags: qr rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 5, ADDITIONAL: 6

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;baidu.com.			IN	A

;; ANSWER SECTION:
baidu.com.		497	IN	A	123.125.115.110
baidu.com.		497	IN	A	220.181.57.216

;; AUTHORITY SECTION:
baidu.com.		86297	IN	NS	dns.baidu.com.
baidu.com.		86297	IN	NS	ns2.baidu.com.
baidu.com.		86297	IN	NS	ns3.baidu.com.
baidu.com.		86297	IN	NS	ns4.baidu.com.
baidu.com.		86297	IN	NS	ns7.baidu.com.

;; ADDITIONAL SECTION:
dns.baidu.com.		86334	IN	A	202.108.22.220
ns2.baidu.com.		86302	IN	A	61.135.165.235
ns3.baidu.com.		86302	IN	A	220.181.37.10
ns4.baidu.com.		86297	IN	A	220.181.38.10
ns7.baidu.com.		172799	IN	A	119.75.219.82

;; Query time: 8 msec
;; SERVER: 202.101.172.35#53(202.101.172.35)
;; WHEN: Fri Oct 26 09:19:55 CST 2018
;; MSG SIZE  rcvd: 240

```

查询的结果包括了五段内容：

1. 第一段为查询的参数和统计

2. 查询内容

   ```
   ;; OPT PSEUDOSECTION:
   ; EDNS: version: 0, flags:; udp: 4096
   ;; QUESTION SECTION:
   ;baidu.com.			IN	A
   ```

   查询域名baidu.com的A记录，其中A表示地址，即address的缩写



#### DNS服务器

想要使用DNS服务，你当然要知道DNS服务器的IP地址，否则你没法上网浏览网页了。DNS服务器的IP地址可以由DHCP协议为子网中的计算机动态分配，也可以使用固定的地址（比如谷歌的8.8.8.8）。Linux中，DNS服务器的IP地址保存在/etc/resolv.conf中。dig命令默认使用本机设置好的DNS服务器地址进行查询，你也可以通过参数”@”，来显示指定其他DNS服务器进行查询。

```
 dig @8.8.8.8 baidu.com
```

#### 域名的分级

域名分为多层，DNS服务器查询域名的IP地址，采用的方式也是分级查询。如果你仔细看前面的查询结果，你会发现每个域名的后面都有一个点”.”。
我们在查询时输入的是baidu.com，而显示的结果却多了一个点。这不是错误，而是所有域名的尾部都有一个根域名。比如，www.example.com真正的域名是www.example.com.root，简写为www.example.com.。因为，根域名.root对于所有域名都是一样的，所以加以省略。 

#### 域名的层级结构如下

```
主机名.次级域名.顶级域名.根域名
```

#### 根域名服务器

DNS根据域名的层级，进行分级查询。 
你需要知道，每一级域名都有自己的NS记录，NS记录指向该级域名的域名服务器。这些服务器知道下一级域名的各种记录。分级查询的就是一级一级的查询域名的NS记录，知道查询到最终的IP地址，大致过程如下：

```
从“根域名服务器”查到“顶级域名服务器”的NS记录和A记录
从“顶级域名服务器”查到“次级域名服务器”的NS记录和A记录
从“次级域名服务器”查到“主机名”的IP地址
```


如过仔细看，你会发现没有提到根域名服务器的IP地址是怎么知道的，答案是，根域名服务器的NS记录和IP地址几乎不会变，所以在DNS服务器中设置好了。根域名服务器全世界一共有十三组，从“A.ROOT-SERVERS.NET——M.ROOT-SERVERS.NET”，你可以访问根域名服务器，来获取根域名服务器的NS和IP地址信息。

#### 分级查询演示

通过dig命令的+trace参数可以显示DNS的整个查询过程

```
dig +trace baidu.com
```

查询的结果

```
; <<>> DiG 9.10.6 <<>> +trace baidu.com
;; global options: +cmd
.			85970	IN	NS	j.root-servers.net.
.			85970	IN	NS	e.root-servers.net.
.			85970	IN	NS	m.root-servers.net.
.			85970	IN	NS	f.root-servers.net.
.			85970	IN	NS	h.root-servers.net.
.			85970	IN	NS	c.root-servers.net.
.			85970	IN	NS	a.root-servers.net.
.			85970	IN	NS	g.root-servers.net.
.			85970	IN	NS	i.root-servers.net.
.			85970	IN	NS	d.root-servers.net.
.			85970	IN	NS	b.root-servers.net.
.			85970	IN	NS	l.root-servers.net.
.			85970	IN	NS	k.root-servers.net.
.			85970	IN	RRSIG	NS 8 0 518400 20181107170000 20181025160000 2134 . neNmFDknOZ053yxo+QPWTA61cNlb+H0QV97KjPuZZWeYiQYXDzpist+n B/5pE9UL9HlLsUGp2w7RyK8FKL/0cJj0MjUsJ3jIGNdq7/u0WcUMwwFp YNm8L++ptIxza27Wqimcd7fCrzhi2Mcjr5Nixgkd38nAY1MKGETy55Hk s85Td027JDibCImhAPz3bcJrzQuCDPBOFF16gdcIEjJ59D2/L2qTM3Jo 8PVcA/LngVxOzAcDvopiMhyJIwjfiPyRE/l4VH6y5LzS+OTjFca9s+6k Ecxm77ojJUGryYGMC8tZ66yeEV8YOWkbKKgrWgYgtTh330Oorrem/sc7 p4ye2Q==
;; Received 1097 bytes from 202.101.172.35#53(202.101.172.35) in 4 ms

com.			172800	IN	NS	m.gtld-servers.net.
com.			172800	IN	NS	k.gtld-servers.net.
com.			172800	IN	NS	h.gtld-servers.net.
com.			172800	IN	NS	d.gtld-servers.net.
com.			172800	IN	NS	b.gtld-servers.net.
com.			172800	IN	NS	f.gtld-servers.net.
com.			172800	IN	NS	e.gtld-servers.net.
com.			172800	IN	NS	c.gtld-servers.net.
com.			172800	IN	NS	j.gtld-servers.net.
com.			172800	IN	NS	g.gtld-servers.net.
com.			172800	IN	NS	a.gtld-servers.net.
com.			172800	IN	NS	i.gtld-servers.net.
com.			172800	IN	NS	l.gtld-servers.net.
com.			86400	IN	DS	30909 8 2 E2D3C916F6DEEAC73294E8268FB5885044A833FC5459588F4A9184CF C41A5766
com.			86400	IN	RRSIG	DS 8 1 86400 20181107170000 20181025160000 2134 . hcsllBQMbanONW+yUDVEcmGsaqGil0Tw02WaaO1nLnLW+72JPqytWT71 dpeNmlXgVzdNZgC/nS7Y87iE8Yp8dRyn60Ng+l/PoM8OyqZU/r/dosN+ IAi0Sp9PQVb1+Z1R4H6TBVqWsh9RakO0pq8RNAFgHnWQuPxmsp6kHBKx bVkPZFPjwTjjj4t3BkkYow7kfRGot3DhlBRUXsdSG8EL4bJsTmGpAP1U AOvm/qmZygcgb81GBftwx4HXlKjM4qyR3d19erirnqGTHffpaYFYDUfi aqQ9JkWzSldr0xSz77OLQ4anyKFW4MV0rpynoceyGzuzoflA/cX9fgFq InyQVA==
;; Received 1169 bytes from 199.9.14.201#53(b.root-servers.net) in 166 ms

baidu.com.		172800	IN	NS	dns.baidu.com.
baidu.com.		172800	IN	NS	ns2.baidu.com.
baidu.com.		172800	IN	NS	ns3.baidu.com.
baidu.com.		172800	IN	NS	ns4.baidu.com.
baidu.com.		172800	IN	NS	ns7.baidu.com.
CK0POJMG874LJREF7EFN8430QVIT8BSM.com. 86400 IN NSEC3 1 1 0 - CK0Q1GIN43N1ARRC9OSM6QPQR81H5M9A  NS SOA RRSIG DNSKEY NSEC3PARAM
CK0POJMG874LJREF7EFN8430QVIT8BSM.com. 86400 IN RRSIG NSEC3 8 2 86400 20181031044301 20181024033301 37490 com. AM3Oob3vmp6ZcMUVhIZ6F4HsKXxdHIiwrLe2NgMwCEaXTKFzWefRtouC TCJ1CNYbvWLvqoHLQFSchVNgekPX6dqDcZ2mob58ki4WSVH3ljjbFQAQ PE77lVHZUsT9FvjxRZ+okCsCbkin3jBY3zkm0zvGq36+83CEuNqbVj0x hSY=
HPVV2B5N85O7HJJRB7690IB5UVF9O9UA.com. 86400 IN NSEC3 1 1 0 - HPVVP23QUO0FP9R0A04URSICJPESKO9J  NS DS RRSIG
HPVV2B5N85O7HJJRB7690IB5UVF9O9UA.com. 86400 IN RRSIG NSEC3 8 2 86400 20181029051452 20181022040452 37490 com. 2rhRUOV2dI6V2bbf1kTzPYm1NtdPUck6iHLqmcnPVM1+ADfQ0i2MBGFZ gyCVh9PiGzY8Ri7SblkNWdSfs+6TFdJyFKjBQJi0BE1PBiQUz2rIlwPY zGWiNyo3hHE0Ozm6dpQWS0jqEjkWgFOv5uPMbq1qVxyrlo9kcF09E6KU q/A=
;; Received 693 bytes from 192.35.51.30#53(f.gtld-servers.net) in 195 ms

baidu.com.		600	IN	A	123.125.115.110
baidu.com.		600	IN	A	220.181.57.216
baidu.com.		86400	IN	NS	dns.baidu.com.
baidu.com.		86400	IN	NS	ns4.baidu.com.
baidu.com.		86400	IN	NS	ns2.baidu.com.
baidu.com.		86400	IN	NS	ns7.baidu.com.
baidu.com.		86400	IN	NS	ns3.baidu.com.
;; Received 240 bytes from 220.181.37.10#53(ns3.baidu.com) in 33 ms
```



#### DNS的记录类型

域名和IP之间的对应关系，称为“记录”（record）。根据使用的目的不同，又分为不同的类型，常见的DNS记录类型如下：

A：地址记录（Address），返回域名指向的IP地址。

NS：域名服务器记录（Name Server），返回保存下一级域名信息的服务器地址。该记录只能设置为域名，不能设置为IP地址。

MX：邮件记录（Mail eXchange），返回接收电子邮件的服务器地址。

CNAME：规范名称记录（Canonical Name），返回另一个域名，即当前查询的域名是另一个域名的跳转。

PTR：逆向查询记录（Pointer Record），只用于从IP地址查询域名。