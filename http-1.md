```
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
type ServeMux struct {
    mu    sync.RWMutex
    m     map[string]muxEntry
    hosts bool 
}

type muxEntry struct {
    explicit bool
    h        Handler
    pattern  string
}
```