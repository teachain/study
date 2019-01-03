# 浅析GPU计算——cuda编程

```
Thread是CUDA最基本的执行单元，多个threads组成一个block，多个blocks组成grid。

现在的GPU一个block最多可以分配1024个threads。

Block中的threads可以是一维、二维、三维的，grid中的blocks也可以是一维、二维、三维的。

Block可以通过内建变量blockIdx索引，block的维度可以用dim3类型的内建变量blockDim访问。

一维：
threadIdx.x: block中的thread索引 

blockIdx.x :Grid中block索引  

blockDim.x：一个Block中thread数量; 

gridDim.x:一个Grid中Block数量  

使用线程与使用线程块相比有什么优势？ ：

一：解决线程块数量的限制 

二：进行部分的数据共享和通讯

从以上的内容我们可以得出，因为Thread是CUDA最基本的执行单元（不可分割，可认为是原子的），

所以blocknum就相当于grid,threadNum就相当于block

所以调用的地方<<<blocknum,threadNum>>> 就是<<<grid,block>>>

又因为他们是三维的，所以(x,y,z)是必然的表示方法
```



```
我们将CPU以及系统的内存称为主机，而将GPU及其内存称为设备。
```

```
在GPU设备上执行的函数通常称为核函数(Kernel)。
```

```
使用__global__来修饰函数，表示该代码执行在设备上

调用时，函数名称<<<>>>(参数...) 使用这样形式来调用。

尖括号表示要将一个参数传递给运行时系统，这些参数并不是传递给设备代码的参数，而是告诉运行时如何启动设备代码。

传递给设备代码本身的参数是放在圆括号中传递的，就像标准的函数调用一样。

```

```
int *dev_c;

int c

 //在设备上分配内存，这个函数的作用是告诉CUDA运行时在设备上分配内存
cudaMalloc((void**)&dev_c,sizeof(int));

//将设备上的内存上的变量值拷贝到主机内存上
//cudaMemcpyDeviceToHost告诉运行时源指针是一个设备指针，而目标指针是一个主机指针
//dest=src,记得是这种形式，自然就记得参数的位置该怎么放
//cudaMemcpyHostToDevice  //主机到设备
//cudaMemcpyDeviceToDevice //都在设备上
//这种函数都是在主机上去操作的时候用的。
//当源指针和目标指针都位于主机上时，自然就是标准的C memcpy()函数。
cudaMemcpy(&c,dev_c,sizeof(int),cudaMemcpyDeviceToHost)

//释放设备上的内存
cudaFree(dev_c); 

```

```
int count;
cudaGetDeviceCount(&count);//获取设备的数量
cudaDeviceProp prop;
for(int i=0;i<count;i++){
    cudaGetDeviceProperties(&prop,i);//获取单个设备的属性，一个主机上可能装N个显卡
}

```

```
<<<n,m>>> n表示设备在执行核函数时使用的并行线程块的数量,m表示CUDA运行时在每个线程块中创建的线程数量。

1<=n<=65535,m不能超过设备属性结构中maxThreadPerBlock域的值。

在核函数中，通过blockIdx.x可以得到线程块的索引,blockIdx是一个内置变量(在cuda中天生就有的，拿来用即可)，

其实是因为它是在CUDA运行时预先定义的。通过threadIdx.x可以得到并行线程的索引，threadIdx是一个内置变量。

blockDim,对于所有线程块来说，这个变量是一个常数，保存的是线程块中每一维的线程数量。

int tid=threadIdx.x+blockIdx.x*blockDim.x


```

| 线程0 | 线程2 | 线程2 | 线程3 |
| ----- | ----- | ----- | ----- |
| 线程0 | 线程2 | 线程2 | 线程3 |
| 线程0 | 线程2 | 线程2 | 线程3 |
| 线程0 | 线程2 | 线程2 | 线程3 |

如果线程表示列，而线程块表示行，那么可以计算得到一个唯一的索引，将线程块索引与每个线程块中的线程数量相乘，然后加上线程在线程块中的索引。也就是

int tid=threadIdx.x+blockIdx.x*blockDim.x

```
共享内存 __share__
常量内存__constant__
```



CPU是整个计算机的核心，它的主要工作是负责调度各种资源，包括其自身的计算资源以及GPU的计算计算资源。比如一个浮点数相乘逻辑，理论上我们可以让其在CPU上执行，也可以在GPU上执行。那这段逻辑到底是在哪个器件上执行的呢？cuda将决定权交给了程序员，我们可以在函数前增加修饰词来指定。

```
关键字              执行位置
__host__           CPU
__global__         GPU
__device__         GPU
```

一般来说，我们只需要2个修饰词就够了，但是cuda却提供了3个——2个执行位置为GPU。这儿要引入一个“调用位置”的概念。父函数调用子函数时，父函数可能运行于CPU或者GPU，相应的子函数也可能运行于CPU或者GPU，但是这绝不是一个2*2的组合关系。因为GPU作为CPU的计算组件，不可以调度CPU去做事，所以不存在父函数运行于GPU，而子函数运行于CPU的情况。

```
关键字             调用位置
__host__           CPU
__global__         CPU
__device__         GPU
```

```
__global__描述的函数就是“被CPU调用，在GPU上运行的代码”，同时它也打通了host和device修饰的函数。
 如果一段代码既需要运行于CPU，也要运行于GPU，怎么办？难道要写两次？当然不用，我们可以同时使用__host__和__device__修饰。这样编译器就会帮我们生成两份代码逻辑。
 
 
特别需要注意： __global__不能和它们（__host__和__device__）中任何一个修饰符一起使用。

 __global__这个修饰符的使命使得它足够的特殊。比如它修饰的函数是异步执行的。 __global__修饰的函数只能是void类型。

```

```
   cuda是一个GPU编程环境，所以它对__device__修饰的函数进行了比较多的优化。比如它会根据它的规则，让某个__device__修饰函数成为内联函数（inline）。这些规则是程序员不可控，但是如果我们的确对是否内联有需求，cuda也提供了方式：使用__noinline__修饰函数不进行内联优化；使用 __forceinline__修饰函数强制进行内联优化。当然这两种修饰符不能同时使用。

   也许你已经发现，__global__函数调用方式非常特别——使用“<<<>>>”标志。这儿就需要引入cuda的并行执行的线程模型来解释了。在同一时刻，一个cuda核只能运行一个线程，而线程作为逻辑的运行载体有其自己的ID。这个ID和我们在linux或windows系统上CPU相关的线程ID有着不同的表达方式。比如在Linux系统上可以使用gettid方法获取一个pid_t值，比如3075。但是cuda的表达方式是一个三维空间，表达这个空间的是一个叫block的概念。比如单个block定义其有(Dx, Dy, 0)个线程，则每个线程ID为x+yDx；再比如有(Dx, Dy, Dz)个线程，则每个线程ID为x+yDx+zDxDy。

```

