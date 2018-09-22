

# linux网络设置

1、运行以下命令可以查看网卡和ip信息：

```
ifconfig -a
```

eth0表示第一块网卡

2、在虚拟机中使用桥接模式以后，运行以下命令可使虚拟机可以随机获取一个ip地址：

```
dhcient
```

3、网卡配置文件

```
vim /etc/sysconfig/network-scripts/ifcfg-eth0
```

```
ONBOOT=yes
BOOTPROTO=static
IPADDR=192.168.1.109
NETMASK=255.255.255.0
GATEWAY=192.168.1.1
DNS1=192.168.1.1
DNS2=8.8.8.8
```

4、保存退出之后，要重启network服务

```
service network restart
```

5、关闭网卡

```
ifdown eth0
```

6、开启网卡

```
ifup  eth0
```


在CentOS 7或RHEL 7或Fedora中防火墙由firewalld来管理

```
firewall-cmd --permanent --add-port=8090/tcp
```

用以下命令来查询端口是否开放

```
firewall-cmd --permanent --query-port=1000/tcp
```

