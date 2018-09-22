##反射##

var w io.Writer
w = r.(io.Writer)
The expression in this assignment is a type assertion; what it asserts is that the item inside r also implements io.Writer, and so we can assign it to w.

也就是说r是一个接口类型的变量，我们就可以使用类型断言，断言说r实现了io.Writer接口，并且分配给w.

<font color="red">Interfaces do not hold interface values.</font>


反射是由reflect包提供支持，它定义了两个重要的类型，Type和Value

一个Type表示一个Go类型，它是一个接口。

* 函数reflect.TypeOf接受任意的interface{}类型，并返回对应动态类型的reflect.Type。（<font color="red">因为reflect.TypeOf返回的是一个动态类型的接口值，它总是返回具体的类型</font>）


* 函数reflect.ValueOf接受任意的interface{}类型，并返回对应动态类型的reflect.Value。

<font color="red">知识点：将一个具体的值转为接口类型会有一个隐式的接口转换操作，它会创建一个包含两个信息的接口值：操作数的动态类型和它的动态的值。</font>

调用reflect.Value的Type方法将返回具体类型所对应的reflect.Type。

```
type MyInt int
var x MyInt=7
v:=reflect.ValueOf(x)
fmt.Println(v.Type())
fmt.Println(v.Kind())
```

从这里我们就可以把变量的类型分为三种类型，动态类型（dynamically typed），静态类型（static type），底层类型（underlying type）

Kind()方法返回的始终是底层类型，Type（）返回的是静态类型，就是字面上定义的类型，如上面的x的底层类型是int，它的静态类型是MyInt
也就是用Type()方法可以区分出MyInt和int类型。而kind并不能区分MyInt和int类型。




