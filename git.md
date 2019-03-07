```


git add .   // 添加文件到版本库（只是添加到缓存区），.代表添加文件夹下所有文件 

git add filename //添加文件到版本库

git commit -m "first commit" // 把添加的文件提交到版本库，并填写提交备注

git remote add origin 你的远程库地址  // 把本地库与远程库关联

git push -u origin master    // 第一次推送时

git push origin master  // 第一次推送后，直接使用该命令即可推送修改

git status #查看当前版本库的状态

git diff filepath #查看文件修改内容

#提交新文件
git add 新添加的文件
git commit -m "注释"

#提交修改
git add 修改的文件
git commit -m "注释"

git log  #查看提交日志

#回滚版本
令git reset --hard commit_id #git reset --hard 1094a(这是commit id)

#用来记录你的每一次命令(用这个命令来查看历史命令)
git reflog

git reset HEAD <file>  #可以把暂存区的修改撤销掉（unstage），重新放回工作区



```



创建目录，并进入目录

```
mkdir myapp
cd myapp
```

初始化版本库

```
git init
```

添加文件到版本库



本地操作(下面的操作只是在本地磁盘上操作，还没有涉及到远程版本库)

git add 表示的是把文件添加到暂存区(stage)

git commit 表示的是把暂存区的所有内容提交到当前分支

git管理的是修改，而不是文件。必须先添加到stage才能提交到分支，如果没有进行add 修改,不管你怎么修改，都

不会commit  ,工作区--->stage--->分支





命令`git checkout -- readme.txt`意思就是，把`readme.txt`文件在工作区的修改全部撤销，这里有两种情况：

一种是`readme.txt`自修改后还没有被放到暂存区，现在，撤销修改就回到和版本库一模一样的状态；

一种是`readme.txt`已经添加到暂存区后，又作了修改，现在，撤销修改就回到添加到暂存区后的状态。

总之，就是让这个文件回到最近一次`git commit`或`git add`时的状态。



场景1：当你改乱了工作区某个文件的内容，想直接丢弃工作区的修改时，用命令`git checkout -- file`。

场景2：当你不但改乱了工作区某个文件的内容，还添加到了暂存区时，想丢弃修改，分两步，第一步用命令`git reset HEAD <file>`，就回到了场景1，第二步按场景1操作。

场景3：已经提交了不合适的修改到版本库时，想要撤销本次提交，参考[版本回退](https://www.liaoxuefeng.com/wiki/0013739516305929606dd18361248578c67b8067c8c017b000/0013744142037508cf42e51debf49668810645e02887691000)一节，不过前提是没有推送到远程库。





从版本库中删除该文件，那就用命令git rm删掉，并且git commit：



```
git push -u origin master #推送到远程版本库
```



要关联一个远程库，使用命令`git remote add origin git@server-name:path/repo-name.git`；

关联后，使用命令`git push -u origin master`第一次推送master分支的所有内容；

此后，每次本地提交后，只要有必要，就可以使用命令`git push origin master`推送最新修改；





要克隆一个仓库，首先必须知道仓库的地址，然后使用`git clone`命令克隆。

Git支持多种协议，包括`https`，但通过`ssh`支持的原生`git`协议速度最快。

git clone url





创建`dev`分支，然后切换到`dev`分支

```
git checkout -b dev
```

`git checkout`命令加上`-b`参数表示创建并切换，相当于以下两条命令：

```
$ git branch dev
$ git checkout dev
```

用`git branch`命令查看当前分支

```
git branch  #查看本地分支
git branch -r #查看远端分支
```

`dev`分支的工作成果合并到`master`分支上

```
git merge dev
```

删除`dev`分支

```
git branch -d dev
```



当Git无法自动合并分支时，就必须首先解决冲突。解决冲突后，再提交，合并完成。

解决冲突就是把Git合并失败的文件手动编辑为我们希望的内容，再提交。

用`git log --graph`命令可以看到分支合并图。



查看远程库的信息

```
git remote
```



显示更详细的信息

```
git remote -v
```

```
origin	https://git.coding.net/baoquan2017/ethereum.git (fetch)
origin	https://git.coding.net/baoquan2017/ethereum.git (push)
```

上面显示了可以抓取和推送的`origin`的地址。如果没有推送权限，就看不到push的地址。

推送分支

```
git push origin master  #推送master分支到远端仓库 origin
git push -u origin 1.8dev_rename #推送1.8dev_rename分支到远端仓库 origin
```

origin表示的是远程版本库的名称，用git remote可以查看得到

master表示的是分支名称，用git branch可以查看得到



tag就是一个让人容易记住的有意义的名字，它跟某个commit绑在一起。

首先，切换到需要打标签的分支上

```
git checkout master
```

打一个新标签

```
git tag v1.0
```

查看所有标签

```
git tag
```

给指定的commit打上标签

```
git tag v0.9 f52c633
```

查看标签信息

```
git show v0.9
```

用`-a`指定标签名，`-m`指定说明文字

```
git tag -a v0.1 -m "version 0.1 released" 1094adb
```

删除标签

```
git tag -d v0.1
```

推送某个标签到远程

```
git push origin v1.0
```

一次性推送全部尚未推送到远程的本地标签：

```
git push origin --tags
```

标签已经推送到远程，要删除远程标签就麻烦一点，先从本地删除：

```
git tag -d v0.9
```

然后，从远程删除。删除命令也是push

```
git push origin :refs/tags/v0.9
```

把远程分支拉到本地

git fetch origin dev（dev为远程仓库的分支名）

工作区文件误删，恢复的方法

```
git reset HEAD   误删文件或文件夹路径 （它的路径可以通过git status来查看）

git checkout     误删文件或文件夹路径
```





**1、先将本地修改存储起来**

```
git stash
```

2、用git stash list可以看到保存的信息

```
git stash list
```

3、暂存了本地修改之后，就可以pull了

```
git pull
```

**4、还原暂存的内容**

```
git stash pop stash@{0}
```



#### 标签

tagname为标签名

1、打一个新标签，当然先切换到你需要打标签的分支上

```
#master是你要打标签的分支
git checkout master 

#tagname就是你的标签名，其实它是git tag tagname commitid 的简写，
#因为它使用了默认最新的commitid
git tag tagname 
```

还可以创建带有说明的标签，用`-a`指定标签名，`-m`指定说明文字

```
git tag -a V0.1 -m "version 0.1 released" 1186b5e
```

2、使用git tag 可以查看现有的标签

```
git tag
```

3、查看标签

```
git show tagname
```

4、删除本地仓库中的标签（用git一定要时刻保有本地和远端的概念）

```
git tag -d  tagname
```

5、推送tag到远端仓库

```
git push origin tagname
```

6、删除远程标签

```
#要分两步，先删除本地，在删除远端
#删除本地
git tag -d tagname
#删除远端
git push origin :refs/tags/tagname

```

