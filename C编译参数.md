## 文件处理

在 Compile Sources 设置某个文件的 Compiler Flags

得到的 warning 信息是 `-Wdeprecated-declarations` 需要改成 `-Wno-deprecated-declarations`

即所有的这类报错都是 `-W错误信息` 的格式，我们需要将 `-W` 替换成 `-Wno-` 即告诉编译器这个文件的这个错误不在提示警告，可以添加多个