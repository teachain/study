graphviz便于数据可视化，用于生成决策树、流程图

Graphviz http://www.graphviz.org/

mac下安装

```
brew install graphviz
```

文本保存为`hello.dot`

```
digraph pic { 
  Hello -> World
}
```

同一目录下终端运行

```
dot hello.dot -T png -o hello.png
```