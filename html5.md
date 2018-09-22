## 在客户端存储数据

HTML5 提供了两种在客户端存储数据的新方法：

- localStorage - 没有时间限制的数据存储
- sessionStorage - 针对一个 session 的数据存储



localStorage 方法存储的数据没有时间限制。第二天、第二周或下一年之后，数据依然可用。

如何创建和访问 localStorage

localStorage是内置对象，开箱即用。

目前所有的浏览器中都会把localStorage的值类型限定为string类型



localStorage与sessionStorage的唯一一点区别就是localStorage属于永久性存储，而sessionStorage属于当会话结束的时候，sessionStorage中的键值对会被清空

localStorage常用的API

| **名称**             | **作用**                                             |
| -------------------- | ---------------------------------------------------- |
| setItem              | 存储数据【增】                                       |
| getItem              | 读取数据【查单个】                                   |
| removeItem           | 删除某个数据【删单个】                               |
| clear                | 删除全部数据【删全部】                               |
| length               | localStorage存储变量的个数【计算数据总数】           |
| key                  | 读取第i个数据的名字或称为键值(从0开始计数)           |
| valueOf              | 获取所有存储的数据【查全部】                         |
| hasOwnProperty       | 检查localStorage上是否保存了变量x，需要传入x【判断】 |
| propertyIsEnumerable | 用来检测属性是否属于某个对象的【判断】               |
| toLocaleString       | 将（数组）转为本地字符串                             |