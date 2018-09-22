##Bootstrap##


####布局容器####

 * <font color="red">.container</font>类用于固定宽度并支持响应式布局的容器。
 * <font color="red">.container-fluid</font>类用于100%宽度，占据全部视口(viewport)的容器.

<font color="#FFA500">注意</font>：这两个容器类不能互相嵌套。


####栅格类####

* 一个行（.row）中最多包含12个列(column),如果多于12个，那么多余部分则整体另起一行来显示。

* 所有的列(column)都必须放在<font color="red">.row</font>内.

* 列的栅格类可以有<font color="red">.col-xl-* </font>、<font color="red">.col-sm- * </font>、 <font color="red">.col-md- * </font>、<font color="red">.col-lg-* </font>  这里的*可取值为1~12之间的值
* .col-md-offset-* 类可以将列向右侧偏移，*表示偏移几格
* 为了使用内置的栅格系统将内容再次嵌套，可以通过添加一个新的.row元素和一系列的.col-sm- *元素到已经存在的.col-sm- * 元素内。


```
  <div class="row">
      <div class="col-md-4">占四个</div>
      <div class="col-md-8">占八格</div>
  </div>
```