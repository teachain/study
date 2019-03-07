# 常用命令

命令格式：命令 -选项  参数

当有多个选项时，可以写在一起。（例如 ls -al,al就是写在一起的选项）

两个特殊的目录.和..,分别代表当前的目录和当前目录的父目录。

命令 、选项、 参数、它们之间必须使用空格分割。  

- cd  （切换目录）
- ls -ald  (查看目录或文件)
- pwd (print working directory,打印当前目录)
- cp -R src dest （拷贝文件或目录）
- mv  src dest (移动文件，重命名，移动并改名)
- touch filename （创建文件）
- mkdir dir (创建目录)
- rm -r  file or dir （删除文件或目录）
- cat  filename (查看文件内容，显示文件内容)
- more filename
- head -num filename
- tail  -num filename (-f 动态查看文件内容)
- ln -s  src  dest (创建软链接-相当于快捷方式，去掉-s就是硬链接-有cp的效果，并且同步更新)



1、权限管理命令：chmod

它的英文愿意是:change the permissions mode of a file

语法是

```
chmod [{ugo}{+-=}{rwx}][文件或目录]  #方式一  chmod u+rwx ./group
 
chmod [mode=421] [文件或目录]       #方式二 r=4 w=2 x=1,这个是重点。

//比如说 chmod 755 ./group ,它就相当于 chmod rwx-wx-wx ./group 要记住的是每三位代表一类用户
```

它的作用是：改变文件或目录的权限。

文件目录权限总结，要特别注意对目录的意义

| 代表字符 | 权限     | 对文件的含义     | 对目录的含义               |
| -------- | -------- | ---------------- | -------------------------- |
| r        | 读权限   | 可以查看文件内容 | 可以列出目录中的内容       |
| w        | 写权限   | 可以修改文件内容 | 可以在目录中创建、删除文件 |
| x        | 执行权限 | 可以执行文件     | 可以进入目录               |

2、chown 

英文原意：change file ownership

功能描述：改变文件或目录的所有者

语法:chown [用户] [文件或目录]



##用户管理

用户信息文件：/etc/passwd

- 每一行代表一个用户信息，格式为：用户名:密码位:UId:GID:注释性描述:宿主目录:命令解释器

密码文件: /etc/shadow

用户组文件:/etc/group

用户组密码文件:/etc/gshadow

用户配置文件：

- /etc/login.defs
- /etc/default/useradd

新用户信息文件:/etc/skel

登录信息:/etc/motd

linux用户分为三种:

- 超级用户（root,UId=0）
- 普通用户(UID 500~60000)
- 伪用户(UID 1~499)

### 用户组

- 每个用户都至少属于一个用户组。
- 每个用户组可以包含多个用户。
- 同一用户组的用户享有该组共有的权限。



用户权限

所有者 所属组 其他人

rwx  可读 可写  可执行  

linux权限规则:缺省创建的文件不能授予可执行x权限

