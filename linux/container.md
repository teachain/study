##linux##

Linux内核中提供了六种命名空间

| Namespace | 系统调用参数        | 隔离内容          |      |
| --------- | ------------- | ------------- | ---- |
| UTS       | CLONE_NEWUTS  | 主机名与域名        |      |
| IPC       | CLONE_NEWIPC  | 信号量、消息队列和共享内存 |      |
| PID       | CLONE_NEWPID  | 进程编号          |      |
| NETWORK   | CLONE_NEWNET  | 网络设备、网络栈，端口等等 |      |
| MOUNT     | CLONE_NEWNS   | 挂载点（文件系统）     |      |
| USER      | CLONE_NEWUSER | 用户和用户组        |      |

在同一个namespace下的进程可以感知彼此的变化，而对外界的进程一无所知。



## PID namespace

PID namespace隔离非常实用，它对进程PID重新标号，即两个不同namespace下的进程可以有同一个PID。每个PID namespace都有自己的计数程序。内核为所有的PID namespace维护了一个树状结构，最顶层的是系统初始时创建的，我们称之为root namespace。他创建的新PID namespace就称之为child namespace（树的子节点），而原先的PID namespace就是新创建的PID namespace的parent namespace（树的父节点）。通过这种方式，不同的PID namespaces会形成一个等级体系。所属的父节点可以看到子节点中的进程，并可以通过信号等方式对子节点中的进程产生影响。反过来，子节点不能看到父节点PID namespace中的任何内容。



linux内核为所有的pid命名空间维护了一个树状结构：最顶层的是系统初始化时创建的root namespace(根名空间)，再创建新的pid namespace.



#### 最重要的是

父节点可以看到子节点中的进程，并可以通过信号等方式对子节点中的进程产生影响。反过来，子节点不能看到父节点名空间的任何内容。



在Docker中，每个Container都是Docker Daemon的子进程，每个Container进程缺省都具有不同的PID名空间。通过名空间技术，Docker实现容器间的进程隔离。



一个PID Namespace为进程提供了一个独立的PID环境，PID Namespace内的PID将从1开始，在Namespace内调用fork，vfork或clone都将产生一个在该Namespace内独立的PID。新创建的Namespace里的第一个进程在该Namespace内的PID将为1，就像一个独立的系统里的init进程一样。该Namespace内的孤儿进程都将以该进程为父进程，当该进程被结束时，该Namespace内所有的进程都会被结束。PID Namespace是层次性，新创建的Namespace将会是创建该Namespace的进程属于的Namespace的子Namespace。子Namespace中的进程对于父Namespace是可见的，一个进程将拥有不止一个PID，而是在所在的Namespace以及所有直系祖先Namespace中都将有一个PID。系统启动时，内核将创建一个默认的PID Namespace，该Namespace是所有以后创建的Namespace的祖先，因此系统所有的进程在该Namespace都是可见的。



### 我们假设我们已经得到了容器的长id(就是64个字符很长的那个容器id)

下面我们就用containerId来指代容器id，真实代码里要用真实的容器id来替代containerId

1、在这样的情况下，我们可以读取 

```
/run/runc/containerId/state.json
```

这个文件，获取到该容器的进程id（pid）,rootfs等信息。

2、再使用

```
ps -ef
```

获得宿主机上所有的进程信息

3、对进程信息进行筛选，也就是进程的父进程(ppid)等于第一步获得的容器进程pid即是容器中的进程。



### CPU的占用率











