## 命令行的参数设置

1、--mine和--dev，当设置了这两个参数之一，将会查找--gasprice(老的)或--miner.gasprice(新的)参数来设置txpool中的gasPrice，并且根据--minerthreads或--miner.threads来设置挖矿的线程数，然后启动挖矿。

2、—etherbase(旧的)或—miner.etherbase(新的)表示设置矿工地址，如果不设置这个参数的话，那么默认找数据目录下的keystore文件下的账户，选择第一个账户为矿工的地址。如果没有账户，则报错。



