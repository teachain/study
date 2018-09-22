必须安装的依赖有：

- Node(版本必须在 v8.0 以上,注意不要使用 cnpm)
- Watchman
- React Native 命令行工具
-  Xcode(编译 iOS 应用所需的工具和环境。[Xcode](https://developer.apple.com/xcode/downloads/) 9.0 或更高版本)

设置一下仓库

```
npm config set registry https://registry.npm.taobao.org --global
npm config set disturl https://npm.taobao.org/dist --global
```

Yarn是 Facebook 提供的替代 npm 的工具，可以加速 node 模块的下载。

React Native 的命令行工具用于执行创建、初始化、更新项目、运行打包服务（packager）等任务。

```
npm install -g yarn react-native-cli
```

安装完 yarn 后同理也要设置镜像源：

```
yarn config set registry https://registry.npm.taobao.org --global
yarn config set disturl https://npm.taobao.org/dist --global
```

创建应用

```
react-native init your_app_name
```

运行应用

```
react-native run-ios #运行ios
react-native run-android# 运行android
```

