## vue基础

“{{}}”是最基本的文本插值方法，它会自动将我们双向绑定的数据实时显示出来。

如果有的时候就是想输出 HTML，而不是将数据解释后的纯文本，可以使用 v-html替代v-model来达到目的。

如果想显示{{}}标签，而不进行替换，使用 v-pre即可跳过这个元素和它的子元素的编译过程。

#### 过滤器

Vue. 支持在{{}}插值的尾部添加一小管道符 “ |” 对数据进行过滤，经常用于格式化文
本，比如字母全部大写、货币千位使用逗号分隔等。过滤的规则是自定义的， 通过给 Vue 实例添
加选项filters来设置。

#### 指令

指令(Directives)是 Vue扣 模板中最常用的一项功能，它带有前缀 v-,如v-model,v-if,v-for。

指令的主要职责就是当其表达式的值改变时，相应地将某些行为应用到 DOM 上

v-bind 的基本用途是动态更新 HTML 元素上的属性，也就是原来的html属性都改为v-bind:为前缀，如果src就改为v-bind:src="data中的变量"。

v-on，它用来绑定事件监听器。比如v-on:click="methods中的函数"

Vue扣的v-bind和v-on指令都提供了语法糖，也可以说是缩写， 比如v-bind，可以省略v-bind,

直接写一个冒号 “:”,v-on也可以直接用 “@” 来缩写



vue.js实战 41页

#### 工程化

安装最新的vue cli

```
sudo npm install -g @vue/cli
```

创建项目

```
vue create myapp
```

安装vue-router组件（前端路由）

```
vue add router
```

安装ivew组件 （页面布局）

```
npm install iview --save
```

安装iview-loader（页面加载）

```
npm install iview-loader --save-dev
```

安装  vuex组件(状态管理)

```
 vue add vuex
```

安装 axios组件（http请求）

```
npm install axios --save
```

在项目的根目录下新建 vue.config.js 文件（是根目录，不是src目录）

```
module.exports = {
  // 基本路径
  baseUrl: '/',
  // 输出文件目录
  outputDir: 'dist',
  assetsDir: "static",

  // eslint-loader 是否在保存的时候检查
  lintOnSave: true,
  // use the full build with in-browser compiler?
  // https://vuejs.org/v2/guide/installation.html#Runtime-Compiler-vs-Runtime-only
  // 生产环境是否生成 sourceMap 文件
  productionSourceMap: true,
  // css相关配置
  css: {
    // 是否使用css分离插件 ExtractTextPlugin
    extract: true,
    // 开启 CSS source maps?
    sourceMap: false,
    // css预设器配置项
    loaderOptions: {},
    // 启用 CSS modules for all css / pre-processor files.
    modules: false
  },
  // use thread-loader for babel & TS in production build
  // enabled by default if the machine has more than 1 cores
  parallel: require('os').cpus().length > 1,

  //https://cli.vuejs.org/zh/guide/webpack.html#%E9%93%BE%E5%BC%8F%E6%93%8D%E4%BD%9C-%E9%AB%98%E7%BA%A7
  chainWebpack: config => {

    config.module
      .rule('vue')
      .use('vue-loader')
      .loader('vue-loader')

    config.module
      .rule('vue')
      .use('iview-loader')
      .loader('iview-loader')
      .options({
        prefix: false
      })
  },

  // PWA 插件相关配置
  // see https://github.com/vuejs/vue-cli/tree/dev/packages/%40vue/cli-plugin-pwa
  pwa: {},
  // webpack-dev-server 相关配置
  devServer: {
    open: process.platform === 'darwin',
    host: '0.0.0.0',
    port: 8080,
    https: false,
    hotOnly: false,
    proxy: 'http://localhost:8090', // 配置跨域处理,只有一个代理
    before: app => {}
  },
  // 第三方插件配置
  pluginOptions: {
    // ...
  }
}
```

在main.js中添加

```
import iView from 'iview'
import 'iview/dist/styles/iview.css'
Vue.use(iView)

```

`vue-router` 默认 hash 模式 —— 使用 URL 的 hash 来模拟一个完整的 URL，于是当 URL 改变时，页面不会重新加载。

如果不想要很丑的 hash，我们可以用路由的 **history 模式**，这种模式充分利用 `history.pushState` API 来完成 URL 跳转而无须重新加载页面。

在Router对象里配置

```
mode: 'history',//添加该参数则为history，没有该参数，默认是hash,history需要后台配置支持。
```

