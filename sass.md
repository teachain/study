

# CSS预处理器(css preprocessor)

它的基本思想是，用一种专门的编程语言，进行网页样式设计，然后再编译成正常的CSS文件。



[SASS](http://sass-lang.com/)是一种CSS的开发工具，提供了许多便利的写法，大大节省了设计者的时间，使得CSS的开发，变得简单和可维护。



SASS是Ruby语言写的，但是两者的语法没有关系。不懂Ruby，照样使用。只是必须先[安装Ruby](http://www.ruby-lang.org/zh_cn/downloads/)，然后再安装SASS。

假定你已经安装好了Ruby，接着在命令行输入下面的命令：

```
sudo gem install sass
```

然后，就可以使用了。

SASS文件就是普通的文本文件，里面可以直接使用CSS语法。文件后缀名是.scss，意思为Sassy CSS。

下面的命令，可以在屏幕上显示.scss文件转化的css代码。（假设文件名为test。）

> 　　sass test.scss

如果要将显示结果保存成文件，后面再跟一个.css文件名。

> 　　sass test.scss test.css

SASS提供四个[编译风格](http://sass-lang.com/docs/yardoc/file.SASS_REFERENCE.html#output_style)的选项：

> 　　* nested：嵌套缩进的css代码，它是默认值。
>
> 　　* expanded：没有缩进的、扩展的css代码。
>
> 　　* compact：简洁格式的css代码。
>
> 　　* compressed：压缩后的css代码。



生产环境当中，一般使用最后一个选项。

> 　　sass --style compressed test.sass test.css

你也可以让SASS监听某个文件或目录，一旦源文件有变动，就自动生成编译后的版本。

> 　　// watch a file
>
> 　　sass --watch input.scss:output.css
>
> 　　// watch a directory
>
> 　　sass --watch app/sass:public/stylesheets





SASS允许使用变量，所有变量以$开头。

如果变量需要镶嵌在字符串之中，就必须需要写在#{}之中。

SASS允许在代码中使用算式



SASS允许选择器嵌套。比如，下面的CSS代码：

```
　div h1 {
　　　　color : red;
　　}
```

可以写成：

```
div {
　　　　hi {
　　　　　　color:red;
　　　　}
　　}
```



属性也可以嵌套，比如border-color属性，可以写成：

```
p {
　　　　border: {
　　　　　　color: red;
　　　　}
　　}
```

注意，border后面必须加上冒号。

在嵌套的代码块内，可以使用&引用父元素。比如a:hover伪类，可以写成：

```
　a {
　　　　&:hover { color: #ffb3ff; }
　　}
```



完整参看

http://www.ruanyifeng.com/blog/2012/06/sass.html



Sass是一种"CSS预处理器"，可以让CSS的开发变得简单和可维护。但是，只有搭配[Compass](http://compass-style.org/)，它才能显出真正的威力。



**一、Compass是什么？**

简单说，Compass是Sass的工具库（toolkit）。

Sass本身只是一个编译器，Compass在它的基础上，封装了一系列有用的模块和模板，补充Sass的功能。它们之间的关系，有点像Javascript和jQuery、Ruby和Rails、python和Django的关系。

Compass是用Ruby语言开发的，所以安装它之前，必须安装Ruby。

假定你的机器（Linux或OS X）已经安装好Ruby，那么在命令行模式下键入：

```
sudo gem install compass
```

**项目初始化**

接下来，要创建一个你的Compass项目，假定它的名字叫做myproject，那么在命令行键入：

```
　　compass create myproject
```

当前目录中就会生成一个myproject子目录。

进入该目录：

```
cd myproject
```

你会看到，里面有一个[config.rb](https://github.com/thesassway/sass-test/blob/master/config.rb)文件，这是你的项目的[配置文件](http://compass-style.org/help/tutorials/configuration-reference/)。还有两个子目录sass和stylesheets，前者存放Sass源文件，后者存放编译后的css文件。

Compass的编译命令是

```
　compass compile
```





npm install 命令会到目录下查找package.json文件，并根据这个文件下载相应的文件

```
npm install
```

安装bower

```
sudo npm install -g bower
```

bower install 到目录下查找bower.json文件，并根据这个文件下载相应的文件

```
bower install
```

```
//安装gulp
sudo npm install --global gulp
//安装pug
sudo npm install -g pug
//安装 pug-cli
npm install -g pug-cli
```


