

http的源码

```
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

// Handle registers the handler for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func Handle(pattern string, handler Handler) { DefaultServeMux.Handle(pattern, handler) }

// HandleFunc registers the handler function for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}
```

这里我们很明显地看到Handler是个接口。而调用http.Handle的时候，第二个参数要提供一个Handler接口值，http.HandleFunc调用的时候，第二个参数就是一个func(ResponseWriter, *Request)，这是他们之间的区别。