## beego的学习

从2015年5月开始接触golang,经过一段时间的学习，感觉这门语言挺好玩的，就继续深入了，到现在已经一年有余，以前一直是用它来写tcp的应用，绝少写http方面的，现在终于不怎么忙了，有空将beego学习一下，也补一下用golang写http方面的短板。

## Session控制

beego中使用session相当方便，可通过两种方式启用

1. 在main入口函数中设置
    
    beego.SessionOn=true
    
2. 在配置文件中配置

   sessionon=true
   
   
## 多数据格式输出

    func(this *MainController){
        mystruct:={...}
        this.Data["json"]=&mystruct
        this.ServeJSON()
    }
这里直接输出json格式的数据，这里的this.Data["json"]是固定的，调用ServeJSON()之后，会设置content-type为application/json,然后同时把数据进行json序列化输出，也就是可以支持以下三种格式的数据输出：

* json----->application/json

```
this.Data["json"]=&mystruct
    
this.ServeJSON() 
```


* xml------>application/xml

``` 
this.Data["xml"]=&mystruct
    
this.ServeXML()
```

* jsonp---->application/javascript

```    
this.Data["jsonp"]=&mystruct
    
this.ServeJSONP()
```
在源码beego/context/output.go里可以看到，是分别设置了content-type的。




