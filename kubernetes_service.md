

# Service

Kubernetes中的Service 是一个抽象的概念，它定义了Pod的逻辑分组和一种可以访问它们的策略，这组Pod能被Service访问，使用YAML或JSON 来定义Service，Service所针对的一组Pod通常由*LabelSelector*实现。

可以通过type在ServiceSpec中指定一个需要的类型的 Service，Service的四种type:

- *ClusterIP*（默认） - 在集群中内部IP上暴露服务。此类型使Service只能从群集中访问。
- *NodePort* - 通过每个 Node 上的 IP 和静态端口（NodePort）暴露服务。NodePort 服务会路由到 ClusterIP 服务，这个 ClusterIP 服务会自动创建。通过请求 <NodeIP>:<NodePort>，可以从集群的外部访问一个 NodePort 服务。
- *LoadBalancer* - 使用云提供商的负载均衡器（如果支持），可以向外部暴露服务。外部的负载均衡器可以路由到 NodePort 服务和 ClusterIP 服务。
- *ExternalName* - 通过返回 `CNAME` 和它的值，可以将服务映射到 `externalName` 字段的内容，没有任何类型代理被创建。这种类型需要v1.7版本或更高版本`kube-dnsc`才支持。

#### Kubernetes中的nodePort，targetPort，port的区别和意义

这是对内或对外都可以服务的

```
apiVersion: v1
kind: Service
metadata:
 name: nginx-service
spec:
 type: NodePort
 ports:
 - port: 30080
   targetPort: 80
   nodePort: 30001
 selector:
  name: nginx-pod
```

```
这里只能对内服务，外部不能访问到。
apiVersion: v1
kind: Service
metadata:
 name: mysql-service
spec:
 ports:
 - port: 33306
   targetPort: 3306
 selector:
  name: mysql-pod
```

## 1. nodePort

外部机器可访问的端口。 
比如一个Web应用需要被其他用户访问，那么需要配置`type=NodePort`，而且配置`nodePort=30001`，那么其他机器就可以通过浏览器访问scheme://node:30001访问到该服务，例如[http://node:30001](http://node:30001/)。 
例如MySQL数据库可能不需要被外界访问，只需被内部服务访问，那么不必设置`NodePort`

## 2. targetPort

容器的端口（最根本的端口入口），与制作容器时暴露的端口一致（DockerFile中EXPOSE）

targetPort很好理解，targetPort是pod上的端口，从port和nodePort上到来的数据最终经过kube-proxy流入到后端pod的targetPort上进入容器。

## 3. port

kubernetes中的服务之间访问的端口，尽管mysql容器暴露了3306端口，但是集群内其他容器需要通过33306端口访问该服务，外部机器不能访问mysql服务，因为他没有配置NodePort类型。



需要启动kube-proxy服务。

总的来说，port和nodePort都是service的端口，前者暴露给集群内客户访问服务，后者暴露给集群外客户访问服务。从这两个端口到来的数据都需要经过反向代理kube-proxy流入后端pod的targetPod，从而到达pod上的容器内。