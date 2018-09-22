##数组##

go的数组类型由两部分组成——类型和长度，二者缺一不可。数组本来就是一块存储相同类型元素的连续内存空间，因此决定一个数组的类型，必然需要决定其存储元素的类型以及存储多少个元素。



##切片slice##

slice由三个部分组成：指针，长度，容量。

一个slice是一个数组某个部分的引用。

在内存中，它是一个包含3个域的结构体

```
type slice struct {
    array unsafe.Pointer//指向slice中第一个元素的指针
    len   int //长度
    cap   int //容量
}
```

* 指向slice中第一个元素的指针

* slice的长度，长度是下标操作的上界，若x[i]中<font color="red">i必须小于长度（len()）也就是说下标index必须永远 index < len() </font>


也就是说以下代码将会报错

```
	s := make([]string, 5, 10)
	//panic: runtime error: index out of range
	s[6] = "hello"
```

* slice的容量，容量是分割操作的上界，如x[i:j]中<font color="red">j不能大于容量(cap()),也就是j<=cap()</font>



```
package main

import (
	"fmt"
)

func main() {
   //len(x)==5,cap(x)==5
	x := []int{2, 3, 5, 7, 11}
	y := x[1:3] //执行该操作后，y的pointor就指向了x中的元素3的这个位置
	z := y[1:cap(y)]
	fmt.Println("len(x)=", len(x), "cap(x)=", cap(x), x)
	fmt.Println("len(y)=", len(y), "cap(y)=", cap(y), y)
	fmt.Println("len(z)=", len(z), "cap(z)=", cap(z), z)
}
```
很有意思的输出

```
len(x)= 5 cap(x)= 5 [2 3 5 7 11]
len(y)= 2 cap(y)= 4 [3 5]
len(z)= 3 cap(z)= 3 [5 7 11]
```

注意点:

* <font color="red">append无论如何都是向slice的尾部追加数据</font>

* <font color="red">要时刻记得slice中的len和cap是会动态变化的。</font>

当你进行append操作的时候务必理解以下情况

* <font color="red">append操作的位置总是当前的len+1（始终跟着len）</font>

* 当len<cap的时候，不会重新分配内存，直接将数据存入slice中。

* 当len==cap的时候，将重新分配内存，然后再将数据存入到slice中。

<font color="red">每次cap改变的时候指向array内存的指针都在变化。当在使用 append 的时候，如果 cap==len 了这个时候就会新开辟一块更大内存，然后把之前的数据复制过去。
实际go在append的时候放大cap是有规律的。在 cap 小于1024的情况下是每次扩大到 2 * cap ，当大于1024之后就每次扩大到 1.25 * cap 。</font>

##创建slice##

声明一个Array通常使用 make ，可以传入2个参数，也可传入3个参数，第一个是数据类型，第二个是 len ，第三个是cap 。 <font color="red">如果不穿入第三个参数，则 cap=len</font> ，append 可以用来向数组末尾追加数据。

##创建子slice##

分割原则

* x[i:j]中的i < len()

* x[i:j]中j不能大于容量(cap()),也就是j<=cap()

当创建子切片的时候，一定要注意它引用的底层数组，我们可以认为子切片的cap就等于引用的底层数组的大小与子切片的数组指针指向的元素的索引之差。
