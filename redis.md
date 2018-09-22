##redis##

redis的数据类型操作命令肯定是这样子的

```
  command key ...  
```


###1.String 存入字符类型###

 * Set name luowen  #设置name = luowen 存储
 
 * Get name    #获取设置好的name的值
 
 * Setnx name luowen 设置name键值为luowen #如果存在,则返回0 不存在返回1
 
 * Mset name luowen age 23 salary 233333 #设置多个键值对 一块存错 全成功,全失败
 
 * Msetnx name maomao age 23 hoby basketball #如果设置多个键值对中有存在返回失败
 
 * Mget name age salary #获取多个键的值
 
 * Getset name maomao #获取name的旧值,并设置新的值为maomao

 * Setrange name 3 maomao #将键name的值从下标3开始进行替换，假设原来name的值是luowen,那么执行该命令之后，name的值就修改为luomaomao了。
 
 * Getrange name 3 6 #获取键name的值,取一个范围下标3~6（3<=index<=6） 

 * Append name .com #给键nane追加.com 

 * Incr age #设置每个值自增（加1）
 
 * Incrby age 6 给name加上6 如果是负数则减

 * Decr 与incr相反
 
 * Decrby 与Incrby相反

 * Strlen 返回键对应的值得字符长度


###2.Hash 方便存对象 键值对 ###

* Hset user:001 name luowen    #设置哈表key为user:001 的name 设置为 luowen

* Hsetnx user:001 name maomao    #设置哈希表名字中的name属性 ,若name属性已存在,设置不成功,若name属性不存在，则设置成功。

* Hget user:001 name #获取user:001的name属性

* Hmset user:001 name maomao age 23 #批量设置（设置user:001的多个属性）

* Hmget user:001 name age #批量获取user:001的属性（获取user:001的多个属性）

* Hincrby user:001 age 3  #给user:001的属性age值加上3,by的值可正可负

* Hexists user:001 name #判断user:001是否存在name属性

* Hlen user:001 #获取user:001属性个数

* Hkeys user:001 #返回user:001的所有字段

* Hvals user:003 #返回user:001所有的属性的值

* Hgetall user:001 #返回user:001所有的属性和属性值

* Hdel user:001 name #删除user:001的name的属性和属性键

##3.list 链表 (双向链表)##

* 栈:先进后出（（lpush和lpop）组合实现或（rpush和rpop）组合实现）

* 队列:先进先出（(rpush和lpop)组合实现或(lpush和rpop)组合实现）

将list链表想象成一根管道，将元素想象成小球，水平放置

以下的<font color="red">管道是一个形象的说法，其实是指list，弹出必然删除元素，压入必然增加元素,然后下标必须是从左开始计算，从0开始，list表中的元素是可重复存在的。</font>

1、lpush 从左边压入

```
Lpush list1 world 
lpush list1 hello
Lrange list1 0 -1 #把链表中的数据从0到尾全部取出(注意没有rrange命令)
 #得到的结果是
 1) "hello"
 2) "world"
``` 

2、rpush 从右边压入
 
 ```
 rpush list2  world
 rpush list2  hello
 lrange list2 0 -1
  #结果是
  1) "world"
  2) "hello"
 ```
 3、linsert 插入数据
 
 ```
Rpush list3 luowen
Rpush list3 maomao
Lrange list3 0 -1
#结果是
	1) "luowen"
	2) "maomao"
    
Linsert list3 before maomao love
Lrange list3 0 -1
#结果是
	1) "luowen"
	2) "love"
	3) "maomao"

	
linsert list3 after maomao forever
Lrange list3 0 -1
 #结果是
	1) "luowen"
	2) "love"
	3) "maomao"
	4) "forever"
 ```
 4 lset 给某个元素（针对下标）赋值
 
 ```
 Rpush list5 luowen
 Rpush list5 maomao
 Lset list5 0 damin #将下标为0的元素赋值为damin
 ```
5、lrem 删除list表中的数据（删除给定的某个值，指定个数）

