##http##

从服务器端来看，我们想要将cookies设置到客户端，那么可以考虑以下方法

* 第一种方式

```
ck:=&http.Cookie{
  Name:"mycookie",
  Value:"hello",
  Path:"/",
  Domain:"localhost",
  MaxAge:120,
}
http.setCookie(w,ck)
```

* 第二种方式

```
ck:=&http.Cookie{
  Name:"mycookie",
  Value:"hello",
  Path:"/",
  Domain:"localhost",
  MaxAge:120,
}
w.Header().Set("Set-Cookie",ck.String())

w为http.ResponseWriter

```
当Value里有空格，换行符之类的，会造成错误的，需要自定处理。

服务器端读取请求方携带上来的cookie,则只需要

```
ck,err:=r.Cookie("mycookie")

mycookie:=ck.Value

```

 