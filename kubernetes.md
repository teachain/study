1、Kubernetes是一个分布式架构，可灵活地进行安装和部署，可以部署在单机，也可以分布式部署，机器可以是物理机，也可以是虚拟机，但是需要运行linux(x86_64)操作系统。

注意：从上面的定义来看，我们就可以选择在一台机器上玩，也可以选择在多台机器上联合起来一台玩。也即是说它需要的这几个软件，可以安装在一台机器上，也可以安装在不同的机器上，只要能够联网，就能够保证他们能够协作工作。

2、所需要的几个软件为

- Kubernetes
- Docker
- Etcd
- Flannel
- Open vSwitch

3、Kubernetes依赖于Etcd,所以需要先部署Etcd,Etcd的github地址是

https://github.com/coreos/etcd/

我们使用以下命令来下载Etcd:#是系统的提示符号啦，因为我习惯切换到root安装软件，所以命令提示是#

```
# wget https://github.com/coreos/etcd/releases/download/v3.2.5/etcd-v3.2.5-linux-amd64.tar.gz
# tar xzvf etcd-v3.2.5-linux-amd64.tar.gz
# cd etcd-v3.2.5-linux-amd64
# cp etcd /usr/bin/etcd 
# cp etcdctl /usr/bin/etcdctl
```

从上面的安装过程我们可以看出，Etcd的发行包里就两个主要的可执行文件:

- etcd
- etcdctl

4、执行以下命令来运行etcd，注意url要用单引号包含起来，多个url用逗号分隔

```
# etcd -name etcd \
-data-dir /var/lib/etcd \
-listen-client-urls 'http://0.0.0.0:2379, http://0.0.0.0:4001' \
-advertise-client-urls 'http://0.0.0.0:2379, http://0.0.0.0:4001' \
>>/var/log/etcd.log 2>&1 &
```

5、检查它的健康状态

```
# etcdctl -C http://localhost:4001 cluster-health
```

当出现类型以下内容时，表示健康

```
member 8e9e05c52164694d is healthy: got healthy result from http://0.0.0.0:2379
cluster is healthy
```

6、使用以下命令查看监听端口

```
netstat -atnup|grep LISTEN
```



kubernetes提供了3中应用部署策略



Pod是容器的集合，容器是真正的执行体。Pod的设计并不是为了运行同一个应用的多个实例，而是运行一个应用多个紧密联系的程序，而每个程序运行在单独的容器中，以Pod的形式组合成一个应用。



##### Pod调度

Pod的调度指的是Pod在创建之后分配到哪一个Node上，调度算法分为两个步骤，第一个筛选出符合条件的Node，第二步选择最优的Node.

可通过在创建Pod的时候指定nodeSelector或nodeName来指定调度到符合条件的node,推荐使用nodeSelector来指定调度到哪个Node.



https://kubernetes.io/docs/tasks/access-application-cluster/communicate-containers-same-pod-shared-volume/



kubernetes api 文档

https://kubernetes.io/docs/api-reference/v1.6/