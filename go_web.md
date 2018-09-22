##关于cookie##
HTTP协议是无状态的协议。一旦数据交换完毕，客户端与服务器端的连接就会关闭，再次交换数据需要建立新的连接。这就意味着服务器无法从连接上跟踪会话。要跟踪该会话，必须引入一种机制。

Cookie就是这样的一种机制。它可以弥补HTTP协议无状态的不足。在Session出现之前，基本上所有的网站都采用Cookie来跟踪会话。

服务器给客户端们颁发一个通行证吧，每人一个，无论谁访问都必须携带自己通行证。这样服务器就能从通行证上确认客户身份了。这就是Cookie的工作原理。
Cookie实际上是一小段的文本信息。客户端请求服务器，如果服务器需要记录该用户状态，就使用response向客户端浏览器颁发一个Cookie。客户端浏览器会把Cookie保存起来。当浏览器再请求该网站时，浏览器把请求的网址连同该Cookie一同提交给服务器。服务器检查该Cookie，以此来辨认用户状态。服务器还可以根据需要修改Cookie的内容。

在go的web框架beego中，默认将产生的sessionId放在了cookie中，所以当web中通过beego.BConfig.WebConfig.Session.SessionOn = true开启了会话支持之后，那么通过控制器去取session中的数据的时候，都是针对sessionId的，也就是说，通过this.getSession(key)这个方法取到的数据是针对sessionId的，就是不同的用户，通过相同的key取到的数据都是针对个人设置的值。也就是其实是map[sessionId]Session 这个样子的。

从客户端这一端来看，就是浏览器在做http请求时，通过请求头中设置一个请求行
Cookie	sid=b00c93; roleId=123456
每一个cookie包含key和value,然后每个cookie之间用分号";"分割
这样客户端浏览器就可以将携带通行证（cookie）给服务器端了。

从服务器端这一端来看，就是响应给客户端的时候，在响应头中设置一个响应行
Set-Cookie sid=b00c93; roleId=123456
每一个cookie包含key和value,然后每个cookie之间用分号";"分割
这样服务器端就可以将携带通行证扔给客户端浏览器（或http客户端）
这样，客户端在解析响应头的时候就可以得到cookie了



