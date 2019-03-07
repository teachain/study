2>&1

1代表标准输出,2代表标准错误

那么命令后面跟 2>&1意思就是把命令执行的标准错误和标准输出都重定向到mr.log里面去.



```
nohup  geth \

--fast --cache=1024 \

--datadir=public_fast \

1>>eth_public_fast.log 2>>eth_public_fast_error.log &
```

这里我们看到将标准输出重定向到文件eth_public_fast.log中

将错误输出重定向到文件eth_public_fast_error.log中



