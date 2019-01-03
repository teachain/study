在节点发现时使用的数据字段

分别存在以下key:(+表示拼接,nodeId表示节点唯一码，它是一个字符串)

```
version  //表示数据的版本

n:+nodeId+:discover+:lastping //表示nodeId对应节点最后一个ping消息的时间

n:+nodeId+:discover+:lastpong //表示nodeId对应节点最后一个pong消息的时间

n:+nodeId+:discover+:findfail //表示nodeId对应节点查找失败次数

n:+nodeId+:discover+:localendpoint //表示nodeId对应节点的本地端点

n:+nodeId+:tickets //nodeId对应节点的数据tickets
```

