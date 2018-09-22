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

