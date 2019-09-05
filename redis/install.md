#### linux下redis的安装

安装必要的编译工具

```
yum -y install gcc  
yum -y install gcc-c++  
```

获取文件并编译安装

```
wget http://download.redis.io/releases/redis-5.0.2.tar.gz
tar xzf redis-5.0.2.tar.gz
cd redis-5.0.2
make & make install #如果不想install的话，make就好了，可执行文件会在src下。
```

