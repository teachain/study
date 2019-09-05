##Redis简介##
REmote DIctionary Server(Redis) 是一个由Salvatore Sanfilippo写的<font color="red">key-value</font>存储系统。
Redis是一个开源的使用ANSI C语言编写、遵守BSD协议、支持网络、可基于内存亦可持久化的日志型、Key-Value数据库，并提供多种语言的API。
它通常被称为数据结构服务器，<font color="red">因为值（value）可以是 字符串(String), 哈希(Map), 列表(list), 集合(sets) 和 有序集合(sorted sets)等类型。</font>


Redis 与其他 key - value 缓存产品有以下三个特点：

* Redis支持数据的持久化，可以将内存中的数据保持在磁盘中，重启的时候可以再次加载进行使用。
* Redis不仅仅支持简单的key-value类型的数据，同时还提供list，set，zset，hash等数据结构的存储。
* Redis支持数据的备份，即master-slave模式的数据备份。



Redis提供了两种不同的持久化方法来将数据存储到硬盘里面：

* 快照(snapshotting),它可以将存在于某一时刻的所有数据都写入硬盘里面。
* 只追加文件(append-only file,aof),它会在执行写命令时，将被执行的写命令复制到硬盘里面。

这两种持久化方法既可以同时使用，又可以单独使用。在某些情况下甚至可以两种方法都不使用。



#### 快照

```
save 60 10000
stop-writes-on-bgsave-error no
rdbcompression yes
dbfilename dump.rdb
dir ./ #保存dump.rdb的位置
```

缺点：redis、系统或者硬件这三者之中的任意一个崩溃了，那么redis将丢失最近一次创建快照之后写入的所有数据。

创建快照的方法

* 客户端通过向redis发送bgsave命令来创建一个快照。redis会调用fork来创建一个子进程，然后子进程负责将快照写入硬盘，而父进程则继续处理命令请求。
* 客户端通过向redis发送save命令来创建一个快照，接到save命令的redis服务器在快照创建完毕之前将不再响应任何其他命令。
* 设置save配置项，比如 save 60 10000,那么从redis最近一次创建快照之后开始算起，当"60秒之内有10000次写入"这个条件被满足时，redis就会自动触发bgsave命令，如果用户设置了多个save配置选项，那么当任意一个save配置选项所设置的条件被满足时，redis就会触发一次bgsave命令。
* 收到shutdown命令，会还行一个save命令，阻塞所有客户端，不再执行客户端发送的任何命令，并在save命令执行完毕之后关闭服务器。
* 当一个redis服务器连接另一个redis服务器，并向对方发送sync命令来开始一次复制操作的时候，如果主服务器目前没有在执行bgsave操作，或者主服务器并非刚刚执行完bgsave操作，那么主服务器就会执行bgsave命令。

总结下来就是，执行bgsave操作或执行save操作。要特别注意save的条件，多长时间多少次***写入***

#### AOF

```
appendonly yes
appendfsync everysec #可用的选项：always、everysec、no
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
dir ./ #保存追加文件的位置
```

简单来说，AOF持久化会将被执行的***写命令***写到AOF文件的末尾，以此来记录数据发生的变化。因此，redis只要从头到尾重新执行一次aof文件包含的所有写命令，就可以恢复AOF文件所记录的数据集。

appendfsync 可用的选项：

* always ,每个redis写命令都要同步写入硬盘。这样做会严重降低redis的速度（不推荐）

* everysec，每秒执行一次同步，显式地将多个写命令同步到硬盘。（推荐使用）
* no，让操作系统来决定应该何时进行同步。（不推荐）

缺点:aof文件的大小，因为redis在重启之后需要通过重新执行AOF文件记录的所有写命令来还原数据集，所以如果AOF文件的体积非常大，那么还原操作执行的时间就可能会非常长。

为了解决AOF文件体积不断增大的问题，用户可以向redis发送bgrewriteaof命令，这个命令会通过移除aof文件中的冗余命令来重写aof文件，使aof文件的体积变得尽可能地小。

跟快照持久化可以通过设置save选项来自动执行bgsave一样，aof持久化也可以通过设置auto-aof-rewrite-percentage选项和auto-aof-rewrite-min-size选项来自动执行bgrewriteaof。

#### 复制

复制可以让其他服务器拥有一个不断地更新的数据副本，从而使得拥有数据副本的服务器可以用于处理客户端发送的读请求。关系数据库通常会使用一个主服务器向多个从服务器发送更新，并使用从服务器来处理所有读请求。redis也采用了同样的方法来实现自己的复制特性，并将其用作扩展性能的一种手段。



当从服务器连接主服务器的时候，主服务器会执行bgsave操作。

***开启从服务器所必须的选项只有saveof一个***

如果用户在启动redis服务器的时候，指定了一个包含slaveof host port选项的配置文件，那么redis服务器将根据该选项给定的ip地址和端口号来连接主服务器。对于一个正在运行的redis服务器，用户可以通过发送slaveof no one 命令来让服务器终止复制操作，不在接受主服务器的数据更新；也可以通过发送slaveof host port命令来让服务器开始复制一个新的主服务器。

从服务器在进行同步时，会清空自己的所有数据。redis不支持主主复制。但从服务器可以有从服务器，这样就形成了主从链也就是 主-->从-->从，这样子的链。也就是如果主因为多个从需要同步而导致系统超载（主-->从），我们可以通过这样方式来扩展从的数量而不会造成主的系统超载（我们使用主-->从-->从这种方式来解决,当然可以有n层从，以树的结构来看）。

***检查Info命令的输出结果中aof_pending_biio_fsync属性的值是否为0，如果是的话，那么就表示服务器已经将一直的所有数据都保存到硬盘里面了。***



***通过同时使用复制和AOF持久化，用户可以增强redis对于系统崩溃的抵抗能力。***

更换故障主服务器的方法

方法一：首先想从服务器发送一个save命令，让它创建一个新的快照文件，接着将快照文件发送给新的主服务器，并在新的主服务器上面启动redis,最后，让从服务器成为主服务器的从服务器（也即是执行slaveof host port命令）。

方法二：将从服务器升级为主服务器，并为升级后的主服务器创建从服务器。

用户接下来要做的就是更新客户端的配置，让它们去读正确的服务器。

### redis事务

redis的事务以特殊命令multi开始，之后跟着用户传入的多个命令，最后以exec为结束。但是由于这种简单的事务在exec命令被调用之前不会执行任何实际操作，所以用户将没办法根据读取到的数据来决定。

page108











