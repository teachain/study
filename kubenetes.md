## Kubernetes节点

每个节点上都要运行Docker。Docker负责所有具体的映像下载和容器运行。

一些有用的链接

```
https://github.com/coreos/flannel         flannel网络

```



### 健康检查

#### 使用Liveness及Readness探针

- Liveness探针主要用于判断Container是否处于运行状态，比如当服务crash或者死锁等情况发生时，kubelet会kill掉Container, 然后根据其设置的restart policy进行相应操作（可能会在本机重新启动Container，或者因为设置Kubernetes QoS，本机没有资源情况下会被分发的其他机器上重新启动）。
- Readness探针主要用于判断服务是否已经正常工作，如果服务没有加载完成或工作异常，服务所在的Pod的IP地址会从服务的endpoints中被移除，也就是说，当服务没有ready时，会将其从服务的load balancer中移除，不会再接受或响应任何请求。



#### 使用建议

1.  建议对全部服务同时设置服务（readiness）和Container（liveness）的健康检查
2. 通过TCP对端口检查形式（TCPSocketAction），仅适用于端口已关闭或进程停止情况。因为即使服务异常，只要端口是打开状态，健康检查仍然是通过的。
3. 基于第二点，一般建议用ExecAction自定义健康检查逻辑，或采用HTTP Get请求进行检查（HTTPGetAction）。
4.  无论采用哪种类型的探针，一般建议设置检查服务（readiness）的时间短于检查Container（liveness）的时间，也可以将检查服务（readiness）的探针与Container（liveness）的探针设置为一致。目的是故障服务先下线，如果过一段时间还无法自动恢复，那么根据重启策略，重启该container、或其他机器重新创建一个pod恢复故障服务。



## master 服务启动（使用flanneld）

```
for SERVICES in etcd kube-apiserver kube-controller-manager kube-scheduler flanneld; do
    systemctl restart $SERVICES
    systemctl enable $SERVICES
    systemctl status $SERVICES
done
```

## node服务启动

```
for SERVICES in kube-proxy kubelet flanneld docker; do
    systemctl restart $SERVICES
    systemctl enable $SERVICES
    systemctl status $SERVICES
done
```



## flannel网络

```
flannel是CoreOS提供用于解决Dokcer集群跨主机通讯的覆盖网络工具。它的主要思路是：预先留出一个网段，每个主机使用其中一部分，然后每个容器被分配不同的ip；让所有的容器认为大家在同一个直连的网络，底层通过UDP/VxLAN等进行报文的封装和转发。
Flannel的设计目的就是为集群中的所有节点重新规划IP地址的使用规则，从而使得不同节点上的容器能够获得“同属一个内网”且”不重复的”IP地址，并让属于不同节点上的容器能够直接通过内网IP通信。
```

![flannel-01](/Users/daminyang/github/study/flannel-01.png)

1. 容器直接使用目标容器的ip访问，默认通过容器内部的eth0发送出去。
2. 报文通过`veth pair`被发送到`vethXXX`。
3. `vethXXX`是直接连接到虚拟交换机`docker0`的，报文通过虚拟`bridge docker0`发送出去。
4. 查找路由表，外部容器ip的报文都会转发到`flannel0`虚拟网卡，这是一个`P2P`的虚拟网卡，然后报文就被转发到监听在另一端的`flanneld`。
5. `flanneld`通过`etcd`维护了各个节点之间的路由表，把原来的报文`UDP`封装一层，通过配置的`iface`发送出去。
6. 报文通过主机之间的网络找到目标主机。
7. 报文继续往上，到传输层，交给监听在8285端口的`flanneld`程序处理。
8. 数据被解包，然后发送给`flannel0`虚拟网卡。
9. 查找路由表，发现对应容器的报文要交给`docker0`。
10. `docker0`找到连到自己的容器，把报文发送过去。

flannel服务需要先于Docker启动。flannel服务启动时主要做了以下几步的工作：

- 从etcd中获取network的配置信息。
- 划分subnet，并在etcd中进行注册。
- 将子网信息记录到`/run/flannel/subnet.env`中。



## Volume

#### EmptyDir

EmptyDir是一个空目录，它的生命周期和所属的pod是完全一致的。它的用处是：可以在同一个pod内的不同容器之间共享工作过程中产生的文件。缺省情况下，EmptyDir 是使用主机磁盘进行存储的，也可以设置emptyDir.medium 字段的值为Memory，来提高运行速度，但是这种设置，对该卷的占用会消耗容器的内存份额。

#### HostPath

这种会把宿主机上的指定卷加载到容器之中，当然，如果 Pod 发生跨主机的重建，其内容就难保证了。



#### configmap

configmap的创建

```
使用kubectl create configmap -h 查看使用方式
```

configmap的使用

```
1、通过环境变量
      env:
        - name: 你定义的环境变量的名字
          valueFrom:
            configMapKeyRef:
              name: 你创建的configmap的名字
              key: 你configmap中定义的key
              
2、通过在pod的命令行下运行的方式
   在命令行下引用时，需要先设置为环境变量，之后 可以通过$(VAR_NAME)设置容器启动命令的启动参数
   command: [ "/bin/sh", "-c", "echo $(SPECIAL_LEVEL_KEY) $(SPECIAL_TYPE_KEY)" ]
      env:
        - name: SPECIAL_LEVEL_KEY
          valueFrom:
            configMapKeyRef:
              name: special-config
              key: special.how
              
3、使用volume的方式挂载入到pod内
	使用volume将ConfigMap作为文件或目录直接挂载，其中每一个key-value键值对都会生成一个文件，key为文件名，      value为内容
	 command: [ "/bin/sh", "-c", "cat /etc/config/special.how" ] //注意最后这里
      volumeMounts:
      - name: config-volume
        mountPath: /etc/config
  volumes:
    - name: config-volume
      configMap: //注意这里
        name: special-config //这里显然是你创建configmap的名字
  restartPolicy: Never
  
  最后需要说明两点：

1、ConfigMap必须在Pod之前创建

2、只有与当前ConfigMap在同一个namespace内的pod才能使用这个ConfigMap，换句话说，ConfigMap不能跨命名空间调用。
```



### secret

创建方式类似configmap,只不过value需要你事先将它进行base64编码

引用方式也是

```
         valueFrom:
 
           secretKeyRef:  //就这里和configMapRef类似的
 
             name: test-secret
 
             key: username
```



### 

## DaemonSet

DaemonSet能够让所有（或者一些特定）的Node节点运行同一个pod。当节点加入到kubernetes集群中，pod会被（DaemonSet）调度到该节点上运行，当节点从kubernetes集群中被移除，被（DaemonSet）调度的pod会被移除，如果删除DaemonSet，所有跟这个DaemonSet相关的pods都会被删除。

Daemon Sets就是让一个pod在所有的k8s集群节点上都运行一个

```
在使用kubernetes来运行应用时，很多时候我们需要在一个区域（zone）或者所有Node上运行同一个守护进程（pod），例如如下场景：

每个Node上运行一个分布式存储的守护进程，例如glusterd，ceph
每个Node上运行日志采集器，例如fluentd，logstash
每个Node上运行监控的采集端，例如Prometheus Node Exporter,collectd等
```

