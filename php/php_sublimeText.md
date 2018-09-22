##sublime text3 配置##

在sublime text3的菜单栏里选择Tools->Build System->new Build System...

在新建的文件里，把原有的内容删掉，然后将以下内容拷贝过去

```
   {
		"cmd": ["php","$file"],
		"file_regex":"php$",
		"selector":"source.php"
   }
```
然后将文件内容保存为php.sublime-build

重启sublime text3 ,然后就可以使用ctrl+b来执行php文件了。