```
Rpush list6 luowen
Rpush list6 luowen1
Rpush list6 luowen2
Rpush list6 luowen3
Rpush list6 luowen4
Lrem list6 1 “luowen”# 删除list6 中值为luowen的值 1表示删除的个数

```

 6、ltrim(保留指定下标【一个范围】的元素) ，其余的将从管道中删除。
 
```
Lpush list7 luowen1
Lpush list7 luowen2
Lpush list7 luowen3
Lpush list7 luowen4
Lpush list7 luowen5
Ltrim list7 1 2 (1 2 为保留的范围)

```
7 、lpop 从链表的左边弹出一个元素

```
Lpush list8 luowen1
Lpush list8 luowen2
Lpush list8 luowen3
Lpop list8  #从左边弹出一个元素返回给client,并将元素从管道中删掉。

```

8、rpop 从链表的右边弹出一个元素

```
Lpush list8 luowen1
Lpush list8 luowen2
Lpush list8 luowen3
rpop list8 #从右边弹出一个元素返回给client,并将元素从管道中删掉。
```
9、rpoplpush 从一个链表右边弹出,在从左边压入到另一个链表
 
 ```      
 Rpoplpush list1 list2 #从list1右边弹出一个元素，然后从list2的左边压入这个元素
 ```
10、lindex 返回一个list下标的索引的元素

```
lindex list11 0 #返回下标为0的元素
```
        
11、llen 返回这个链表的元素的长度

```
llen list12 #返回链表的元素的长度

```


##4、set集合（无序集合）(元素唯一，不存在重复元素)##
1 sadd 向集合中插入一条数据

```
Sadd myset1 luowen #如果myset1中不存在luowen,那么插入成功，否则插入失败
```

2 srem 删除集合中的一个元素

```
Srem myset1 luowen
```
3 smembers 查看集合中的元素(获取集合中的元素)

```
Smembers myset1
```

4 spop 从集合随机弹出一个元素,返回元素（元素会从集合中删除）

```
spop myset1

```

5 sdiff 返回两个集合的差集 返回两个集合不一样的,<font color="red">根据第一个集合为标准</font>

```  
sdiff setdemo1 setdemo2

```

6 sdiffstroe 将两个差集存储到另外一个集合

<font color="red">务必记住目标集合是第一个参数</font>

Redis Sdiffstore 命令将给定集合之间的差集存储在指定的集合中。如果指定的集合 key 已存在，则会被覆盖。比如说假设setdemo3已经存在，那么setdemo1 setdemo2的diff结果会覆盖掉原来setdemo3的内容

```    
Sdiffstore setdemo3 setdemo1 setdemo2 #将setdemo1 setdemo2的查集保存到setdemo3

```
7 sinter 求两个集合的交集（取相同的元素）

```
sinter set1 set2
```

8 sinterstore 将两个集合的交集存储到另外一个集合中

```
sinterstore myset6 myset3 myset4#将myset3 myset4的交集保存到myset6中

```

9 sunion 	求两个集合并集

```
  sunion myset3 myset4
```
10 sunionstore 将两个集合并集并存储到另外一个集合中

```
  sunionstore myset7 myset3 myset4将myset3 myset4的并集保存到myset7中
 
```
11 smove 将以个集合中的元素移动到另外一个集合中

```
smove myset3 myset4 damin5 #将元素damin5从myset3中移动到myset4中 

```

12 scard 查看集合中元素的个数

```
 scard myset3 #返回myset3集合中元素的个数
```

13 sismember 判断是否是集合中的元素

```
Sismember myset3 luowen #判断luowen是否是myset3中的元素

```

14 srandmember 从集合中随机读取一个元素

```
srandmember myset3 #随机取出myset3 中的元素(不会删除元素)

```

##5、有序集合zset##

1、 zadd 添加到有序集合中(redis自动根据score来安排元素所在位置也就是index)

```
Zadd myzset 1 luowen1 #1 表示score
```

