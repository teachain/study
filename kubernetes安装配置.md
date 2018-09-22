





kubernetes apiServer(master，主节点)

minion(从节点)

1. 安装Open vSwitch及配置GRE

   为了解决跨minion之间Pod的通信问题，我们在每个minion上安装Open vSwtich，并使用GRE或者VxLAN使得跨机器之间Pod能相互通信，本文使用GRE，而VxLAN通常用在需要隔离的大规模网络中。

2. 根据Kubernetes的设计架构，需要在minion上部署docker, kubelet, kube-proxy

- etcd：存储flannel相关信息；


- flannel：基于etcd的容器网络管理工具；


- docker：容器引擎；



java 已经安装以下版本：

```
java version "1.8.0_144"
Java(TM) SE Runtime Environment (build 1.8.0_144-b01)
Java HotSpot(TM) 64-Bit Server VM (build 25.144-b01, mixed mode)
```

maven已经安装以下版本：

```
Apache Maven 3.5.0 (ff8f5e7444045639af65f6095c62210b5713f426; 2017-04-04T03:39:06+08:00)
Maven home: /usr/local/apache-maven-3.5.0
Java version: 1.8.0_144, vendor: Oracle Corporation
Java home: /usr/java/jdk1.8.0_144/jre
Default locale: en_US, platform encoding: UTF-8
OS name: "linux", version: "3.10.0-514.el7.x86_64", arch: "amd64", family: "unix"
```

etcd 安装在/usr/bin/目录下,包括etcd和etcdctl这两个可执行文件。

1. 安装mysql

   ```
   wget http://dev.mysql.com/get/mysql-community-release-el7-5.noarch.rpm
   rpm -ivh mysql-community-release-el7-5.noarch.rpm
   yum install mysql-community-server
   ```

2. 重启mysql服务

   ```
   service mysqld restart
   ```

3. 修改密码

   ```
   set password for 'root'@'localhost' =password('password');
   ```

4. 远程连接

   把在所有数据库的所有表的所有权限赋值给位于所有IP地址的root用户。

   ```
   grant all privileges on *.* to root@'%'identified by 'password';
   ```

5. 创建用户

   ```
   create user 'username'@'%' identified by 'password';  
   ```

   ​



  