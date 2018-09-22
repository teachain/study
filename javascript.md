# javascript

现代web应用程序一般都把全部javascript引用放在<body>元素中页面的内容后面。

在JavaScript中，用var申明的变量实际上是有作用域的。

由于JavaScript的函数可以嵌套，此时，内部函数可以访问外部函数定义的变量，反过来则不行.

JavaScript的函数在查找变量时从自身函数定义开始，从“内”向“外”查找。如果内部函数定义了与外部函数重名的变量，则内部函数的变量将“屏蔽”外部函数的变量。

JavaScript的函数定义有个特点，它会先扫描整个函数体的语句，把所有申明的变量“提升”到函数顶部

```
'use strict';

function foo() {
    var x = 'Hello, ' + y;
    alert(x);
    var y = 'Bob';
}

foo();
```

不在任何函数内定义的变量就具有全局作用域。实际上，JavaScript默认有一个全局对象`window`，全局作用域的变量实际上被绑定到`window`的一个属性.

为了解决块级作用域，ES6引入了新的关键字`let`，用`let`替代`var`可以申明一个块级作用域的变量：

ES6标准引入了新的关键字`const`来定义常量，`const`与`let`都具有块级作用域.

要指定函数的`this`指向哪个对象，可以用函数本身的`apply`方法，它接收两个参数，第一个参数就是需要绑定的`this`变量，第二个参数是`Array`，表示函数本身的参数

另一个与`apply()`类似的方法是`call()`，唯一区别是：

- `apply()`把参数打包成`Array`再传入；
- `call()`把参数按顺序传入。

#### 高阶函数

JavaScript的函数其实都指向某个变量。既然变量可以指向函数，函数的参数能接收变量，那么一个函数就可以接收另一个函数作为参数，这种函数就称之为高阶函数。也就是说用函数作为参数的函数就是高阶函数。

##### map

由于`map()`方法定义在JavaScript的`Array`中，我们调用`Array`的`map()`方法，传入我们自己的函数，就得到了一个新的`Array`作为结果

```
function pow(x) {
    return x * x;
}

var arr = [1, 2, 3, 4, 5, 6, 7, 8, 9];
arr.map(pow); // [1, 4, 9, 16, 25, 36, 49, 64, 81]
```

##### reduce

Array的`reduce()`把一个函数作用在这个`Array`的`[x1, x2, x3...]`上，这个函数必须接收两个参数，`reduce()`把结果继续和序列的下一个元素做累积计算，其效果就是：

```
[x1, x2, x3, x4].reduce(f) = f(f(f(x1, x2), x3), x4)
```

比方说对一个`Array`求和，就可以用`reduce`实现：

```
var arr = [1, 3, 5, 7, 9];
arr.reduce(function (x, y) {
    return x + y;
}); // 25
```

##### filter

filter也是一个常用的操作，它用于把`Array`的某些元素过滤掉，然后返回剩下的元素。

和`map()`类似，`Array`的`filter()`也接收一个函数。和`map()`不同的是，`filter()`把传入的函数依次作用于每个元素，然后根据返回值是`true`还是`false`决定保留还是丢弃该元素。

#### 箭头函数

ES6标准新增了一种新的函数：Arrow Function（箭头函数）。

为什么叫Arrow Function？因为它的定义用的就是一个箭头：

```
x => x * x
```

ES6中定义类的方式

定义一个class

每一个使用class方式定义的类默认都有一个**constructor**函数， 这个函数是构造函数的主函数， 该函数体内部的**this**指向生成的实例，类必须先声明再使用。

定义函数的静态方法

如果定义函数的时候， 大括号内部， **函数名**前声明了**static**， 那么这个函数就为静态函数， 就为静态方法， 和原型没啥关系。

# 定义原型方法

定义**原型方法**，直接这样声明： **函数名 () {} **即可， 小括号内部为参数列表， 大括号内部为代码块。ES5中要定义原型方法是通过：** 构造函数.prototype.原型方法名() {}** , 这种书写形式很繁琐， 使用ES6定义原型的方式有点像java和C#了， 这些都是比较高级语言的特征。

可以通过"对象.属性"或"对象[属性]"这两种方式来访问对象的属性。

创建对象的方法:

1. 构造方法

   ```
   function Rect(w,h){
     this.width=w;
     this.height=h;
     this.area=function(){
        return this.width*this.height;
     }
   }
   //不是直接制造对象
   //通过this来定义成员
   //没有return
   var r=new Rect(5,10)；//这里我们创建对象，上面的函数是构造方法
   ```

2. 原型对象

```
function Person(){}//一个空的函数
Person.prototype.name="mike";
var p=new Person();//创建对象
p.name="jack";//给name属性赋值，则对象自己有了这个属性，如果不赋值，那么p.name就取原型上的值mike值。
```



# ES6

ES6 新增了`let`命令，用来声明变量。它的用法类似于`var`，但是所声明的变量，只在`let`命令所在的代码块内有效。

定义“类”的方法的时候，前面不需要加上`function`这个关键字，直接把函数定义放进去了就可以了。另外，方法之间不需要逗号分隔，加了会报错,类里面有一个`constructor`方法，这就是构造方法。

```
//定义类
class Point {
  constructor(x, y) {
    this.x = x;
    this.y = y;
  }

  toString() {
    return '(' + this.x + ', ' + this.y + ')';
  }
}
```

`constructor`方法是类的默认方法，通过`new`命令生成对象实例时，自动调用该方法。一个类必须有`constructor`方法，如果没有显式定义，一个空的`constructor`方法会被默认添加。



异步操作的模式

1、回调函数

```
function f1(f2) {
  // ...
  //干完上面的事情，在接着干f2
  f2()
}

function f2() {
  // ...
}


```

2、事件监听

```
f1.on('done', f2);

function f1() {
  setTimeout(function () {
    // ...
    f1.trigger('done');
  }, 1000);
}
```

当`f1`发生`done`事件，就执行`f2`

3、发布/订阅

```
jQuery.subscribe('done', f2);
function f1() {
  setTimeout(function () {
    // ...
    jQuery.publish('done');
  }, 1000);
}
jQuery.unsubscribe('done', f2);
```

