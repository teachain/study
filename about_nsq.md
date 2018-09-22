##NSQ是由四个重要组件构成##

* nsqd：一个负责接收、排队、转发消息到客户端的守护进程
* nsqlookupd：管理拓扑信息并提供最终一致性的发现服务的守护进程
* nsqadmin：一套Web用户界面，可实时查看集群的统计数据和执行各种各样的管理任务
* utilities：常见基础功能、数据流处理工具，如nsq_stat、nsq_tail、nsq_to_file、nsq_to_http、nsq_to_nsq、to_nsq

这些组件中，nsqd是可以独立部署的，也就是说可以没有其余三部分

关于安装，方式有二（其实没有安装过程）

* 使用源码编译安装（其实就是把go源码编译可执行文件）（https://github.com/nsqio/nsq）
* 直接下载预编译的可执行文件(http://nsq.io/deployment/installing.html)

因为是可执行文件，所以，只要能够找到他们就可以执行他们

要执行以下命令，一般有两种方式:

 * 第一种方式就是将这些可执行文件所有的目录加入到path环境变量中，
 * 第二种方式是cd到这些可执行文件所在目录

在terminal（终端）下，使用一下命令

```
 $ nsqlookupd 
```
使用该命令将启动nsqlookupd监控服务，在该命令后面加一个&即可以守护进程的方式启动服务。这个服务将会监听两个默认的端口4160和4161,在4160上是tcp服务，在4161上开启的是http服务。

```
$ nsqd --lookupd-tcp-address=127.0.0.1:4160
```

启动nsqd服务，并通知nsqlookupd,意思是说nsqd跟nsqlookupd说，我启动了,默认nsqd在4151和4150上监听,在4150上是tcp服务在4151上开始的是http服务。当你需要开启https时，可以-https-address来指定地址和端口，详细可参考下文中的启动命令行参数

```
$ nsqadmin --lookupd-http-address=127.0.0.1:4161

```
启动nsqadmin服务，并通知nsqlookupd,意思是说nsqadmin跟nsqlookupd说，我启动了，默认nsqadmin在4171上监听

###nsqd启动时可用的命令行参数###

```
-auth-http-address=: <addr>:<port> to query auth server (may be given multiple times)
-broadcast-address="": address that will be registered with lookupd (defaults to the OS hostname)
-config="": path to config file
-data-path="": path to store disk-backed messages
-deflate=true: enable deflate feature negotiation (client compression)
-e2e-processing-latency-percentile=: message processing time percentiles (as float (0, 1.0]) to track (can be specified multiple times or comma separated '1.0,0.99,0.95', default none)
-e2e-processing-latency-window-time=10m0s: calculate end to end latency quantiles for this duration of time (ie: 60s would only show quantile calculations from the past 60 seconds)
-http-address="0.0.0.0:4151": <addr>:<port> to listen on for HTTP clients
-https-address="": <addr>:<port> to listen on for HTTPS clients
-lookupd-tcp-address=: lookupd TCP address (may be given multiple times)
-max-body-size=5123840: maximum size of a single command body
-max-bytes-per-file=104857600: number of bytes per diskqueue file before rolling
-max-deflate-level=6: max deflate compression level a client can negotiate (> values == > nsqd CPU usage)
-max-heartbeat-interval=1m0s: maximum client configurable duration of time between client heartbeats
-max-msg-size=1024768: maximum size of a single message in bytes
-max-msg-timeout=15m0s: maximum duration before a message will timeout
-max-output-buffer-size=65536: maximum client configurable size (in bytes) for a client output buffer
-max-output-buffer-timeout=1s: maximum client configurable duration of time between flushing to a client
-max-rdy-count=2500: maximum RDY count for a client
-max-req-timeout=1h0m0s: maximum requeuing timeout for a message
-mem-queue-size=10000: number of messages to keep in memory (per topic/channel)
-msg-timeout="60s": duration to wait before auto-requeing a message
-snappy=true: enable snappy feature negotiation (client compression)
-statsd-address="": UDP <addr>:<port> of a statsd daemon for pushing stats
-statsd-interval="60s": duration between pushing to statsd
-statsd-mem-stats=true: toggle sending memory and GC stats to statsd
-statsd-prefix="nsq.%s": prefix used for keys sent to statsd (%s for host replacement)
-sync-every=2500: number of messages per diskqueue fsync
-sync-timeout=2s: duration of time per diskqueue fsync
-tcp-address="0.0.0.0:4150": <addr>:<port> to listen on for TCP clients
-tls-cert="": path to certificate file
-tls-client-auth-policy="": client certificate auth policy ('require' or 'require-verify')
-tls-key="": path to private key file
-tls-required=false: require TLS for client connections
-tls-root-ca-file="": path to private certificate authority pem
-verbose=false: enable verbose logging
-version=false: print version string
-worker-id=0: unique identifier (int) for this worker (will default to a hash of hostname)

```

###nsqlookupd启动时可用的命令行参数###
```
-http-address="0.0.0.0:4161": <addr>:<port> to listen on for HTTP clients
-inactive-producer-timeout=5m0s: duration of time a producer will remain in the active list since its last ping
-tcp-address="0.0.0.0:4160": <addr>:<port> to listen on for TCP clients
-broadcast-address: external address of this lookupd node, (default to the OS hostname)
-tombstone-lifetime=45s: duration of time a producer will remain tombstoned if registration remains
-verbose=false: enable verbose logging
-version=false: print version string

```

###nsqadmin启动时可用的命令行参数###

```
-graphite-url="": URL to graphite HTTP address
-http-address="0.0.0.0:4171": <addr>:<port> to listen on for HTTP clients
-lookupd-http-address=[]: lookupd HTTP address (may be given multiple times)
-notification-http-endpoint="": HTTP endpoint (fully qualified) to which POST notifications of admin actions will be sent
-nsqd-http-address=[]: nsqd HTTP address (may be given multiple times)
-proxy-graphite=false: Proxy HTTP requests to graphite
-statsd-interval=1m0s: time interval nsqd is configured to push to statsd (must match nsqd)
-statsd-prefix="nsq.%s": prefix used for keys sent to statsd (%s for host replacement, must
-template-dir="": path to templates directory
-use-statsd-prefixes=true: expect statsd prefixed keys in graphite (ie: 'stats_counts.')
-version=false: print version string

```



目前从测试出来的数据来看，一个topic会同时将数据同步发布给各个channel,一个channel下可以有多个consumer,对一个channel而言，一条消息只会交给一个consumer来处理


