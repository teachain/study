##命令行参数##

os.Args变量是一个字符串的切片，os.Args的第一个元素os.Args[0]是命令行本身的名字，其它的元素则是程序启动时传递给它的参数。

<font color="red">make 仅能用于slices,maps,channels的内存分配和初始化,这三种类型都是引用类型</font>,它的返回值是非指针类型的对象。也就是返回的类型是T,而不是*T

想看看结构体的类型以及字段名和值的话

```
fmt.Printf("%v",结构体变量) //结构体字段的值
fmt.Printf("%+v",结构体变量) //附带字段
fmt.Printf("%#v",结构体变量) //附带结构体类型
```

当然不仅仅是结构体类型，任何类型你都可以试一试。

常量取值 at compile time
变量取值 at run time