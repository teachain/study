# RLP数据定义

RLP编码的定义只处理以下两类数据：

- 字符串（string）是指字节数组。例如，空串""，再如单词"cat"，以及句子"Lorem ipsum dolor sit amet, consectetur adipisicing elit"等。
- 列表（list）是一个可嵌套结构，里面可包含字符串和列表。例如，空列表[]，再如一个包含两个字符串的列表["cat","dog"]，在比如嵌套列表的复杂列表["cat", ["puppy", "cow"], "horse", [[]], "pig",  [""], "sheep"]。

其他类型的数据需要转成以上的两类数据，才能编码。转换的规则RLP编码不统一规定，可以自定义转换规则。例如struct可以转成列表，int可以转成二进制序列（属于字符串这一类, 必须去掉首部0，必须用大端模式表示）。