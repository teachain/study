##map##

<font color="red">map的底层是用hashmap实现的</font>

<font color="red">在可预测map中的元素数量的情况下，可以使用make(map[key]value,capacity)这种方式来预先分配好，也就是指定好capacity，可稳定提升性能。</font>

在go语言中，一个map就是一个哈希表的引用，map类型可以写为<font color="red">map[KeyType]ValueType</font>，map中所有的key都有相同的类型，所有的value也有着相同的类型，但是key和value之间可以是不同的数据类型。<font color="red">key必须是支持==比较运算符的数据类型</font>

<font color="red">map类型的零值是nil,向一个nil值的map存入元素将导致一个panic</font>,所以必须采用以下两种方式之一创建一个map

####创建一个map可以使用以下两种方式####

* <font color="red">使用内置函数make来创建一个map</font>
  
   ```
   ages:=make(map[string]int)
   ages["alice"]=18
   ages["david"]=22
   
   ```
   
* <font color="red">使用map字面值得语法来创建一个map</font>
 
 ```
   ages:=map[string]int{
      "alice":18,
      "david":22,
   }
 ```
 
 当然这里的一般形式就是
 
 ```
   ages := map[string]int{}
   ages["alice"] = 18
   ages["david"] = 22
 ```

####访问map中的元素####

<font color="red">通过key对应的下标语法访问</font>

```
ages["alice"]=32  //赋值(添加元素或修改元素值)

fmt.Println(ages["alice"]) //读取

```

<font color="red">使用内置的delete函数可以删除元素</font>

```
   delete (ages,"alice")
```

<font color="red">注意:不能对map的元素进行取址操作</font>

```
    p:=&ages["alice"] //错误的写法
```


map之间不能直接通过==和！=进行比较，map只能和nil使用==和！=进行比较，如果要判断两个map是否包含相同的key和value，我们必须通过一个循环实现。


<font color="red">从Go 1.6开始，Go运行时系统对字典的非并发安全访问采取零容忍的态度</font>









