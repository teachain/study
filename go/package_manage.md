### go的包管理工具

1. govendor
2. godep
3. go mod

### 1、govendor

#### govendor安装

```
go get -u github.com/kardianos/govendor
```

go vendor 是glang1.6正式引入管理包依赖的方式

其基本思路是，将引用的外部包的源代码放在当前工程的vendor目录下面

编译go代码会优先从vendor目录先寻找依赖包

1.解决的问题：
将源码拷贝到当前工程的vendor目录下，这样打包当前的工程代码到任意机器的$GOPATH/src下都可以通过编译，避免项目代码外部依赖过多，迁移后，
需要多次go get 外包依赖包；而且通过go get 重新拉去的外部依赖包的版本可能和工程开发时使用的不一致，导致编译错误问题。

2.未解决的问题：
无法精确的引用外部包进行版本控制，不能指定引用某个特定版本的外部包；只是在开发时，将其拷贝过来，但是一旦外部包升级,vendor下的代码不会跟着升级，
而且vendor下面并没有元文件记录引用包的版本信息，这个引用外部包升级产生很大的问题，无法评估升级带来的风险；

#### govendor使用

在项目的根目录下执行

```
govendor init
```

这个命令执行之后，会在根目录下生成vendor目录和vendor/vendor.json文件

## Quickstart

```
# Setup your project.
cd "my project in GOPATH"
# 初始化 vendor 目录, project 下出现 vendor 目录
govendor init

# Add existing GOPATH files to vendor.
govendor add +external #把依赖都载入到vendor中，前提是已经go get 下来了

# View your work.
govendor list

# Look at what is using a package
govendor list -v fmt

# Specify a specific version or revision to fetch，这种方式直接在vendor目录中，不会在gopath下
govendor fetch golang.org/x/net/context@a4bbce9fcae005b22ae5443f6af064d80a6f5a55

# Get latest v1.*.* tag or branch.
govendor fetch golang.org/x/net/context@v1   

# Get the tag or branch named "v1".
govendor fetch golang.org/x/net/context@=v1  

# Update a package to latest, given any prior version constraint
govendor fetch golang.org/x/net/context

# Format your repository only
govendor fmt +local

# Build everything in your repository only
govendor install +local

# Test your repository only
govendor test +local
```

## Sub-commands

```
init     创建 vendor 文件夹和 vendor.json 文件
list     列出已经存在的依赖包
add      从 $GOPATH 中添加依赖包，会加到 vendor.json
update   从 $GOPATH 升级依赖包
remove   从 vendor 文件夹删除依赖
status   列出本地丢失的、过期的和修改的package
fetch   从远端库增加新的，或者更新 vendor 文件中的依赖包
sync     Pull packages into vendor folder from remote repository with revisions
migrate  Move packages from a legacy tool to the vendor folder with metadata.
get     类似 go get，但是会把依赖包拷贝到 vendor 目录
license  List discovered licenses for the given status or import paths.
shell    Run a "shell" to make multiple sub-commands more efficient for large projects.

go tool commands that are wrapped:
      `+<status>` package selection may be used with them
    fmt, build, install, clean, test, vet, generate, tool
```

##### 在做代码提交的时候，我们需要把vendor目录提交。

### 2、godep

godep和dep是两个工具，不是同一个工具。

安装

```
go get -u -v github.com/tools/godep
```

成功安装后，在`$GOPATH的bin目录`下会有一个godep可执行的二进制文件，后面执行的命令都是用这个，建议这个目录加入到全局环境变量中

#### 以下命令都是在工程的根目录运行

#### 项目中只有Godeps.json文件，而没有包含第三库文件时使用godep restore

```
godep restore
```

这种命令适用的情况是：项目中只有Godeps.json文件，而没有包含第三库。

可以使用godep restore这个命令将所有的依赖库下载到`$GOPATH\src`中 ，也就是godep会按照`Godeps/Godeps.json`内列表，依次执行`go get -d -v`来下载对应依赖包到GOPATH/src路径下。

#### 项目已经可以构建，需要提交到仓库，别人拉下来就可以直接使用的情况下，我们使用godep save来做。

```
godep save
```

- 自动扫描当前目录所属包中import的所有外部依赖库（非系统库）
- 将所有的依赖库下来下来到当前工程中，产生文件 `Godeps\Godeps.json` 文件
- 在没有 `Godeps` 文件的情况下，生成模组依赖目录`vendor`文件夹
- 在这种情况下，vendor目录下并没有vendor.json文件。

这个命令达成的效果是：在项目根目录下生成Godeps\Godeps.json 文件和vendor文件夹，会把第三方依赖都放入vendor中。

也就是一般我们开发的时候，在没有用govendor的情况下，我们是通过go get这种方式来下载第三方依赖的，这样第三方依赖并没有被包含在项目中，在这种情况下，如果我们直接把项目提交了，别人拉取你的代码以后，还是需要逐个go get各个第三方依赖。为了避免这种情况的出现，我们可以在go get到第三方依赖以后，项目已经能够成功go build了，那么我们就可以使用godep save命令来生成Godeps\Godeps.json` 文件和vendor`文件夹，然后把Godeps\Godeps.json` 文件和vendor`文件夹都提交到仓库中，别人拉取代码下来就可以直接编译了。

当项目下存在Godeps\Godeps.json` 文件和vendor`文件夹之后，我们可以直接在原来的go run 这些命令前加godep来运行。

##### 在做代码提交的时候，我们需要把Godeps\Godeps.json` 文件和vendor`文件夹都提交。

## 3、go mod

首先确认你的go编译器的版本必须1.11版本以上，这个才不用配置环境变量就可以玩。

想要使用go mod，首先项目根目录下必须先有go.mod文件，这个可以通过以下命令来达成：

```
go mod init 模块名
```

这个命令执行完毕之后，就会在当前的目录下生成go.mod文件。

当使用go build命令的时候，就会在目录下生成go.sum文件。

当前目录下有了go.mod文件，go compiler将工作在module-aware模式下，自动分析项目的依赖、确定

项目依赖包的初始版本，并下载这些版本的依赖包缓存到特定目录下（目前是存放在$GOPATH/pkg/mod下面）

从这里可以看出，go build就是根据目录下有没有go.mod这个文件来判断是使用module-aware模式还是使用

gopath模式（gopath模式就是在$GOPATH/src下面找相关的import package）。

### 用why解释为何依赖，给出依赖路径

go.mod中的依赖项由go相关命令自动生成和维护。但是如果开发人员想知道为什么会依赖某个package，可以通过go mod why命令来查询原因。go mod why命令默认会给出一个main包到要查询的packge的最短依赖路径。如果go mod why使用 -m flag，则后面的参数将被看成是module，并给出main包到每个module中每个package的最短依赖路径（如果依赖的话）