在包含有c代码的go文件中，所以c的编译参数什么的，就是在go文件中设置的。

```
//#cgo windows LDFLAGS:-LC:/cuda/v5.0/lib/x64 
//#cgo windows CFLAGS: -IC:/cuda/v5.0/include 
```

CFLAGS指示了头文件地址

LDFLAGS表明了库文件地址

-I (大写)指示了头文件目录

-L 指示了库文件目录 

-l(L小写)指示所用的具体的某个库文件






