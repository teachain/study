OpenCL中的数据并行性表述为一个N维计算域，其中N=1、2或3（其实就是一维、二维、三维）。N-D域定义了可以并行执行的工作项的总数。

关键字

```
//等价
__kernel和kernel
__global和global
__constant和constant
__local和local
__private和private
```

get_global_id(0)返回各个工作项的1维全局ID(get_global_id(1)二维,get_global_id(2)三维)



如果内核函数的参数声明为某个类型的指针，那么这样的参数只能指向以下地址空间global、local或constant。

工作条目就是一个最小的执行单元，工作条目可以组成工作组（work group）。这样的划分也与存储器有关，在opencl中，存储分为三大类：Global memory，Local memory，以及Private memory。Global是可以让所有的工作组和工作条目都可见，Local是只有当前工作组中的工作条目可见，而Private是只有单独一个工作条目可见。这样的存储访问控制，可以有效利用高速缓存提高效率，而不是每一次数据访问都需要外部DDR。

（1）work_dims：the number of dimensions in the data ( if you deal with image object, you should probably set work_dims equal 2 or 3. But for buffer objects, you can set whatever dimensionality you think best. For a buffer object containing a two-dimensional matrix, you might set work-dims equal 2.)其实就是工作的维数（一维、二维或者三维）

（2）global_work_offset：the global ID offset in each dimension （每一维的 id的偏移量）

（3）global_work_size：the number of work items in each dimension  (the global_work_size argument of clEnqueueNDRangeKernel identifies how many work-items need to be processed for each dimension. )每一维的工作项的个数

（4）local_work_size：the number of work_items in a work_group，in each dimension  (local_work_size less than the global_work_size）在一个工作组里的工作项个数。



OpenCL C中使用类型修饰符__global(或global)描述的指针就是全局内存。全局内存的数据可以被执行内核中的所有工作项所访问到。不同类型的内存对象的作用域是不同的。

__global可以用来描述一个全局数组。数组中可以存放任意类型的数据：标量、矢量或自定义类型。无论数组中存放着什么类型的数据，其是通过指针顺序的方式进行访问，并且执行内核可以对数组进行读写，也可以设置成针对内核只读或只写的数组。数组数据和C语言的数组几乎没有任何区别，C语言怎么使用的，就可以在OpenCL内核代码中怎么使用





任何以`clEnqueue`开头的OpenCL API都能向命令队列提交一个命令，并且这些API都需要一个命令队列对象作为输入参数。例如，`clEnqueueReadBuffer()`将device上的数据传递到host，`clEnqueueNDRangeKernel()`申请一个内核在对应device执行。



OpenCL C上的并发执行单元称为工作项(*work-item*)。每一个工作项都会执行内核函数体。

当OpenCL设备开始执行内核，OpenCL C中提供的内置函数可以让工作项知道自己的编号。编程者调

用`get_global_id(0)`来获取当前工作项的位置。



global_work_size参数指定NDRange在每个维度上有多少个工作项，

local_work_size参数指定WorkGroup在每个维度上有多少个工作组。

所以global_work_size除以local_work_size得到的结果不能大于max work group size



内核执行之前，通常需要将主机端的数据拷贝到OpenCL内存对象的所分配的空间中。创建数组或图像可以调用不同的创建API(`clCreate*()`)。将主机指针作为`clCreate*()`的参数用于初始化OpenCL内存对象。这种方式可以隐式的进行数据传输，并不需要编程者为之担心。内存对象初始化之后，运行时就需要保证数据依据依赖关系，以正确的顺序和时间转移到设备端。



假设我们的内存对象是一个数组，主机端和设备端的内存互传需要使用到下面两个API：`clEnqueueWriteBuffer()`和`clEnqueueReadBuffer()`

clEnqueueWriteBuffer将主机内存传输到设备。

clEnqueueReadBuffer将设备内存传输到主机。



全局内存对于执行内核中的每个工作项都是可见的(类似于CPU上的内存)。当数据从主机端传输到设备端，数据就存储在全局内存中。有数据需要从设备端传回到主机端，那么对应的数据需要存储在全局内存中。其关键字为`global`或`__global`，关键字加在指针类型描述符的前面，用来表示该指针指向的数据存储在全局内存中。



常量内存并非为只读数据设计，但其能让所有工作项同时对该数据进行访问。这里存储的值通常不会变化(比如，某个数据变量存储着π的值)。OpenCL的内存模型中，常量内存为全局内存的子集，所以内存对象传输到全局内存的数据可以指定为“常量”。使用关键字`constant`或`__constant`将相应的数据映射到常量内存。

- 局部内存中的数据，只有在同一工作组内的工作项可以共享。通常情况下，局部内存会映射到片上的物理内存，例如：软件管理的暂存式存储器。比起全局内存，局部内存具有更短的访问延迟，以及更高的传输带宽。调用`clSetKernelArg()`设置局部内存时，只需要传递大小，而无需传递对应的指针，相应的局部内存会由运行时进行开辟。OpenCL内核中，使用`local`或`__local`关键字来描述指针，从而来定义局部内存(例如，`local int *sharedData`)。不过，数据也可以通过关键字`local`，静态申明成局部内存变量(例如，`local int[64]`)。
- 私有内存只能由工作项自己进行访问。局部变量和非指针内核参数通常都在私有内存上开辟。实践中，私有变量通常都与寄存器对应。不过，当寄存器不够私有数组使用是，这些溢出的数据通常会存储到非片上内存(高延迟的内存空间)上。



创建并执行一个简单的OpenCL应用大致需要以下几步：

1. 查询平台和设备信息
2. 创建一个上下文
3. 为每个设备创建一个命令队列
4. 创建一个内存对象(数组)用于存储数据
5. 拷贝输入数据到设备端
6. 使用OpenCL C代码创建并编译出一个程序
7. 从编译好的OpenCL程序中提取内核
8. 执行内核
9. 拷贝输出数据到主机端
10. 释放资源

OpenCL C上的并发执行单元称为工作项(*work-item*)。每一个工作项都会执行内核函数体。这里就不用手动的去划分任务，这里将每一次循环操作映射到一个工作项中。OpenCL运行时可以创建很多工作项，其个数可以和输入输出数组的长度相对应，并且工作项在运行时，以一种默认合适的方式映射到底层硬件上(CPU或GPU核)。当要执行一个内核时，编程者需要指定每个维度上工作项的数量(NDRange)。一个NDRange可以是一维、二维、三维的，其不同维度上的工作项ID映射的是相应的输入或输出数据。NDRange的每个维度的工作项数量由`size_t`类型指定。