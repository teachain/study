### 使用本地代码库

Go现在通过 `go.mod` 文件来配置模块加载。

默认使用 github.com/zhouzme/snail-go 包会到 github 上去下载，但这个包还在本地开发中并未push到线上，那么可以通过 replace 配置来重定向当前项目对该包的加载路径：

```
replace github.com/zhouzme/snail-go => E:\Go\snail-go
```

这里 `E:\Go\snail-go` 为本地包的绝对路径，这样写就可以了，当本地`snail-go`包代码修改后就可以在当前项目看到实时效果了，注意中间符号是 `=>`，这是特别需要注意的。