给一个模块加上

```
if __name == "__main__":
```

常常是为了测试这个模块，因为这个语句块只有当module被作为script直接传给Interpreter的时候才会被执行。

### 模块包(package)

包(package)可以理解为是组织起来的module的一个层次结构，也就是package是一个directory，它包含sub-package或者是module，而module是.py文件，要让Python Interpreter把一个目录作为package，则该目录下必须有__init__.py文件，__init__.py可以为空，当然也可以有对象定义和语句，用来做初始化工作，__init__.py还有个作用就是设置__all__变量。

package本身就可以来作为一个module使用，只是它所包含的sub-module或module可以通过package name用package.module的名称形式去引用，这更有利于组织一系列相关的module，避免module间定义的名称的混乱。

package在实际工程中非常常用，__init__.py也常常不会为空，而会有对象定义和初始化代码来让这个包，也就是这个module，包含其该有的item定义。