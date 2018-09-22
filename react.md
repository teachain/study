



HTML 语言直接写在 JavaScript 语言之中，不加任何引号，这就是 [JSX 的语法](http://facebook.github.io/react/docs/displaying-data.html#jsx-syntax)，它允许 HTML 与 JavaScript 的混写。

遇到 HTML 标签（以 `<` 开头），就用 HTML 规则解析；遇到代码块（以 `{`开头），就用 JavaScript 规则解析。



`<script>` 标签的 `type` 属性为 `text/babel`，凡是使用 JSX 的地方，都要加上 `type="text/babel"` ，也就是

```
如果我们需要使用 JSX，则 <script> 标签的 type 属性需要设置为 text/babel。
就是在javascript中直接写html的代码，不需要用双引号包起来。
也即是这种
<script type="text/babel"> //这里是必须的
//这个表示渲染一个组件
ReactDOM.render(
    <h1>Hello, world!</h1>, //这里的html内容并不需要用""包起来，这就是jsx语法。
    document.getElementById('example')
);
</script>
```



Components must return a single root element

这意思是必须用一个根元素来包含子元素，不能用多个同级元素返回，也就是要求用一个顶层的元素包含起来

就像vue中的最外层的<template></template>

```
//你可以这样写
class Webcome extends React.Component{
  render(){
     return <h1>Hello,{this.props.name}</h1>;
  }
}
//但你不能这么写
class Webcome extends React.Component{
  render(){
     return (<h1>Hello,{this.props.name}</h1>
     <h1>Hello,{this.props.name}</h1>);
  }
}
//当然你就可以这样写咯,感谢root元素div
class Webcome extends React.Component{
  render(){
     return (
      <div>
      <h1>Hello,{this.props.name}</h1>
      <h1>Hello,{this.props.name}</h1>
      </div>
      );
  }
}

```



**All React components must act like pure functions with respect to their props**

这啥意思，就是

```
class Webcome extends React.Component{
  render(){
     //在这里不要去修改props中的值
     return <h1>Hello,{this.props.name}</h1>;
  }
}
function Comment(props) {
  //在这个函数内不要去修改props中的值
  return (
    <div className="Comment">
      <div className="UserInfo">
        <img className="Avatar"
          src={props.author.avatarUrl}
          alt={props.author.name}
        />
        <div className="UserInfo-name">
          {props.author.name}
        </div>
      </div>
      <div className="Comment-text">
        {props.text}
      </div>
      <div className="Comment-date">
        {formatDate(props.date)}
      </div>
    </div>
  );
}
```



Local state is exactly that: a feature available only to classes.

We want to [set up a timer](https://developer.mozilla.org/en-US/docs/Web/API/WindowTimers/setInterval) whenever the `Clock` is rendered to the DOM for the first time. This is called "mounting" in React.

componentDidMount()

We also want to [clear that timer](https://developer.mozilla.org/en-US/docs/Web/API/WindowTimers/clearInterval) whenever the DOM produced by the `Clock` is removed. This is called "unmounting" in React.

componentWillUnmount()



There are three things you should know about `setState()`.

### 1、Do Not Modify State Directly

也就是在想要修改状态的地方不要直接修改属性，this.state=xxx只在构造函数中使用。也就是以下英文的意思。

The only place where you can assign `this.state` is the constructor

```
//错误的做法
this.state.comment = 'Hello';
//正确的做法
this.setState({comment: 'Hello'});
```



事件处理函数

```
//原来html的写法
<button onclick="dothing()"></button>
//jsx的写法
<button onClick={dothing}>
//然后停止事件传播要使用e.preventDefault();
//取代原来的return false

```

provide a listener when the element is initially rendered.

this.handleClick = this.handleClick.bind(this);

event.target.name

event.target.value

event.target.type



## 使用 create-react-app 快速构建 React 开发环境

create-react-app 是来自于 Facebook，通过该命令我们无需配置就能快速构建 React 开发环境。

create-react-app 自动创建的项目是基于 Webpack + ES6 。

```
$ npm install -g create-react-app
$ create-react-app my-app
$ cd my-app/
$ npm start
```

我们用 React 开发应用时一般只会定义一个根节点。但如果你是在一个已有的项目当中引入 React 的话，你可能会需要在不同的部分单独定义 React 根节点。

要将React元素渲染到根DOM节点中，我们通过把它们都传递给 ReactDOM.render() 的方法来将其渲染到页面上：

```
const element = <h1>Hello, world!</h1>;
ReactDOM.render(
    element,
    document.getElementById('example')
);
```

**React 只会更新必要的部分**

值得注意的是 React DOM 首先会比较元素内容先后的不同，而在渲染过程中只会更新改变了的部分。

我们可以在 JSX 中使用 JavaScript 表达式。表达式写在花括号 **{}** 中

```
ReactDOM.render(
    <div>
      <h1>{1+1}</h1> //{}这里是javascript的表达式
    </div>
    ,
    document.getElementById('example')
);
```

React 的 JSX 使用大、小写的约定来区分本地组件的类和 HTML 标签。

比如<App></App> //react  <div></div> //html



我们可以使用函数定义了一个组件：

```
function HelloMessage(props) {
    return <h1>Hello World!</h1>;
}
```

你也可以使用 ES6 class 来定义一个组件:

```
class Welcome extends React.Component {
  render() {
    return <h1>Hello World!</h1>;
  }
}
```

注意，在添加属性时， class 属性需要写成 className ，for 属性需要写成 htmlFor ，这是因为 class 和 for 是 JavaScript 的保留字。



向组件传递参数，可以使用 **this.props** 对象，也就是从父级元素从调用时通过属性传入的，都是用这个对象来访问**this.props**。



React 里，只需更新组件的 state，然后根据新的 state 重新渲染用户界面（不要操作 DOM）。

```
始终调用this.setState({})来修改数据。//一定要记住这样用。
```

state 和 props 主要的区别在于 **props** 是不可变的，而 state 可以根据与用户交互来改变。这就是为什么有些容器组件需要定义 state 来更新和修改数据。 而子组件只能通过 props 来传递数据。

组件类的 defaultProps 属性为 props 设置默认值

```
class HelloMessage extends React.Component {
  render() {
    return (
      <h1>Hello, {this.props.name}</h1>
    );
  }
}
 
HelloMessage.defaultProps = {
  name: 'Runoob'
};
 
const element = <HelloMessage/>;
 
ReactDOM.render(
  element,
  document.getElementById('example')
);
```



### 将生命周期方法添加到类中

componentDidMount() 与 componentWillUnmount() 方法被称作生命周期钩子。



# React 事件处理

React 元素的事件处理和 DOM 元素类似。但是有一点语法上的不同:

- React 事件绑定属性的命名采用驼峰式写法，而不是小写。

- 如果采用 JSX 的语法你需要传入一个函数作为事件处理函数，而不是一个字符串(DOM 元素的写法)

  ```
  <button onClick={activateLasers}> //驼峰法，使用{},不需要双引号包含
    激活按钮
  </button>
  ```

  在 React 中另一个不同是你不能使用返回 **false** 的方式阻止默认行为， 你必须明确的使用 preventDefault。 



当你使用 ES6 class 语法来定义一个组件的时候，事件处理器会成为类的一个方法。

```
class Toggle extends React.Component {
  constructor(props) {
    super(props);
    this.state = {isToggleOn: true};
 
    // 这边绑定是必要的，这样 `this` 才能在回调函数中使用
    this.handleClick = this.handleClick.bind(this);
  }
 
  handleClick() {
    this.setState(prevState => ({
      isToggleOn: !prevState.isToggleOn
    }));
  }
 
  render() {
    return (
      <button onClick={this.handleClick}>
        {this.state.isToggleOn ? 'ON' : 'OFF'}
      </button>
    );
  }
}
 
ReactDOM.render(
  <Toggle />,
  document.getElementById('example')
);
```

## redux

##### action

##### Action 本质上是 JavaScript 普通对象。我们约定，action 内必须使用一个字符串类型的 `type` 字段来表示将要执行的动作。多数情况下，`type` 会被定义成字符串常量。当应用规模越来越大时，建议使用单独的模块或文件来存放 action。

一般来说你会通过 [`store.dispatch()`](http://www.redux.org.cn/docs/api/Store.html#dispatch) 将 action 传到 store。(这意思就是我要改变数据了)

actions 只是描述了*有事情发生了*这一事实，并没有描述应用如何更新 state。

用store.dispatch()分发action(相当于把事件广播出去)

##### reducer(就是对action进行处理)

reducer 就是一个纯函数，接收旧的 state 和 action，返回新的 state。（相当于收到事件了以后，做一些改变）

```js
(previousState, action) => newState
```

使用 [action](http://www.redux.org.cn/docs/basics/Actions.html) 来描述“发生了什么”，和使用 [reducers](http://www.redux.org.cn/docs/basics/Reducers.html) 来根据 action 更新 state 的用法。

**Store** 就是把它们联系到一起的对象。Store 有以下职责：

- 维持应用的 state；
- 提供 [`getState()`](http://www.redux.org.cn/docs/api/Store.html#getState) 方法获取 state；
- 提供 [`dispatch(action)`](http://www.redux.org.cn/docs/api/Store.html#dispatch) 方法更新 state；
- 通过 [`subscribe(listener)`](http://www.redux.org.cn/docs/api/Store.html#subscribe) 注册监听器;
- 通过 [`subscribe(listener)`](http://www.redux.org.cn/docs/api/Store.html#subscribe) 返回的函数注销监听器。