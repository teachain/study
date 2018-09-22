##HTML5 Canvas##

HTML5 canvas 元素用户图形的绘制，通过脚本（通常是javascript）来完成，它只是一个图形容器，必须使用脚本来绘制图形。

```
<canvas id="myCanvas" width="200" height="100"
style="border:1px solid #000000;">
</canvas>

```
注意:<font color="red">canvas默认没有边框和内容</font>

canvas 元素本身是没有绘图能力的。所有的绘制工作必须在 JavaScript 内部完成。

例如：

```
<script>
var c=document.getElementById("myCanvas");
var ctx=c.getContext("2d");
ctx.fillStyle="#FF0000";
ctx.fillRect(0,0,150,75);
</script>
```

Canvas 坐标

canvas 是一个二维网格。

<font color="red">canvas 的左上角坐标为 (0,0)</font>

上面的 fillRect 方法拥有参数 (0,0,150,75)。

意思是：在画布上绘制 150x75 的矩形，从左上角开始 (0,0)。





