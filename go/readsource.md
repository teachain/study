##源码阅读笔记##
io/ioutil包依赖于bytes,io,os,sort,sync包。
我们可以从源码里可以从import部分看出来。

```
import (
	"bytes"
	"io"
	"os"
	"sort"
	"sync"
)

```