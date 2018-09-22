首先需要安装以下依赖：

- [Node.js](https://nodejs.org/) `v7.x` (use the preferred installation method for your OS)
- [Meteor](https://www.meteor.com/install) javascript app framework
- [Yarn](https://yarnpkg.com/) package manager



修改备注：

一、版本更新(注意多语言文件的修改)

1、版本号

2、下载链接

3、md5

跳过更新，下载新版本

"____name__"表示的是app的名称

修改语言文件（interface/i18n/）



interface/client/lib/helpers/helperFunction.js 中的

Helpers.detectNetwork 为绑定网络

检查更新链接设置

modules/clientBinaryManager.js（这个是管geth的）

const BINARY_URL = 'https://raw.githubusercontent.com/ethereum/mist/master/clientBinaries.json';

modules/updateChecker.jso（这个是管mist或wallet的）

https://api.github.com/repos/ethereum/mist/releases/latest 

modules/menuItems.js

https://gitter.im/ethereum/mist

modules/constants.js

https://mainnet.infura.io



目前存在问题

当以--mode  wallet 运行时，在mac下是一片空白。

要在mac 下编译win的安装文件，需要brew install makensis

### Run Mist

打开一个terminal window,输入以下命令:

```
yarn dev:meteor
```

等到出现Client modified -- refreshing 这行内容之后，再另开一个terminal window，输入以下命令

```
yarn dev:electron
```

这个时候，如果能够正常运行起来的话，你将看到mist的窗口了。

进入到meteor-dapp-wallet/app目录下，执行以下命令

```
git submodule update --recursive --init
meteor npm install --save @babel/runtime@7.0.0-beta.49
```

