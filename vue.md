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

### 基本的指令

- v-once是一个不需要表达式的指令，作用是定义它的元素或组件只渲染一次，包括元素或组件的所有子节点。首次渲染后，不再随数据的变化而重新渲染，将被视为静态内容。
- v-if、v-else-if、v-else。
- v-show,v-show不能用在<template>元素上。
- v-for需要结合in来使用，类似 item in items的形式。也可以配合of使用，类似item of items,也可以是v-for="(item ,index) in items"这种形式。
- 视图会变化的数组方法：push()、pop()、shift()、unshift()、splice()、sort()、reverse()。有些方法是不会改变原数组的方法、filter()、concat()、slice()。
- 事件修饰符，在@绑定的事件后加小圆点"."，再跟一个后缀来使用修饰符。比如
  - .stop
  - .prevent
  - .capture
  - .self
  - .once
- v-model，用于在表单类元素上双向绑定数据。

与事件的修饰符类似，v-model也有修饰符，用于控制数据同步的时机。

1. .lazy
2. .number (把输入转换为Number类型)
3. .trim （可以自动过滤输入的首尾空格）

### 组件

组件需要注册后才可以使用。注册有全局注册和局部注册两种方式。全局注册后，任何Vue实例都可以使用。

```
Vue.component('my-component',{})
```

在Vue实例中，使用components选项可以局部注册组件，注册后的组件只有在该实例作用域下有效。组件中也可以使用components选项来注册组件。

组件不仅仅是要把模板的内容进行复用，更重要的是组件间要进行通信。通常父组件的模板中包含子组件，父组件要正向地向子组件传递数据或参数，子组件接收到后根据参数的不同来渲染不同的内容或执行操作。这个正向传递数据的过程就是通过props来实现的。

在组件中，使用选项props来声明需要从父级接收的数据，props的值可以是两种，一种是字符串数组，一种是对象。 

props中声明的数据与组件data函数return的数据主要区别就是props的来自父级，而data中的是组件自己的数据，作用域是组件本身，这两种数据都可以在模板template及计算属性computed和方法methods中使用。在组件的自定义标签上直接写该props的名称，如果要传递多个数据，在props数组中添加项即可。props:["参数名1","参数名2"]类似这样子。当使用DOM模板时，驼峰命名的props名称要转为短横分割命名（-中间横线）。这里要注意的是，这样的数据传递是单向的，也就是父组件传向子组件，当在子组件中修改数据，是不会反映到父组件的。但也要注意，在javascript中对象和数组是引用类型，指向同一个内存空间，所以props是对象和数组时，在子组件内改变时会影响父组件的。

组件之间的通信方法

父组件—>子组件 通过props

子组件—>父组件 通过自定义事件  子组件用$emit()来触发事件，父组件用$on()来监听子组件的事件。

非父子组件，通过使用一个空的Vue实例作为中央事件总线

```
var bus=new Vue();
bus.$emit("事件"，参数);//触发事件
bus.$on("事件"，方法) //监听事件
```

### slot分发

编译作用域

slot分发的内容，作用域是在父组件上的。

### 监听

watch选项用来监听某个prop或data的改变，当它们发生变化时，就会触发watch配置的函数。watch监听的数据的回调函数有2个参数可用，第一个是新的值，第二个是旧的值。（当然如果你不用，可以不写。）

### 自定义指令

自定义指令的注册方法和组件很像，也分全局注册和局部注册。

```
//全局注册
Vue.directive('focus',{//focus为指令的名称
    //指令选项
});

//局部注册就是作为选项，在vue实例中定义。
var app=new Vue({el:'#app',
directive:{
    focus:{
    //指令选项
    }
}
});
```

自定义指令的选项是由几个钩子函数组成的，每个都是可选的。

- bind:只调用一次，指令第一次绑定到元素时调用，用这个钩子函数可以定义一个在绑定时执行一次的初始化动作。
- inserted:被绑定元素插入父节点时调用（父节点存在即可调用，不必存在于document中）。
- update：被绑定元素所在的模板更新时调用，而不论绑定值 是否变化，通过比较更新前后的绑定值，可以忽略不必要的模板更新。
- componentUpdated:被绑定元素所在模板完成一次更新周期时调用。
- unbind:只调用一次，指令与元素解绑时调用。



注意：

1、Vue提供了一个特殊变量$event,用于访问原生DOM事件。

```
<a href=”http://www.apple.com” @click="handleClick ('禁止打开’，$event)”〉
```

2、使用v-model时，如果是用中文输入法输入中文，一般在没有选定词组钱，也就是在拼音阶段，Vue是不会更新数据的，当敲下汉字时才会触发更新。如果希望总是实时更新，可以用@input来替代v-model。



vue实战144页。

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

