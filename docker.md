##Docker##
<font color="red">Docker背后的想法是创建软件程序可移植的轻量容器，让其可以在任何安装了Docker的机器上运行，而不用关心底层操作系统。</font>（这意思就是它来做好跨平台咯，就像装了一个jre一样）

Docker是一个开源的应用容器引擎，允许开发者将他们的应用以及依赖包打包到一个可移植的容器中，然后发布到任何装有Docker的linux机器上。

<font color="red">Docker仓库用来保存我们的images，当我们创建了自己的image之后我们就可以使用push命令将它上传到公有或者私有仓库，这样下次要在另外一台机器上使用这个image时候，只需要从仓库上pull下来就可以了。</font>

Docker镜像中包含了运行环境和配置，所以Docker可以简化部署多种应用实例工作。比如Web应用、后台应用、数据库应用、大数据应用比如Hadoop集群、消息队列等等都可以打包成一个镜像部署。

Docker简单来说就是用于容器化应用，简化部署工作。由于一个Docker容器已经包含项目运行所需要的所有依赖，只要它能在你的笔记本上运行，就能在任何一个别的远程服务器的生产环境上运行，包括Amazon的EC2和DigitalOcean上的VPS。

##Docker的安装##

Docker有很多种安装的选择，我们推荐您在Ubuntu下面安装，<font color="red">因为docker是在Ubuntu下面开发的</font>，安装包测试比较充分，可以保证软件包的可用性。Mac, windows和其他的一些linux发行版本无法原生运行Docker，<font color="red">可以使用虚拟软件创建一个ubuntu的虚拟机并在里面运行docker</font>。


##Docker 组件与元素##

###三个组件###

* Docker Client 是用户界面，它支持用户与Docker Daemon之间通信。
* Docker Daemon 运行于主机上，处理服务请求。
* Docker Index 是中央registry,支持拥有公有和私有访问权限的Docker容器镜像的备份。

###三个基本要求###
* Docker Containers 负责应用程序的运行，包括操作系统，用户添加的文件以及元数据。
* Docker Images 是一个只读模板，用来运行Docker容器
* DockerFile是文件指令集，用来说明如何自动创建Docker镜像。

关键点：<font color="red">运行任何应用程序，都需要有两个基本步骤：

* 构建一个镜像
* 运行容器

</font>

docker常用命令

<font color="red">Docker命令工具需要root权限才能工作。</font>

* sudo docker version #查看docker的版本
* sudo docker   #查看Docker的所有命令
* sudo docker command --help  #查看单个Docker命令的帮助，如Docker run --help
* docker info 
* docker pull
* docker run

在运行run命令时，你可指定链接、卷、端口、窗口名称（如果你没提供，Docker将分配一个默认名称）等等。

* docker logs
* docker help
* docker start
* docker stop
* docker restart
* docker rm      # 如果要完全移除容器，需要向将容器停止，然后才能删除
* docker commit  # 将容器的状态保存为镜像
* docker images  # 查看所有镜像的列表
* docker search  (image-name)  # 查找
* docker history  (image-name) # 查看历史版本
* docker push (image_name) #将镜像推送到registry

把learn/ping镜像发布到docker的index网站。

提示：

1. docker images命令可以列出所有安装过的镜像。
2. docker push命令可以将某一个镜像发布到官方网站。
3. 你只能将镜像发布到自己的空间下面


* --rm：告诉Docker一旦运行的进程退出就删除容器。
* -ti：告诉Docker分配一个伪终端并进入交互模式。







<font color="red">注意点：</font>

* 镜像名称只能取字符[a-z]和数字[0-9]。
* Docker 对系统资源的利用率很高，一台主机上可以同时运行数千个 Docker 容器


##镜像（image）##
docker image是一个构建容器的只读模板，它包含了容器启动所需的所有信息，包括运行程序和配置数据。





##容器##

容器是设计来运行一个应用的，而非一台机器。Docker提供了用于分离应用与数据的工具，使得你可以快捷地更新运行中的代码/系统，而不影响数据。

其实容器会一直存在，除非你删除它们,容器是持久的，直到你删除他们，并且你只能这样做：

```
  docker rm my_container #my_container是容器的名字或Id
```

如果您没有执行此命令，那么你的容器会一直存在，依旧可以启动、停止等。


##数据卷##
数据卷让你可以不受容器生命周期影响进行数据持久化，他们表现为容器内的空间，但实际保存在容器之外，从而允许你再不影响数据的情况下销毁、重建、修改、丢弃容器。Docker允许你定义应用部分和数据部分，并提供工具让你可以将它们分开。使用Docker时必须做出的最大思维变化之一就是:

<font color="red">容器应该是短暂和一次性的。</font>

数据卷是针对容器的，你可以使用同一个镜像创建多个容器并定义不同的卷。<font color="red">卷保存在运行Docker的宿主文件系统上，你可以指定卷存放的目录，或让Docker保存在默认位置。</font>
保存在其他类型文件系统上的都不是一个卷。

<font color="red">卷还可以用来在容器间共享数据。</font>

