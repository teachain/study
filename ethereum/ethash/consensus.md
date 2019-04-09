

代码段阅读

```
// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications.
func (ethash *Ethash) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	// If we're running a full engine faking, accept any input as valid
	if ethash.config.PowMode == ModeFullFake || len(headers) == 0 {
		abort, results := make(chan struct{}), make(chan error, len(headers))
		for i := 0; i < len(headers); i++ {
			results <- nil
		}
		return abort, results
	}

	// Spawn as many workers as allowed threads
	//获取当前最大的并行数作为工人的数量
	workers := runtime.GOMAXPROCS(0)
	//如果要处理的区块头数目小于最大的并行数，那么就将工人的数量设置为区块头数即可
	//也就是不要雇佣那么多工人（工作就那么一点）
	if len(headers) < workers {
		workers = len(headers)
	}

	// Create a task channel and spawn the verifiers
	var (
		inputs = make(chan int)
		done   = make(chan int, workers)
		errors = make([]error, len(headers))
		abort  = make(chan struct{})
	)
	//让固定的几个工人都进入准备状态中
	for i := 0; i < workers; i++ {
		go func() {
			//只要他从inputs能够取到任务，他就干，干完了就通知done
			for index := range inputs {
				//我们知道的是index就是任务的下标
				errors[index] = ethash.verifyHeaderWorker(chain, headers, seals, index)
				//干完了就通知done
				done <- index
			}
		}()
	}

	errorsOut := make(chan error, len(headers))
	go func() {
		defer close(inputs)
		var (
			in, out = 0, 0
			checked = make([]bool, len(headers))
			inputs  = inputs
		)
		for {
			select {
			case inputs <- in:
				//这里就是in++
				//然后是if in==len(headers) {inputs = nil}的意思
				if in++; in == len(headers) {
					// Reached end of headers. Stop sending to workers.
					//将这个inputs参数设置为nil,也就是worker那边就会退出
					//这里也可以用clone(inputs)来达到目的。
					inputs = nil
				}
			case index := <-done:
				//这里的技巧就是index为干完的工作的下标，我们必须马上设置为true,表示它已经完成了
				//是否将它放入到errorsOut中，则需要看checked[out]是否为true，如果为true
				//就会将errors[out]放入errorsOut,并且out就会自增，
				//所以即使index不是按着顺序来的，也没有关系
				//由out来保证errorsOut中得到的绝对是按顺序来处理的结果
				for checked[index] = true; checked[out]; out++ {
					errorsOut <- errors[out]
					if out == len(headers)-1 {
						return
					}
				}
			case <-abort:
				return
			}
		}
	}()
	return abort, errorsOut
}
```

