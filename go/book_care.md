内存泄漏使用工具memprof来定位问题

提供性能靠cpuprof和自己做的一些统计信息来定位问题

GOGCTRACE这个环境变量可以开启gc调试信息的打印

go的gc是固定每两分钟执行一次，每次执行都是暂停整个程序的，300多毫秒应该足以导致可感受到的响应延迟。

对象多了会增加gc负担，导致gc时间过长

<font color="red">特别需要注意：</font>

return xxx 这一句语句并不是一条原子指令！

整个return过程，没有defer之前是，先在栈中写一个值，这个值被会当

作返回值。然后再调用RET指令返回。return xxx语句汇编后是先给

<font color="red">返回值</font>赋值，再做一个空的return: ( 赋值指令 ＋ RET指令),也就是对return 语句进行了拆分。

defer的执行是被插入到return指令之前的

有了defer之后，就变成了 (赋值指令 + CALL defer指令 + RET指令)