为了能够保存数据以及共享容器间的数据，docker提出了volume的概念，简单来说，volume就是目录或者文件，它可以绕过默认的联合文件系统，而以正常的文件或目录的形式存在于宿主机上。

<font color="red">Volume并不是为了持久化。</font>

Volume可以将容器以及容器产生的数据分离开来，这样，当你使用docker rm my_container删除容器时，不会影响相关的数据


docker inspect命令找到Volume在主机上的存储位置


可以通过两种方式来初始化volume

* 在Dockerfile中指定VOLUME /some/dir

* 在运行时使用-v来声明volume

```
docker run -v /home/adrian/data:/data debian ls /data

```

将挂载<font color="red">宿主机的/home/adrian/data目录</font>到<font color="red">容器内的/data目录</font>上


无论哪种方式都是做了同样的事情。它们告诉Docker在主机上创建一个目录（默认情况下是在/var/lib/docker下），然后将其挂载到指定的路径（例子中是：/some/dir）。当删除使用该Volume的容器时，Volume本身不会受到影响，它可以一直存在下去。
如果在容器中不存在指定的路径，那么该目录将会被自动创建。
你可以告诉Docker同时删除容器和其Volume：

```
docker rm -v my_container
```

##链接##

容器启动时，将被分配一个随机的私有IP,其它容器可以使用这个IP地址与其进行通讯。这点非常重要，原因有二：

* 它提供了一个容器间相互通信的渠道
* 容器将共享一个本地网络。

要开启容器间通讯，Docker允许你再创建一个新容器时引用其它现存容器，在你刚创建的容器里被引用的容器将获得一个（你指定的）别名，我们就说，这两个容器链接在了一起。

因此，如果db容器已经在运行，我可以创建web服务器容器，并在创建时引用这个DB容器，给它一个别名，比如dbapp。在这个新建的web服务器容器里，我可以在任何时候使用主机名dbapp与DB容器进行通讯。

##网络工具##

Docker的本地网络能力为容器间的连接提供两种方案。

* 公开一个容器的端口，并可选择性的映射到宿主机上并为外部路由服务。可以自己决定使用宿主机的端口来映射，也可以让Docker随机的选择一个未使用的高位端口。这是一种在大多数场景中用来提供对容器访问的友好方式。

* 采用Docker的'links'来允许容器间通信。一个关联的容器将会获得它的对应连接信息，在它处理了那些变量后允许它自动连接。这样就使得同一个宿主机上的容器不需要知道对应服务的端口和地址，就可以直接进行通信。


##容器话应用的架构##

独立容器的设计的一些特点：

* 他们不应该依赖或者关心宿主机上的任何细节
* 每一个组件应该提供一致性的接口，使得调用者可以访问服务
* 每一个服务应该在初始化阶段从环境变量中获取参数
* 应用产生的数据应该通过Valumes存储在<font color="red">容器外部</font>或者<font color="red">数据容器</font>中


##网络与通讯##
当Docker进程启动之后，它会配置一个虚拟的网桥叫docker0在宿主机上。这个接口允许Docker去分配虚拟的子网给即将启动的容器们。这个网桥在容器内的网络和宿主机网络之间将作为接口的主节点。

Docker容器启动后，将创建一个新的虚拟接口并分配一个网桥子网内的IP地址。这个IP地址嵌在容器内网络中，用于提供容器网络到宿主机docker0网桥上的一个通道。Docker自动配置iptables规则来放行并配置NAT，连通宿主机上的docker0。



##服务发现是怎么工作呢？##

每一个服务发现工具都会提供一套API，使得组件可以用其来设置或搜索数据。正是如此，对于每一个组件，服务发现的地址要么强制编码到程序或容器内部，要么在运行时以参数形式提供。通常来说，发现服务用键值对形式实现，采用标准http协议交互。

服务发现门户的工作方式是：当每一个服务启动上线之后，他们通过发现工具来注册自身信息。它记录了一个相关组件若想使用某服务时的全部必要信息。例如，一个MySQL数据库服务会在这注册它运行的ip和端口，如有必要，登录时的用户名和密码也会留下。

当一个服务的消费者上线时，它能够在预设的终端查询该服务的相关信息。然后它就可以基于查到的信息与其需要的组件进行交互。负载均衡就是一个很好的例子，它可以通过查询服务发现得到各个后端节点承受的流量数，然后根据这个信息来调整配置。


##配置存储是如何关联起来的？##
全局分布式服务发现系统的一个主要优势是它可以存储任何类型的组件运行时所需的配置信息。这就意味着可以从容器内将更多的配置信息抽取出去，并放入更大的运行环境中。

通常来说，为了让这个过程更有效率，应用在设计时应该赋上合理的默认值，并且在运行时可以通过查询配置存储来覆盖这些值。这使得运用配置存储跟在执行命令行标记时的工作方式类似。区别在于，通过一个全局配置存储，可以不做额外工作就能够对所有组件的实例进行同样的配置操作。




##docker##

使用docker每次都要用sudo，为了让非root用户使用docker，可将当前用户添加到docker用户组：

```
sudo groupadd docker
sudo gpasswd -a ${USER} docker  # 当前用户添加到docker group

```

 





