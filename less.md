## LESS 原理及使用方式

本质上，LESS 包含一套自定义的语法及一个解析器，用户根据这些语法定义自己的样式规则，这些规则最终会通过解析器，编译生成对应的 CSS 文件。LESS 并没有裁剪 CSS 原有的特性，更不是用来取代 CSS 的，而是在现有 CSS 语法的基础上，为 CSS 加入程序式语言的特性。下面是一个简单的例子：

```
@color: #4D926F; 
#header { 
 color: @color; 
} 
h2 { 
 color: @color; 
}
```

经过编译生成的 CSS 文件如下：

```
#header { 
 color: #4D926F; 
} 
h2 { 
 color: #4D926F; 
}
```

这就像我们C语言的宏替换一样。

LESS 可以直接在客户端使用，也可以在服务器端使用。在实际项目开发中，我们更推荐使用第三种方式，将 LESS 文件编译生成静态 CSS 文件，并在 HTML 文档中应用。

### 客户端

我们可以直接在客户端使用 .less（LESS 源文件），只需要从 [http://lesscss.org](http://lesscss.org/)下载 less.js 文件，然后在我们需要引入 LESS 源文件的 HTML 中加入如下代码

```
<link rel="stylesheet/less" type="text/css" href="styles.less">
```

重要：LESS 源文件的引入方式与标准 CSS 文件引入方式一样，需要注意的是：在引入 .less 文件时，rel 属性要设置为“stylesheet/less”。还有更重要的一点需要注意的是：LESS 源文件一定要在 less.js 引入之前引入，这样才能保证 LESS 源文件正确编译解析。

### 服务器端

LESS 在服务器端的使用主要是借助于 LESS 的编译器，将 LESS 源文件编译生成最终的 CSS 文件，目前常用的方式是利用 node 的包管理器 (npm) 安装 LESS，安装成功后就可以在 node 环境中对 LESS 源文件进行编译。

在项目开发初期，我们无论采用客户端还是服务器端的用法，我们都需要想办法将我们要用到的 CSS 或 LESS 文件引入到我们的 HTML 页面或是桥接文件中，LESS 提供了一个我们很熟悉的功能— Importing。我们可以通过这个关键字引入我们需要的 .less 或 .css 文件。 如：

@import “variables.less”;

.less 文件也可以省略后缀名，像这样：

@import “variables”;

引入 CSS 同 LESS 文件一样，只是 .css 后缀名不能省略。

### 使用编译生成的静态 CSS 文件

我们可以通过 LESS 的编译器，将 LESS 文件编译成为 CSS 文件，在 HTML 文章中引入使用。这里要强调的一点，LESS 是完全兼容 CSS 语法的，也就是说，我们可以将标准的 CSS 文件直接改成 .less 格式，LESS 编译器可以完全识别。



像 JavaScript 中 **arguments**一样，Mixins 也有这样一个变量：**@arguments**。@arguments 在 Mixins 中具是一个很特别的参数，当 Mixins 引用这个参数时，该参数表示所有的变量，很多情况下，这个参数可以省去你很多代码。

