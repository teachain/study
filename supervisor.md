### EVENT

这是事件的抽象类型，supervisor不会将这个事件发送给监听器的，但当使用这个事件配置在

```
[eventlistener:你自己实现的监听器]
events=EVENT
```

这个部分的时候，那么所有继承于EVENT的事件发生的时候，都会通知到你实现的监听器。



### PROCESS_STATE

它也是一个抽象类型，同样遵循继承的原则，也就是当把它配置在eventlistener里时，所有继承于它的事件发生的时候，都会通知到你实现的监听器。







