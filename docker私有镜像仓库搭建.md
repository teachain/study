1. 首先是在机器上先安装好docker。

2. 下载registry镜像，也就是执行以下命令：

   ```
   docker pull registry
   ```

3. 启动一个容器

   ```
   docker run -d -p 5000:5000 registry 
   ```

   ​

4. 测试

   ```
   sudo docker pull busybox  
   sudo docker tag busybox 192.168.0.102:5000/busybox  
   sudo docker tes、 192.168.0.102:5000/busybox 
   vim /usr/lib/systemd/system/docker.service
   --exec-opt native.cgroupdriver=systemd \  
   --insecure-registry=192.168.0.102:5000 \ 
   systemctl restart docker  
   ```

   ​

5. ​

  