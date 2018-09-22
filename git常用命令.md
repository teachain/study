创建版本库

cd到你的指定文件夹

然后执行以下命令即可创建版本库

```
git init
```
把一个文件放到Git仓库只需要两步

第一步，用命令git add告诉Git，把文件添加到仓库

```
git add readme.txt(这是你需要添加的文件名)  

```
第二步，用命令git commit告诉Git，把文件提交到仓库

```
git commit -m "wrote a readme file" (带备注提交)
```

查看当前版本库的状态

```
git status
```

查看比对版本库的修改

```
git diff

```

注意：<font color="red">提交修改和提交新文件是一样的，都需要两步</font>

查看commit日志

```
git log
```

回退版本

```
git reset --hard HEAD^ (git reset --hard commit_id)
```

在Git中，用HEAD表示当前版本,上一个版本就是HEAD^，上上一个版本就是HEAD^^，当然往上100个版本写100个^比较容易数不过来，所以写成HEAD~100。

查看你在版本库里使用的命令历史

```
git reflog

```


工作区有一个隐藏目录.git，这个不算工作区，而是Git的版本库。

Git的版本库里存了很多东西，其中最重要的就是称为stage（或者叫index）的暂存区，还有Git为我们自动创建的第一个分支master，以及指向master的一个指针叫HEAD。
我们把文件往Git版本库里添加的时候，是分两步执行的：

第一步是用git add把文件添加进去，实际上就是把文件修改添加到暂存区；

第二步是用git commit提交更改，实际上就是把暂存区的所有内容提交到当前分支。

因为我们创建Git版本库时，Git自动为我们创建了唯一一个master分支，所以，现在，git commit就是往master分支上提交更改。


注意：<font color="red">每次修改，如果不add到暂存区，那就不会加入到commit中</font>(也就是说每次提交，都是提交暂存区的内容，如果你不把修改add到暂存区中，就不会有东西提交，一旦你执行了git commit就会把暂存区的内容提交，并清空暂存区)


丢弃工作区的修改

```
git checkout -- file可以丢弃工作区的修改

```

<font color="red">git checkout -- file命令中的--很重要，没有--，就变成了“切换到另一个分支”的命令</font>

命令git reset HEAD file可以把暂存区的修改撤销掉（unstage），重新放回工作区

```
git reset HEAD readme.txt

```

git reset命令既可以回退版本，也可以把暂存区的修改回退到工作区。当我们用HEAD时，表示最新的版本

git checkout其实是用版本库里的版本替换工作区的版本，无论工作区是修改还是删除，都可以“一键还原”。

从版本库中删除该文件，那就用命令git rm删掉，并且git commit

```
git rm
git commit
```

把一个已有的本地仓库与之关联

```
git remote add origin git@github.com:你自己的GitHub账户名/learngit.git

```


把本地库的内容推送到远程，用git push命令，实际上是把当前分支master推送到远程。

```
git push -u origin master
```

小结：

要关联一个远程库，使用命令git remote add origin git@server-name:path/repo-name.git；

关联后，使用命令git push -u origin master第一次推送master分支的所有内容；

此后，每次本地提交后，只要有必要，就可以使用命令git push origin master推送最新修改；









