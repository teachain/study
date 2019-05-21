# raft

http://thesecretlivesofdata.com/raft/  动画演示

```
All our nodes start in the follower state
```

```
If followers don't hear from a leader then they can become a candidate.
```

```
The candidate then requests votes from other nodes.
```

```
Nodes will reply with their vote
```

```
The candidate becomes the leader if it gets votes from a majority of nodes.
```

```
This process is called Leader Election.     1.选举
```

```
All changes to the system now go through the leader.
```

```
Each change is added as an entry in the node's log.
```

```
This log entry is currently uncommitted so it won't update the node's value.
```

```
To commit the entry the node first replicates it to the follower nodes
```

```
then the leader waits until a majority of nodes have written the entry.
```

```
The entry is now committed on the leader node and the node state is "5".  
//我们举例子说设置一个值为5这样子一个操作。
```

```
The leader then notifies the followers that the entry is committed.
```

```
The cluster has now come to consensus about the system state.
```

```
This process is called Log Replication.   2、日志复制
```



对于raft算法，raft算法的的容错只支持容错故障节点，不支持容错作恶节点。什么是故障节点呢？就是节点因为系统繁忙、宕机或者网络问题等其它异常情况导致的无响应，出现这种情况的节点就是故障节点。那什么是作恶节点呢？作恶节点除了可以故意对集群的其它节点的请求无响应之外，还可以故意发送错误的数据，或者给不同的其它节点发送不同的数据，使整个集群的节点最终无法达成共识，这种节点就是作恶节点。

raft 算法只支持容错故障节点，假设集群总节点数为n，故障节点为 f ，根据小数服从多数的原则，集群里正常节点只需要比 f 个节点再多一个节点，即 f+1 个节点，正确节点的数量就会比故障节点数量多，那么集群就能达成共识。因此 raft 算法支持的最大容错节点数量是（n-1）/2。

n=2f+1,所以至少需要3个节点。



一个团队一定会有一个老大和普通成员。对于 raft 算法，共识过程就是：只要老大还没挂，老大说什么，我们（团队普通成员）就做什么，坚决执行。那什么时候重新老大呢？只有当老大挂了才重选老大，不然生是老大的人，死是老大的鬼。