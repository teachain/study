https://www.cnblogs.com/yjf512/category/385369.html

## web 服务器

如果你的go程序是用http包启动的web服务器，你想查看自己的web服务器的状态。这个时候就可以选择net/http/pprof。你只需要引入包_"net/http/pprof"，然后就可以在浏览器中使用<http://localhost:port/debug/pprof/>直接看到当前web服务的状态，包括CPU占用情况和内存使用情况等。

## 服务进程

如果你的go程序不是web服务器，而是一个服务进程，那么你也可以选择使用net/http/pprof包，同样引入包net/http/pprof，然后在开启另外一个goroutine来开启端口监听。

比如：

```
go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil)) 
```

}()

## 应用程序

如果你的go程序只是一个应用程序，比如计算fabonacci数列，那么你就不能使用net/http/pprof包了，你就需要使用到runtime/pprof。具体做法就是用到pprof.StartCPUProfile和pprof.StopCPUProfile。比如下面的例子：

```
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
```

…

运行程序的时候加一个--cpuprofile参数，比如fabonacci --cpuprofile=fabonacci.prof

这样程序运行的时候的cpu信息就会记录到XXX.prof中了。

下一步就可以使用这个prof信息做出性能分析图了（需要安装graphviz）。

使用go tool pprof (应用程序) （应用程序的prof文件）

进入到pprof，使用web命令就会在/tmp下生成svg文件，svg文件是可以在浏览器下看的。

```
sudo apt-get install xrdp vnc4server xbase-clients

sudo apt-get install dconf-editor
```

