## go plugin

1、定义plugin(pluginhello.go),它的package必须是main

```
package main

import (
    "fmt"
)

func Hello() {
    fmt.Println("Hello World From Plugin!")
}
```

2、把它编译成动态库

```
go build --buildmode=plugin -o pluginhello.so pluginhello.go
```

3、加载plugin

```
package main

import (
    "fmt"
    "os"
    "plugin"
)

func main() {
    p, err := plugin.Open("./pluginhello.so")
    if err != nil {
        fmt.Println("error open plugin: ", err)
        os.Exit(-1)
    }
    s, err := p.Lookup("Hello")
    if err != nil {
        fmt.Println("error lookup Hello: ", err)
        os.Exit(-1)
    }
    if hello, ok := s.(func()); ok {
        hello()
    }
}
```

