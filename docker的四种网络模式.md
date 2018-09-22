1. host模式

   ```
   docker run 使用--net=host指定，使用的网络实际上和宿主机一样
   ```

2. container模式

   ```
   使用--net=container:container_id/container_name
   多个容器使用共同的网络，看到的ip是一样的
   ```

3. none模式

   ```
   使用--net=none指定
   这种模式下，不会配置任何网络
   ```

4. bridge模式

   ```
   使用--net=bridge指定
   默认模式，不需要指定
   此模式会为每一个容器分配一个独立的network namespace
   ```

   ​

   ​

   ​