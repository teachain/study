##测试##
在测试环境下只需要一个表示参数就可以生成各种分析文件

```
go test -cpuprofile=cpu.out
go test -blockprofile=block.out
go test -memprofile=mem.out

```



所有以_test.go为后缀名的源文件并不是go build构建包的一部分，它们是go test测试的一部分。

在*_test.go文件中，有三种类型的函数:

* 测试函数（一个测试函数是以Test为函数名前缀的函数，用于测试程序的一个逻辑行为是否正确。）

* 基准测试函数 （基准测试函数是以Benchmark为函数名前缀的函数，它们用来衡量一些函数的性能。）

* 示例函数（示例函数是以Example为函数名前缀的函数，提供一个由编译器保证正确性的示例文档。）


###测试函数###

每个测试函数都必须导入testing包

```
import "testing"

func TestName(t * testing.T){
}

```

在运行测试时，也就是go test时加-v,也就是go test -v ，我们就可以打印每个测试函数的名字和运行时间。

参数-run 对应一个正则表达式，只有测试函数名被它正确匹配的测试函数才会被go test测试命令运行。


###基准测试###

```
import "testing"

func BenchmarkName(b *testing.B){
}

```

* testing.B参数除了提供和* testing.T类似的方法，还有额外一些和性能测量相关的方法。它还提供了一个整数N,用于指定操作执行的循环次数。


用go test -bench=.或go test -bench=具体要测试的函数名 命令来运行一个基准测试。

再加一个 -benchmem命令行标志参数将在报告中包含内存的分配数据统计。


###示例函数###

以Example为函数名开头，没有函数参数和返回值。






