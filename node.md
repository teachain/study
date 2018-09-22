# Yarn 

Yarn 对你的代码来说是一个包管理器， 你可以通过它使用全世界开发者的代码， 或者分享自己的代码。

与npm功能差不多的一个工具。

如果已经安装了npm，想换一个方式，可以使用

```
npm install -g yarn
```

安装yarn

裸装

#### MacOS安装

```
curl -o- -L https://yarnpkg.com/install.sh | bash
```

#### Ubuntu安装

```
sudo apt-key adv --keyserver pgp.mit.edu --recv D101F7899D41F3C3 
echo "deb http://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list
sudo apt-get update && sudo apt-get install yarn
```

#### windows

windows 下需要下载msi文件 ，下载地址：<https://yarnpkg.com/latest.msi>



# Node.js 事件循环

## 事件驱动程序

```
// 1、引入 events 模块
var events = require('events');
// 2、创建 eventEmitter 对象
var eventEmitter = new events.EventEmitter();
// 3、绑定事件及事件的处理程序，eventHandler是一个function
eventEmitter.on('eventName', eventHandler);
// 4、触发事件,也就是我触发eventName事件，那么就会调用eventHandler函数
eventEmitter.emit('eventName');

```

我们从这里就可以看到，nodejs使用事件是很简单的，一般就四步。

当事件触发时，注册到这个事件的事件监听器被依次调用，事件参数作为回调函数参数传递。

也就是说对一个事件，你可以绑定多个处理函数，这些函数当事件被触发的时候，会被依次调用。

```
addListener(event, listener)
为指定事件添加一个监听器到监听器数组的尾部。
on(event, listener)
为指定事件注册一个监听器，接受一个字符串 event 和一个回调函数。
once(event, listener)
为指定事件注册一个单次监听器，即 监听器最多只会触发一次，触发后立刻解除该监听器。
removeListener(event, listener)
移除指定事件的某个监听器，监听器必须是该事件已经注册过的监听器。
它接受两个参数，第一个是事件名称，第二个是回调函数名称。
removeAllListeners([event])
移除所有事件的所有监听器， 如果指定事件，则移除指定事件的所有监听器。
setMaxListeners(n)
默认情况下， EventEmitters 如果你添加的监听器超过 10 个就会输出警告信息。 setMaxListeners 函数用于提高监听器的默认限制的数量。
listeners(event)
返回指定事件的监听器数组。
emit(event, [arg1], [arg2], [...])
按参数的顺序执行每个监听器，如果事件有注册监听返回 true，否则返回 false。


```

## error 事件

EventEmitter 定义了一个特殊的事件 error，它包含了错误的语义，我们在遇到 异常的时候通常会触发 error 事件。

当 error 被触发时，EventEmitter 规定如果没有响 应的监听器，Node.js 会把它当作异常，退出程序并输出错误信息。

我们一般要为会触发 error 事件的对象设置监听器，避免遇到错误后整个程序崩溃。

在 Node 应用程序中，执行异步操作的函数将回调函数作为最后一个参数， 回调函数接收错误对象作为第一个参数。

```
var fs = require("fs");

fs.readFile('input.txt', function (err, data) {
   if (err){
      console.log(err.stack);
      return;
   }
   console.log(data.toString());
});
console.log("程序执行完毕");
```

# Node.js模块系统

为了让Node.js的文件可以相互调用，Node.js提供了一个简单的模块系统。

模块是Node.js 应用程序的基本组成部分，文件和模块是一一对应的。换言之，一个 Node.js 文件就是一个模块，这个文件可能是JavaScript 代码、JSON 或者编译过的C/C++ 扩展。

Node.js 提供了 exports 和 require 两个对象，其中 exports 是模块公开的接口，require 用于从外部获取一个模块的接口，即所获取模块的 exports 对象。



# 前端开发桌面程序(electron)

## electron项目和web项目的区别

electron核心我们可以分成2个部分，主进程和渲染进程。主进程连接着操作系统和渲染进程，可以把她看做页面和计算机沟通的桥梁。渲染进程就是我们所熟悉前端环境了。只是载体改变了，从浏览器变成了window。传统的web环境我们是不能对用户的系统就行操作的。而electron相当于node环境，我们可以在项目里使用所有的node api 。

简单理解：
给web项目套上一个node环境的壳。

###### 主进程

在 Electron 里，运行 package.json 里 main 脚本的进程被称为主进程。在主进程运行的脚本可以以创建 web 页面的形式展示 GUI。

###### 渲染进程

由于 Electron 使用 Chromium 来展示页面，所以 Chromium 的多进程结构也被充分利用。每个 Electron 的页面都在运行着自己的进程，这样的进程我们称之为渲染进程。
 在一般浏览器中，网页通常会在沙盒环境下运行，并且不允许访问原生资源。然而，Electron 用户拥有在网页中调用 io.js 的 APIs 的能力，可以与底层操作系统直接交互。

###### 主进程与渲染进程的区别

主进程使用 BrowserWindow 实例创建网页。每个 BrowserWindow 实例都在自己的渲染进程里运行着一个网页。当一个 BrowserWindow 实例被销毁后，相应的渲染进程也会被终止。
 主进程管理所有页面和与之对应的渲染进程。每个渲染进程都是相互独立的，并且只关心他们自己的网页。
 由于在网页里管理原生 GUI 资源是非常危险而且容易造成资源泄露，所以在网页面调用 GUI 相关的 APIs 是不被允许的。如果你想在网页里使用 GUI 操作，其对应的渲染进程必须与主进程进行通讯，请求主进程进行相关的 GUI 操作。
 在 Electron，我们提供用于在主进程与渲染进程之间通讯的 [ipc](https://link.jianshu.com?t=https://github.com/electron/electron/blob/master/docs-translations/zh-CN/api/ipc-main-process.md) 模块。并且也有一个远程进程调用风格的通讯模块 [remote](https://link.jianshu.com?t=https://www.w3cschool.cn/electronmanual/electronmanual-remote.html)。

electron

electron-builder



Error: Command failed: makensis -DVERSIONMAJOR=0 -DVERSIONMINOR=11 -DVERSIONBUILD=2 -DTYPE=mist -DAPPNAME=Mist scripts/windows-installer.nsi

/bin/sh: makensis: command not found

