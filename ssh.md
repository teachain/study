#### 作为用户直接操作的机器（敲键盘的机器，我们称它为主机A）

1、配置代理转发

打开/etc/ssh/ssh_config配置文件（注意不是sshd_config文件，不要搞混了）

```
vi /etc/ssh/ssh_config
```

找到ForwardAgent no，把no改为yes,把前面的#这个注释符号去掉。

2、启动代理

```
ssh-agent bash
```

3、将私钥添加到ssh代理中

```
ssh-add ~/.ssh/id_rsa
```

#### 作为代理的机器(主机B)

我们需要用这台机器来转发

1、修改配置/etc/ssh/sshd_config(因为是ssh服务器端扮演了中间代理的角色，所以我们改的是sshd)

```
vim /etc/ssh/sshd_config #一定不要和ssh_config搞混了
```

将#AllowAgentForwarding yes 前面的注释符号#去掉

作为代理的机器(主机c)要执行同样的操作。



前提是将主机A的id_rsa.pub都已经拷贝到了主机B和主机C。