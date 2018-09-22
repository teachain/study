##Docker监控技术基础##

<font color="red">只要是这四个部分：CPU、内存、磁盘、网络。</font>

* docker stats
* Remote API
* 伪文件系统

<font color="red">最好的方式应该就是使用伪文件系统。</font>

大概位置

```
/sys/fs/cgroup/{memory,cpuacct,blkio}/system.slice/${docker ps --no-trunc}.scope

```


##API##

使用API的话需要注意一下，那就是不要给Docker daemon带来性能负担。如果你有一台主机有200个容器，如果非常频繁的采集系统性能可能会大量占据CPU时间。


###内存(Memory)###

内存的很多性能指标都来自于 memory.stat 文件

Putting everything together to look at the memory metrics for a Docker container, take a look at

```
 /sys/fs/cgroup/memory/docker/<longid>/

```


###处理器(CPU)###

* cpuacct.stat 文件
* docker.cpu.system
* docker.cpu.user

###磁盘(IO)###

* blkio.throttle.io_service_bytes ，读写字节数 

* blkio.throttle.io_serviced ，读写次数

文件

```
/sys/fs/cgroup/blkio/docker/容器id/blkio.throttle.io_service_bytes

```
从这里获取的数值是<font color="red">io的操作字节</font>，是实际操作，而<font color="red">并非是docker的io限制</font>

文件

```
/sys/fs/cgroup/blkio/docker/容器id/blkio.throttle.io_serviced 
```

blkio.throttle.io_serviced 是操作次数，实际操作而非限制

* 读取
* 写入
* 同步
* 异步


###网络###

/sys/class/net/veth559b656/statistics


| 名称        | 描述           | 指标类型  |
| ------------- |:-------------:| -----:|
| Bytes     | 网络流量发送和接收字节大小 | 资源利用率（标准指标） |
|Packets     | 网络封包发送和接收计数器      |   资源利用率（标准指标） |

```
#指定容器id
$containerId=""

#取得容器进程id
$containerPid="docker inspect -f '{{.State.Pid}}' $containerId"

#查看网络数据
cat /proc/$containerPid/net/dev

```

网络从 1.6.1版本以后才支持，和以上的路径有所不同，获取使用容器Pid获取，<font color="red">注意Host模式获取的是主机网络数据，所以 host 模式无法从容器数据统计网络数据.</font>







