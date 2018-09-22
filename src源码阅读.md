##HTTP##

implicit  隐形地

explicitly 明确地

deadlines  最后期限

gone away  消失

guarded  保护

1、http.ResponseWriter与http.Request操作注意点:

要坚持一个原则，先操作完http.Request，再去操作http.ResponseWriter，在文档里有提到，如果说你先去操作http.ResponseWriter，那么一旦触发了flush操作，那么http.Request的Body将变得不可用。

2、Hijacker接口可以接管http的connection

```
type Hijacker interface {
	// Hijack lets the caller take over the connection.
	// After a call to Hijack the HTTP server library
	// will not do anything else with the connection.
	//
	// It becomes the caller's responsibility to manage
	// and close the connection.
	//
	// The returned net.Conn may have read or write deadlines
	// already set, depending on the configuration of the
	// Server. It is the caller's responsibility to set
	// or clear those deadlines as needed.
	//
	// The returned bufio.Reader may contain unprocessed buffered
	// data from the client.
	Hijack() (net.Conn, *bufio.ReadWriter, error)
}
```
