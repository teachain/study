# 弹性盒模型(flex box)

弹性盒布局中有两个互相垂直的坐标轴：一个称之为主轴（main axis），另外一个称之为交叉轴（cross axis）。主轴并不固定为水平方向的 X 轴，交叉轴也不固定为垂直方向的 Y 轴。在使用时，通过 CSS 属性声明首先定义主轴的方向（水平或垂直），则交叉轴的方向也相应确定下来。容器中的条目可以排列成单行或多行。主轴确定了容器中每一行上条目的排列方向，而交叉轴则确定行本身的排列方向。可以根据不同的页面设计要求来确定合适的主轴方向。有些容器中的条目要求从左到右水平排列，则主轴应该是水平方向的；而另外一些容器中的条目要求从上到下垂直排列，则主轴应该是垂直方向的。

确定主轴和交叉轴的方向之后，还需要确定它们各自的排列方向。对于水平方向上的轴，可以从左到右或从右到左来排列；对于垂直方向上的轴，则可以从上到下或从下到上来排列。对于主轴来说，排列条目的起始和结束位置分别称为主轴起始（main start）和主轴结束（main end）；对于交叉轴来说，排列行的起始和结束位置分别称为交叉轴起始（cross start）和交叉轴结束（cross end）。在容器进行布局时，在每一行中会把其中的条目从主轴起始位置开始，依次排列到主轴结束位置；而当容器中存在多行时，会把每一行从交叉轴起始位置开始，依次排列到交叉轴结束位置。

弹性盒布局中的条目有两个尺寸：主轴尺寸和交叉轴尺寸，分别对应其 DOM 元素在主轴和交叉轴上的大小。如果主轴是水平方向，则主轴尺寸和交叉轴尺寸分别对应于 DOM 元素的宽度和高度；如果主轴是垂直方向，则两个尺寸要反过来。与主轴和交叉轴尺寸对应的是主轴尺寸属性和交叉轴尺寸属性，指的是 CSS 中的属性 width 或 height。比如，当主轴是水平方向时，主轴尺寸属性是 width，而 width 的值是主轴尺寸的大小。

对于弹性盒布局的容器，使用"display: flex"声明使用弹性盒布局。CSS 属性声明"flex-direction"用来确定主轴的方向，从而确定基本的条目排列方式。

##### flex-direction

| 属性值            | 含义                                       |
| -------------- | ---------------------------------------- |
| row（默认值）       | 主轴为水平方向。排列顺序与页面的文档顺序相同。如果文档顺序是 ltr，则排列顺序是从左到右；如果文档顺序是 rtl，则排列顺序是从右到左。 |
| row-reverse    | 主轴为水平方向。排列顺序与页面的文档顺序相反。                  |
| column         | 主轴为垂直方向。排列顺序为从上到下。                       |
| column-reverse | 主轴为垂直方向。排列顺序为从下到上。                       |

默认情况下，弹性盒容器中的条目会尽量占满容器在主轴方向上的一行。当容器的主轴尺寸小于其所有条目总的主轴尺寸时，会出现条目之间互相重叠或超出容器范围的现象。CSS 属性"flex-wrap"用来声明当容器中条目的尺寸超过主轴尺寸时应采取的行为。

属性值含义nowrap（默认值）容器中的条目只占满容器在主轴方向上的一行，可能出现条目互相重叠或超出容器范围的现象。wrap当容器中的条目超出容器在主轴方向上的一行时，会把条目排列到下一行。下一行的位置与交叉轴的方向一致。wrap-reverse与 wrap 的含义类似，不同的是下一行的位置与交叉轴的方向相反。

默认情况下，容器中条目的顺序取决于它们在 HTML 标记代码中的出现顺序。可以通过"order"属性来改变条目在容器中的出现顺序。



## CSS预处理器

- 定义了一种新的专门的编程语言，编译后成正常的CSS文件。为CSS增加一些编程的特性，无需考虑浏览器的兼容问题，让CSS更加简洁，适应性更强，可读性更佳，更易于代码的维护等诸多好处。
- CSS预处理器语言：scss（sass）、LESS等；