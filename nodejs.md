

Node.js®是一个基于[Chrome V8 引擎](https://developers.google.com/v8/)的 JavaScript 运行时。 Node.js 使用高效、轻量级的事件驱动、非阻塞 I/O 模型。Node.js 之生态系统是目前最大的开源包管理系统。

nodejs中含有npm

其实npm是nodejs的包管理器（package manager）。我们在[node.js](http://lib.csdn.net/base/nodejs)上开发时，会用到很多别人已经写好的javascript代码，如果每当我们需要别人的代码时，都根据名字搜索一下，下载源码，解压，再使用，会非常麻烦。于是就出现了包管理器npm。大家把自己写好的源码上传到npm官网上，如果要用某个或某些个，直接通过npm安装就可以了，不用管那个源码在哪里。并且如果我们要使用模块A，而模块A又依赖模块B，模块B又依赖模块C和D，此时npm会根据依赖关系，把所有依赖的包都下载下来并且管理起来。试想如果这些工作全靠我们自己去完成会多么麻烦！



Grunt和所有grunt插件都是基于nodejs来运行的.

### 安装grunt-CLI

“CLI”被翻译为“命令行”。要想使用grunt，首先必须将grunt-cli安装到全局环境中，使用nodejs的“npm install…”进行安装。也就是以下命令：

```
npm install -g grunt-cli
```

grunt是一个“构建工具”。Grunt没有具体的作用，但是它能把有具体作用的一个一个插件组合起来，形成一个整体效应。grunt编译依赖文件

- Gruntfile.js 
- package.json



```
npm install grunt --save-dev
```

—save-dev”的意思是，在当前目录安装grunt的同时，顺便把grunt保存为这个目录的开发依赖项。执行完该命令后，在你的项目下就会多一个node_modules文件夹，它是存储grunt源文件的地方。



一些插件



- Contrib-jshint——javascript语法错误检查；
- Contrib-watch——实时监控文件变化、调用相应的任务重新执行；
- Contrib-clean——清空文件、文件夹；
- Contrib-uglify——压缩javascript代码
- Contrib-copy——复制文件、文件夹
- Contrib-concat——合并多个文件的代码到一个文件中
- karma——前端自动化测试工具