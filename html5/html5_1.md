<font color="red">在 HTML5 中，拖放是标准的一部分，任何元素都能够拖放。</font>

首先，为了使元素可拖动，把 draggable 属性设置为 true 


## HTML5 (视频)##


`<video>` 元素提供了 播放、暂停和音量控件来控制视频。

同时`<video>` 元素元素也提供了 width 和 height 属性控制视频的尺寸.如果设置的高度和宽度，所需的视频空间会在页面加载时保留。。如果没有设置这些属性，浏览器不知道大小的视频，浏览器就不能再加载时保留特定的空间，页面就会根据原始视频的大小而改变。

`<video>` 与`</video>` 标签之间插入的内容是提供给不支持 video 元素的浏览器显示的。

`<video>` 元素支持多个 `<source>` 元素. 

`<source>` 元素可以链接不同的视频文件。浏览器将使用第一个可识别的格式：

`<video>` 元素支持三种视频格式： MP4, WebM, 和 Ogg

目前MP4是最佳的选择

```

<video width="320" height="240" controls>
  <source src="movie.mp4" type="video/mp4">
  <source src="movie.ogg" type="video/ogg">
  <source src="movie.webm" type="video/webm">
您的浏览器不支持Video标签。
</video>
```

##HTML5 Audio##

```
<audio controls>
  <source src="horse.ogg" type="audio/ogg">
  <source src="horse.mp3" type="audio/mpeg">
  <source src="horse.wav" type="audio/wav">
您的浏览器不支持 audio 元素。
</audio>
```

目前最佳的选择是mp3


control 属性供添加播放、暂停和音量控件。

在`<audio>` 与 `</audio>` 之间你需要插入浏览器不支持的`<audio>`元素的提示文本 。

`<audio>` 元素允许使用多个 `<source>` 元素.

`<source>` 元素可以链接不同的音频文件，浏览器将使用第一个支持的音频文件


##HTML5 Web 存储##

localStorage和sessionStorage

客户端存储数据的两个对象为：

* localStorage 没有时间限制的数据存储
* sessionStorage 针对一个session的数据存储

在使用 web 存储前,应检查浏览器是否支持 localStorage 和sessionStorage:


```
if(typeof(Storage)!=="undefined")
{
    // 是的! 支持 localStorage  sessionStorage 对象!
    // 一些代码.....
} else {
    // 抱歉! 不支持 web 存储。
}

```

##localStorage 对象##

localStorage 对象存储的数据没有时间限制。第二天、第二周或下一年之后，数据依然可用。

```
localStorage.lastname="Smith";
document.getElementById("result").innerHTML="Last name: "
+ localStorage.lastname;

```

##sessionStorage 对象##
sessionStorage 方法针对一个 session 进行数据存储。当用户关闭浏览器窗口后，数据会被删除。

```
if (sessionStorage.clickcount)
{
    sessionStorage.clickcount=Number(sessionStorage.clickcount)+1;
}
else
{
    sessionStorage.clickcount=1;
}
document.getElementById("result").innerHTML="在这个会话中你已经点击了该按钮 " + sessionStorage.clickcount + " 次 ";

```


