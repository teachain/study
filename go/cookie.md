##Cookie##
我们先从net/http/cookie.go这个文件来看一下cookie的定义

```
type Cookie struct {
	Name  string
	Value string

	Path       string    // optional
	Domain     string    // optional
	Expires    time.Time // optional
	RawExpires string    // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}

```

然后我们看往客户端写入cookie的方法（server->client）

```

func SetCookie(w ResponseWriter, cookie *Cookie) {
	if v := cookie.String(); v != "" {
		w.Header().Add("Set-Cookie", v)
	}
}
```

从这个方法的实现，我们可以知道,我们只要类似这样子定义一个cookie,然后就可以调用http.SetCookie方法就可以往客户端写入cookie，而且http.SetCookie是可以调用多次的，也就是可以写入多个不同的Cookie。

```
expitation := time.Now().Add(24 * time.Hour)

var username string

cookie := http.Cookie{Name: "username", Value: username, Expires: expitation}

http.SetCookie(w, &cookie)

```

调用http.Request的Cookies()方法就可以读取到客户端携带过来的所有Cookie,这个方法返回的是[]*Cookie，我们可以得知它是一个slice,我们只要迭代就可以得到每一个Cookie了。

我们从http协议里得知

* 当服务器往客户端写入cookie的时候，是在<font color="red">响应头</font>里使用<font color="red">Set-Cookie: </font>
* 当客户端携带cookie提交给服务器的时候，是在<font color="red">请求头</font>里带上了<font color="red">Cookie:</font>

