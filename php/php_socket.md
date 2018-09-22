##socket编程##

<font color="red">注意:网络字节顺序采用big endian排序方式。</font>

```
string pack ( string $format [, mixed $args [, mixed $... ]] )
```

Pack given arguments into a binary string according to format.

例如我们需要发送这样一个数据包

包结构是这样子的:

字段  | 字节数 | 说明
------------- | ------------- | ---
包头  | 定长  | 每一个通信消息必须包含的内容
包体  | 不定长  | 根据每个通信消息的不同产生变化

其中包头详细内容如下:

字段  | 字节数 | 类型 | 说明
------------- | ------------- | --- | ----
pkg_len  | 2  | ushort | 整个包的长度，不超过4k
version | 1  | uchar | 通讯协议版本号
command_id  | 1  | uchar | 消息命令ID
result  | 1  | uchar | 请求时不起作用，请求返回时使用

例如一个登陆数据包的body如下

字段 | 字节数 | 类型 | 说明
---- | -----| ---- | ---
用户名 | 30 | uchar[30] | 登陆用户名
密码 | 32 | uchar[32] | 登陆密码



包头是定长的，通过计算可知包头占7个字节，并且包头在包体之前。比如用户陈一回需要登录，密码是123456，则代码如下：

```
<?php
$version    = 1;
$result     = 0;
$command_id = 1001;
$username   = "陈一回";
$password   = md5("123456");
// 构造包体
$bin_body   = pack("a30a32", $username, $password);
// 包体长度
$body_len   = strlen($bin_body);
$bin_head   = pack("nCns", $body_len, $version, $command_id, $result);
$bin_data   = $bin_head . $bin_body;
// 发送数据
socket_write($socket, $bin_data, strlen($bin_data));
socket_close($socket);
```