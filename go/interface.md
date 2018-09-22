##interface##

interface是一组method的组合。

interface类型定义了一组方法，如果某个对象实现了某个接口的所有方法，则此对象就实现了此接口。

一个对象可以实现任意多个interface

一个函数把interface{}作为参数，那么他可以接受任意类型的值作为参数，如果一个函数返回interface{},那么也就可以返回任意类型的值.


<font color="red">interface在内存上实际由两个成员组成</font>

* <font color="red">tab 指向虚表（虚表描述了实际的类型信息及该接口所需要的方法集）</font>
* <font color="red">data 指向实际引用的数据</font>

itable(tab指向)的结构，包括描述type信息的一些元数据，以及满足Stringer接口的函数指针列表。

golang是在运行时生成虚表的。

理解了golang的内存结构，再来分析诸如类型断言等情况的效率问题就很容易了，<font color="red">当判定一种类型是否满足某个接口时，golang使用类型的方法集和接口所需要的方法集进行匹配，如果类型的方法集完全包含接口的方法集，则可认为该类型满足该接口。</font>例如某类型有m个方法，某接口有n个方法，则很容易知道这种判定的时间复杂度为O(mXn)，不过可以使用预先排序的方式进行优化，实际的时间复杂度为O(m+n)。

