##Supervisor在centos7上的安装

在terminal里切换至root用户，然后输入以下命令即可安装

```
easy_install supervisor

```

接着运行以下命令

```
echo_supervisord_conf

```

运行以下命令生成配置文件

```
echo_supervisord_conf > /etc/supervisord.conf

```

##配置文件##

```
/etc/supervisord.conf 

```

当我们执行以下命令

```
supervisorctl staus

```
出现以下错误时

```
http://127.0.0.1:9001 refused connection

```

我们需要将文件/etc/supervisord.conf 中的以下内容将前面的分号去掉，然后

```
[inet_http_server]         ; inet (TCP) server disabled by default
port=127.0.0.1:9001        ; ip_address:port specifier, *:port for all iface

```
将supervisord进程杀掉，然后重新执行以下命令

```
supervisord -c /etc/supervisord.conf

```

web.conf

```
[program:webapp]
directory=/home/david/app/web
command=/home/david/app/web/main
autostart=true
startsecs=5
autorestart=true
startretries=3
user=david
redirect_stderr=true
stdout_logfile = /home/david/app/web/web_stdout.log

```
