##Gradle##

下载解压后，直接在将文件夹下的bin目录路径配置到path环境变量中，就可以直接使用gradle命令了。

* gradle -v #查看版本号
* gradle -b xxx #指定要加载执行的文件
* gradle -q     #输出日志级别


##Gradle wrapper##

Wrapper，顾名思义，其实就是对Gradle的一层包装，便于在团队开发过程中统一Gradle构建的版本，这样大家都可以使用统一的Gradle版本进行构建，避免因为Gradle版本不统一带来的不必要的问题。
在这里特别介绍的目的是因为我们在项目开发过程中，用的都是wrapper这种方式，而不是使用我们自己下载ZIP压缩包，配置Gradle的环境变量的方式。也就是一般工程里要协作，为了统一而做的一层包装，我们对一个工程进行新建的时候，使用以下命令来创建：

```
 gradle wrapper
```
<font color="red">也就是说我们在项目的根目录下执行这样一条命令，就将wrapper做好了</font>

根目录下就生成了一些目录和文件

* gradle 目录 （下面有wrapper目录，wrapper目录下有gradle-wrapper.jar文件和gradle-wrapper.properties文件）  
* gradlew #linux下的可执行脚本
* gradlew.bat  #windows下的可执行脚本


我们一窥gradle-wrapper.properties这个文件的内容

```
#Sat Feb 04 09:56:40 CST 2017
distributionBase=GRADLE_USER_HOME
distributionPath=wrapper/dists
zipStoreBase=GRADLE_USER_HOME
zipStorePath=wrapper/dists
distributionUrl=https\://services.gradle.org/distributions/gradle-3.3-bin.zip

```

我的机器上配置的gradle环境的版本是3.3的，这里是新建项目，所以我们看到了gradle-3.3-bin.zip,这样的后来的协助者，或者说当我们在服务器上构建的时候，我们并不需要配置gradle的环境，而是根据gradle-wrapper.properties里配置信息进行操作。这样构建就一致了。这些生成的wrapper文件可以作为你项目工程的一部分提交到代码版本控制系统里(git)，这样其他开发人员就会使用这里配置好的统一的gradle进行构建开发。

