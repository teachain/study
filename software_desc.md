##Supervisor##

 Supervisor是一个客户/服务器系统，它可以在类Unix系统中管理控制大量进程。Supervisor使用python开发，有多年历史，目前很多生产环境下的服务器都在使用Supervisor。
 
supervisord会帮你把管理的应用程序转成daemon程序，而且可以方便的通过命令开启、关闭、重启等操作，而且它管理的进程一旦崩溃会自动重启，这样就可以保证程序执行中断后的情况下有自我修复的功能。

###服务器###
Supervisor的服务器端称为supervisord，主要负责在启动自身时启动管理的子进程，响应客户端的命令，重启崩溃或退出的子进程，记录子进程stdout和stderr输出，生成和处理子进程生命周期中的事件。可以在一个配置文件中配置相关参数，包括Supervisord自身的状态，其管理的各个子进程的相关属性。配置文件一般位于/etc/supervisord.conf。


###客户端###
Supervisor的客户端称为supervisorctl，它提供了一个类shell的接口（即命令行）来使用supervisord服务端提供的功能。通过supervisorctl，用户可以连接到supervisord服务器进程，获得服务器进程控制的子进程的状态，启动和停止子进程，获得正在运行的进程列表。客户端通过Unix域套接字或者TCP套接字与服务端进行通信，服务器端具有身份凭证认证机制，可以有效提升安全性。当客户端和服务器位于同一台机器上时，客户端与服务器共用同一个配置文件/etc/supervisord.conf，通过不同标签来区分两者的配置。

平台要求

```
Supervisor可以运行在大多数Unix系统上，但不支持在Windows系统上运行。

Supervisor需要Python2.4及以上版本，但任何Python 3版本都不支持。

```


配置文件模版在/etc/supervisord.conf中生成。当然配置文件的位置不一定放在/etc/supervisord.conf中。在不指定配置文件位置的情况下，supervisord和supervisorctl会按照以下三个顺序来搜索配置文件的位置。

* $CWD/supervisord.conf
* $CWD/etc/supervisord.conf
* /etc/supervisord.conf


进程配置在supervisord.conf

然后使用supervisorctl来管理

