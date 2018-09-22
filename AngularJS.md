## AngularJS 

 Angular中的每一样东西都是用依赖注入(DI)的方式来创建和使用的，比如指令(Directive)、过滤器(Filter)、控制器(Controller)、服务(Service)。 在Angular中，依赖注入(DI)的容器(container)叫做"[注入器(injector)](http://www.angularjs.net.cn/tutorial/17.html)"。

要想进行依赖注入，你必须先把这些需要协同工作的对象和函数注册(Register)到某个地方。在Angular中，这个地方叫做[“]()[模块(module)](http://www.angularjs.net.cn/tutorial/6.html)”。

对于 HTML 应用程序，通常建议把所有的脚本都放置在 <body> 元素的最底部。

这会提高网页加载速度，因为 HTML 加载不受制于脚本加载。



当网页加载完毕，AngularJS 自动开启。

**ng-app** 指令告诉 AngularJS，<div> 元素是 AngularJS **应用程序** 的"所有者"。

**ng-model** 指令把输入域的值绑定到应用程序变量 **name**。

**ng-bind** 指令把应用程序变量 name 绑定到某个段落的 innerHTML。

AngularJS 指令是以 **ng** 作为前缀的 HTML 属性。ng-init 指令初始化 AngularJS 应用程序变量。

AngularJS 表达式写在双大括号内：**{{ expression }}**。

AngularJS 表达式把数据绑定到 HTML，这与 **ng-bind** 指令有异曲同工之妙。

AngularJS 将在表达式书写的位置"输出"数据。

**AngularJS 表达式** 很像 **JavaScript 表达式**：它们可以包含文字、运算符和变量。

AngularJS **模块（Module）** 定义了 AngularJS 应用。

AngularJS **控制器（Controller）** 用于控制 AngularJS 应用。 

**ng-app**指令定义了应用,  **ng-controller** 定义了控制器。





## 模块

模块定义了一个应用程序。 

模块是应用程序中不同部分的容器。

模块是应用控制器的容器。

控制器通常属于一个模块。

##### 创建模块

你可以通过 AngularJS 的 **angular.module** 函数来创建模块

```
<div ng-app="myApp">...</div>

<script>

var app = angular.module("myApp", []); 

</script>
```



## 控制器



## 指令

指令本质上就是AngularJS扩展具有自定义功能的HTML元素的途径。


通过Angularjs模块API中的.directive()方法，我们可以通过传入一个字符串和一个函数来注册一个新指令，其中字符串是这个指令的名字，指令名应该是驼峰命名风格的，函数应该返回一个对象。

- ng-show 的指令值可以true或false,true时，控制的元素显示，false时，控制的元素隐藏。
- ng-hide  true时隐藏，false时显示
- ng-true-value 用在多选框中选中时的表示值
- ng-false-value 用在多选框中不选中时的表示值
- ng-options="v.value as v.name for v in data",value和name是v的属性。
- ng-init
- ng-disabled
- ng-click
- ng-if
- ng-model
- ng-repeat="v in data" 或者 "(k,v) in data"
- ng-href
- ng-src
- ng-checked
- ng-readonly
- ng-selected
- ng-class
- ng-style







## 服务

在 AngularJS 中，服务是一个函数或对象，可在你的 AngularJS 应用中使用。

```
//第一种方式
var m=angular.module('myapp',[])
//自定义服务
m.factory('myservice',['$http',function($http){
  return {}
}]);
//第二种方式
m.service('myservice2',['$http',function($http){
  ......
}]);
```



## $scope



## $watch

$watch函数会监视$scope上的某个属性。只要属性发生变化就会调用 对应的函数。可以使用$watch函数在$scope上某个属性发生变化时直接运行一个自定 义函数。



## $digest



## $parse





## $interpolte



## 工具函数

- angular.fromJson() 将json字符串转换成javascript对象
- angualr.toJson()将javascript对象转换成json字符串

## 指令作用域

当在指令的构造行数中，定义scope属性时，这个scope将会成为隔离作用域，指令中就访问不到外部的属性了。使用@符号可以达到单项文本绑定，也就是说可以从元素的属性中获取到scope中。使用=符号可以达到双向绑定。用&来绑定父级中定义的函数。

使用link属性来进行dom操作

```
var m=angular.module('hd',[]);
m.directive('myDirective',[function(){
  return {
      restrict:'E',
      link:function(scope,element,atrribute){
        
      }
  }
}]);
```

