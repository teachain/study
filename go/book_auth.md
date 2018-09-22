##使用Json Web Token设计Passport系统##

规范链接 <a href="https://tools.ietf.org/html/draft-ietf-oauth-json-web-token-32">点我</a>

基于token的身份验证是无状态的，我们不将用户信息保存在服务器或session中，相比原始的Cookie+Session方式，更适合分布式系统的用户认证，绕开了传统的分布式session一致性等问题。

基于token的身份验证的主流程如下：

* 用户通过用户名和密码发送请求。
* 后端程序验证。
* 后端程序返回一个签名的token给客户端。
* <font color="red">客户端存储token</font>,并且每次发送请求都带上该token。

###JSON Web Token标准的设计###

JWT 标准的 Token 有三个部分：

```
header.payload.signature

```

三个部分中间用点分隔开，<font color="red">并且都使用 Base64 编码</font>，所以生成的 Token 类似这样

```
ewogICJ0eXAiOiAiSldUIiwKICAiYWxnIjogIkhTMjU2Igp9.ewogImlzcyI6ICJjaGJsb2dzLmNvbSIsCiAiZXhwIjogIjE0NzA3MzAxODIiLAogInVpZCI6ICIxMjM0NWFiY2RlIiwKfQ.9q2eq8sa374ao2uq9607r6qu6

```

####1、Header报头####

header 部分主要包括两部分，一个是 Token 的类型，另一个是使用的算法

```
{
"typ": "JWT",
"alg": "HS256"
}

```

####2、Payload载荷部分####

Payload 里面是 Token 的具体内容，这部分内容可以自定义，JWT有标准字段，也可以添加其它需要的内容。

```
标准字段：
iss：Issuer，发行者
sub：Subject，主题
aud：Audience，观众
exp：Expiration time，过期时间
nbf：Not before
iat：Issued at，发行时间
jti：JWT ID

```

这是一个典型的payload信息，包含了发行者（网站）、过期时间和用户id：

```
{
"iss": "chblogs.com",
"exp": "1470730182",
"uid": "12345abcde",
}
```

####3、Signature签名部分####

签名部分主要和token的安全性有关，Signature的生成依赖前面两部分。
首先将Base64编码后的Header和Payload用.连接在一起，对这个字符串使用HmacSHA256算法进行加密，这个密钥secret存储在服务端，前端不可见，然后将Signature和前面两部分拼接起来，得到最后的token：


从上面的组成中，我们可以看到它为什么叫json web token了，header和payload部分都是json组成的。


<font color="red">常规的token保存在sessionStorage或者localStorage中，每次请求时将token加在http请求的Header中</font>

####如何保证token的安全性####

客户端不需要持有密钥，由服务端通过密钥生成Token；

在JWT中，不应该在Payload里面加入任何敏感的数据，如用户密码等信息，因为payload并没有做加密，只是一个Base64的编码，
攻击者拿到token以后就可以得到用户敏感信息；


应用将JWT字符串作为该请求Cookie的一部分返回给用户。注意，在这里必须使用HttpOnly属性来防止Cookie被JavaScript读取，从而避免跨站脚本攻击（XSS攻击）。


Session方式存储用户id的最大弊病在于要占用大量服务器内存，对于较大型应用而言可能还要保存许多的状态。一般而言，大型应用还需要借助一些KV数据库和一系列缓存机制来实现Session的存储。

而JWT方式将用户状态分散到了客户端中，可以明显减轻服务端的内存压力。

虽说JWT方式让服务器有一些计算压力（例如加密、编码和解码），但是这些压力相比磁盘I/O而言或许是半斤八两。具体是否采用，需要在不同场景下用数据说话。


当用户点击了“注销”按钮，用户的令牌在Angular端会从授权认证服务AuthenticationService中移除，但是此令牌仍旧是有效的，还可以被攻击者窃取到，用于API调用，直至jsonwebtoken的有效时间结束。
为了避免此情况的发生，可以使用Redis数据库来存储已撤销的令牌——当用户点击注销按钮时。且令牌在Redis存储的时间与令牌在jsonwebtoken中定义的有效时间相同。当有效时间到了后，令牌会自动被Redis删除。
