##字符串##

字符串是一个不可改变的字节序列。

内置的len函数可以返回一个字符串中的字节数目（不是rune字符数目）。


字符串在Go语言内存模型中用一个2<font color="red">字长</font>的数据结构表示。它包含一个<font color="red">指向字符串存储结构的指针</font>和一个<font color="red">长度</font>数.