

```

package vm

import (
	"fmt"

	"github.com/simplechain-org/go-simplechain/params"
)

func makeStackFunc(pop, push int) stackValidationFunc {
	return func(stack *Stack) error {
		//栈里是否还有pop个元素
		if err := stack.require(pop); err != nil {
			return err
		}
        //能否弹出pop个元素又压入push个元素
		if stack.len()+push-pop > int(params.StackLimit) {
			return fmt.Errorf("stack limit reached %d (%d)", stack.len(), params.StackLimit)
		}
		return nil
	}
}

//复制一个元素，如果它的返回值是nil,表示可以复制
func makeDupStackFunc(n int) stackValidationFunc {
	return makeStackFunc(n, n+1)
}
//交换一个元素，如果它的返回值是nil,表示可以交换
func makeSwapStackFunc(n int) stackValidationFunc {
	return makeStackFunc(n, n)
}
```

