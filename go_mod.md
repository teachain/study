使用要求：

* Go的sdk版本在1.12及以上。
* 在一个非go path的路径中。
* go mod init  moduleName
* 



一些特点：

* 主要是通过GOPATH/pkg/mod下的缓存包来对工程进行构建。

* go mod 更类似于maven这种本地缓存库的管理方式,不论你有多少个工程，只要你引用的依赖的版本是一致的，那么在本地就只会有一份依赖文件的存在

* 在go mod中，增加一条类似这样的替换

  ```
  replace github.com/test/commons v1.1.1 => /Users/test/Workspace/bizgocommons
  ```