##接口##

An interface variable can store any concrete (non-interface) value as long as that value implements the interface's methods.

当你拿到一个<font color="red">具体的类型</font>时你就知道它的本身是什么和你可以用它来做什么。

而接口类型是一种<font color="red">抽象的类型</font>,它不会暴露出它所代表的对象的内部值得结构和这个对象支持的基础操作的集合，它们只会展示出它们自己的方法。也就是说<font color="red">当你看到一个接口类型的值时，你不知道它是什么，唯一知道的就是可以通过它的方法来做什么。</font>

接口通常用在函数或方法的参数中。也就是作为形参这种情况比较多。

比如fmt.FprintF(w io.Writer,format string,args...interface{})(int,error)这个函数，w的类型io.Writer就是一个接口类型。

<font color="red">里氏替换原则：一个类型可以自由的使用另一个满足相同接口的类型来进行替换。</font>


接口类型具体描述了一系列方法的集合，一个实现了这些方法的具体类型是这个接口类型的实例。

概念上讲一个接口的值，<font color="red">接口值由两个部分组成，一个具体的类型和那个类型的值，它们被称为接口的动态类型和动态值。</font>

###类型断言###

类型断言是一个使用在接口值上的操作。语法上它看起来像<font color="red">x.(T),被称为类型断言</font>，这里的x表示一个接口的类型和T表示一个类型（具体类型或接口类型）。一个类型断言检查它操作对象的动态类型是否和断言的类型匹配。

###类型开关###

一个类型开关像普通的switch语句一样，它的运算对象是<font color="red">x.(type)</font>,它使用了关键词字面量type,并且每个case有一到多个类型。

```
   switch x.(type){
     	case nil:
     	   //todo
     	case int ,uint:
     	   //todo
     	case bool:
     	   //todo
     	case string:
     	   //todo
     	default:
     	   //todo
   }
```

接口只有当有两个或两个以上的具体类型必须以相同的方式进行处理时才需要。

























