测试的基础



<table>
<thead>
<tr>
  <th>关键点</th>
  <th>说明</th>
</tr>
</thead>
<tbody><tr>
  <td>导入需要的包</td>
  <td>import testing (如果你是goland,那么可以忽略，因为ide就自动帮你加上)</td>
</tr>
<tr>
  <td>执行命令</td>
  <td>go test file_test.go</td>
</tr>
<tr>
  <td>测试文件命名</td>
  <td>必须以_test.go结尾</td>
</tr>
<tr>
  <td>功能测试的用力函数</td>
  <td>必须以Test开头&amp;&amp;后头跟着的函数名不能以小写字母开头，比如：Testcbs 就是不行的，TestCbs就是ok的</td>
</tr>
<tr>
  <td>功能测试参数</td>
  <td>testing.T</td>
</tr>
<tr>
  <td>压力测试用例函数</td>
  <td>必须以Benchmark开头&amp;&amp;其后的函数名不能以小写字母开头(例子同上)</td>
</tr>
<tr>
  <td>压力测试参数</td>
  <td>testing.B</td>
</tr>
<tr>
  <td>测试信息</td>
  <td>.Log方法，默认情况下是不会显示的，只有在go test -v的时候显示</td>
</tr>
<tr>
  <td>测试控制</td>
  <td>通过Error/Errorf/FailNow/Fatal等来进行测试是否是失败，或者在失败的情况下的控制</td>
</tr>
<tr>
  <td>压力测试命令</td>
  <td>go test -test.bench file_test.go</td>
</tr>
<tr>
  <td>压力测试的循环体</td>
  <td>使用test.B.N</td>
</tr>
</tbody></table>



- 功能测试 Test开头
- 压力(性能,基准)测试Benchmark开头



```
go test      #功能测试
go test -bench=.  #压力测试命令 ,使用go help testflag命令来查看怎么使用
go test -cover  #代码的覆盖率测试
go test -v -test.run TestRefreshAccessToken #测试单个方法
go test -v  wechat_test.go wechat.go  #带上测试文件 特定的
```

