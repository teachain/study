#### 上传和下载内容是Swarm的存在理由。

必备条件：

已经运行swarm客户端，并且有一个正在运行的geth节点，默认情况下会在端口8500上进行侦听。

###### 上传和下载数据:

可通过以下任选一种方式进行

1. swarm 终端上的命令行界面(CLI)
2. 通过HTTP接口进行`http://localhost:8500`。

swarm 终端上的命令行界面(CLI) 上传文件

```
swarm up ./smart.sol  #也就是swarm up 文件路径，即可上传一个文件到swarm节点，并同步给其它节点
```

通过这样上传一个文件之后，它会返回一个表示该资源的一个hash给你。你用这个hash就可以访问到这个文件

假设./smart.sol上传之后，返回的hash为：

5f79254ac0f0b4f5fe2abcc79ec2aaa8c2a16d88ab8083f089ee22ead5edf3d8

swarm 终端上的命令行界面(CLI) 下载文件

```
swarm down bzz:/5f79254ac0f0b4f5fe2abcc79ec2aaa8c2a16d88ab8083f089ee22ead5edf3d8
```

这个命令执行之后，就会把上述的./smart.sol这个文件的内容下载到

5f79254ac0f0b4f5fe2abcc79ec2aaa8c2a16d88ab8083f089ee22ead5edf3d8文件中

以hash值来命名这个文件。内容会跟./smart.sol一模一样。

```
将文件上载到本地Swarm节点后，您的节点将同步数据块与网络上的其他节点。因此，即使原始节点脱机，该文件最终也将在网络上可用。
```

使用`--recursive`标志实现上载目录。