

```
package vm

import (
	"fmt"
	"math/big"
)

// Stack is an object for basic stack operations. Items popped to the stack are
// expected to be changed and modified. stack does not take care of adding newly
// initialised objects.

//栈
type Stack struct {
	data []*big.Int
}

//创建一个栈（虽然设置了cap为1024,但其实它是无限栈--除非没有内存了）
func newstack() *Stack {
	return &Stack{data: make([]*big.Int, 0, 1024)}
}

//返回栈中所有的数据
// Data returns the underlying big.Int array.
func (st *Stack) Data() []*big.Int {
	return st.data
}
//将元素压入栈
func (st *Stack) push(d *big.Int) {
	// NOTE push limit (1024) is checked in baseCheck
	//stackItem := new(big.Int).Set(d)
	//st.data = append(st.data, stackItem)
	st.data = append(st.data, d)
}
//将多个元素压入栈
func (st *Stack) pushN(ds ...*big.Int) {
	st.data = append(st.data, ds...)
}

//将栈顶元素弹出栈
func (st *Stack) pop() (ret *big.Int) {
	//我们要时刻记得，栈只有一端可以操作，那就是栈顶
	ret = st.data[len(st.data)-1]
	//移除栈顶，让下一个元素成为栈顶
	st.data = st.data[:len(st.data)-1]
	return
}

//返回栈中元素个数
func (st *Stack) len() int {
	return len(st.data)
}

//将索引n元素与栈顶元素进行交换
//从st.len()-n我们知道，它其实就是弹出栈顶n次的那个元素
//st.len()为栈中有效的元素个数
func (st *Stack) swap(n int) {
	st.data[st.len()-n], st.data[st.len()-1] = st.data[st.len()-1], st.data[st.len()-n]
}

//复制弹出栈顶n次的那个元素压入到栈顶
func (st *Stack) dup(pool *intPool, n int) {
	st.push(pool.get().Set(st.data[st.len()-n]))
}

//查看栈顶元素
func (st *Stack) peek() *big.Int {
	return st.data[st.len()-1]
}

//从栈顶开始往栈底方向数数，从0开始数，第n个
// Back returns the n'th item in stack
func (st *Stack) Back(n int) *big.Int {
	return st.data[st.len()-n-1]
}

//意思是能否弹出n个元素
//也即是要求栈中有n个元素
func (st *Stack) require(n int) error {
	if st.len() < n {
		return fmt.Errorf("stack underflow (%d <=> %d)", len(st.data), n)
	}
	return nil
}

//打印栈中的元素
// Print dumps the content of the stack
func (st *Stack) Print() {
	fmt.Println("### stack ###")
	if len(st.data) > 0 {
		for i, val := range st.data {
			fmt.Printf("%-3d  %v\n", i, val)
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("#############")
}
```