2 、Zrange 获取有序集合的元素

```
Zrange myzset 0 -1 withscores #获取有序集合的元素 0 -1表示下标，withscores表示连score一起返回

```

3、zrem 删除有序集合中的元素

```
zrem myzset luowen1 删除myzsent集合中的luowen1

```

4、zincrby(修改元素的score,数值为正则增加，数组为负则减少) 

```
zincrby myzset 3 luowen2 #如果luowen2没有,就创建他,如果有则增加3

```
5、zrank 返回元素索引

```
zrank myzset luowen2 #返回luowen2在myzset中的索引

```

6、zrevrank 返回倒序的索引

```
 zrevrank myzset luowen2
 
```
7、zrangebyscore key min max [withscores] [limit offset count]

默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)。

```
zrangebyscore myzset 1 5  #小于等于或大于等于

zrangebyscore myzset (1 (5  #小于5或大于1

```
8、zcount 返回指定score空间的数量

```
zcount myzset3 3 6  #返回score在3~6之间一共有多少个元素

```
9、 zcard 返回集合中所有元素的个数

```
zcard myzset3

```

10、zremrangbyrank key start stop 删除集合中指定区间的元素

```
zremrangebyrank myzset3 0 1 #将索引为0~1之间的元素删除

```
11、zremrangebyscore key min max 删除集合中指定元素

```
zremrangebyscore myzset3 3 6#将score为3~6之间的元素删除

```

5 Redis常用命令

    Key-values
    1 keys * 匹配键所有的键. 模糊匹配 keys my* 取出所有已my开头的键
    2 exists 判断是否键 exists name判断是否有name这个键是否存在
    3 del 删除键 del name 删除name的键
    4 expire 设置过期时间 expire key time 
    5 ttl key 查看键的过期时间
    6 select database 选择数据库
    7 move key dababase1 讲key移动dao database1中的数据库中
    8 persist 取消键的过期时间 
    9 randomkey 随机返回一个键的值
    10 rename 重命名一个键
    11 type key 判断key的数据类型

    Server
    1 ping ping我们的主机能否链接 链接是否存活
    2 echo 命令 echo demo直接输出
    3 select 选择数据库 select 0-16个数据库
    4 quit exit 退出链接
    5 dbsize 返回数据库的键的个数
    6 info 返回服务器相关信息
    7 config get 返回服务配置信息
    8 flush db 清空数据库
    9 flushall 删除所有数据库中所有的键

6 Redis 高级应用

    1 在配置文件里面设置 requirepass password
    2 进入后 auth 密码 进行授权 方法二: 登入或在后面加上 –a 加上密码
    3 主从复制:
        One: 一个master服务器可以拥有多个slave
        Two: 一个salve可以有多个master 并且还可以与其他的salve相连接
    配置salve
        打开salveof注释 并添加主机的ip以及端口
        主机加了密码的时候还需要配置masterauth 密码
    4 redis 的事务处理
        输入:multi 打开一个上下文
        Set age 10
        Set age 144
        -----------------------------------------------------------
        上面的全部放入队列最后执行
    Exec 
        最后age为144
        回滚
        Discard    
        
        Watch 监视键的命令
    5 Redis的持久化
        方式一:  snapshotting (快照)将内存的数据写入到文件中 save 500 32 500秒内有32个键发生变化则发起快照到文件中
        方式二: append only file 将没次写修改的命令保存到文件中
        配置:打开append only
                Appendfsync yes
                Appendfsync always 每次都写入
                Appendfsync everysec    每个一秒写入
                Appendfsync no    不写入
    6 发布和订阅消息
        订阅:
    　　 Subscribe tv1 tv2 订阅了两个频道
        发布:
        Publish tv1 luweo
        注:publish tv1的信息 订阅的信息都可以收到
    7 虚拟内存
        方式一:暂时把不使用的数据放到硬盘里面
    　  方式二:可以把数据分割到其他的slave数据服务器中