```
package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/log"
)

const (
	jsonrpcVersion           = "2.0"
	serviceMethodSeparator   = "_"  //这就是分隔符

//rpc 接口的方法名为:Namespace+serviceMethodSeparator+Service的方法（第一个字母要改为小写）
例如
{
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
},

func (api *PrivateMinerAPI) Start(threads *int) error {....}

那么rpc的方法名就是 miner_start

// parseRequest will parse a single request from the given RawMessage. It will return
// the parsed request, an indication if the request was a batch or an error when
// the request could not be parsed.
func parseRequest(incomingMsg json.RawMessage) ([]rpcRequest, bool, Error) {
	var in jsonRequest
	if err := json.Unmarshal(incomingMsg, &in); err != nil {
		return nil, false, &invalidMessageError{err.Error()}
	}

	if err := checkReqId(in.Id); err != nil {
		return nil, false, &invalidMessageError{err.Error()}
	}

	// subscribe are special, they will always use `subscribeMethod` as first param in the payload
	if strings.HasSuffix(in.Method, subscribeMethodSuffix) {
		reqs := []rpcRequest{{id: &in.Id, isPubSub: true}}
		if len(in.Payload) > 0 {
			// first param must be subscription name
			var subscribeMethod [1]string
			if err := json.Unmarshal(in.Payload, &subscribeMethod); err != nil {
				log.Debug(fmt.Sprintf("Unable to parse subscription method: %v\n", err))
				return nil, false, &invalidRequestError{"Unable to parse subscription request"}
			}

			reqs[0].service, reqs[0].method = strings.TrimSuffix(in.Method, subscribeMethodSuffix), subscribeMethod[0]
			reqs[0].params = in.Payload
			return reqs, false, nil
		}
		return nil, false, &invalidRequestError{"Unable to parse subscription request"}
	}

	if strings.HasSuffix(in.Method, unsubscribeMethodSuffix) {
		return []rpcRequest{{id: &in.Id, isPubSub: true,
			method: in.Method, params: in.Payload}}, false, nil
	}
    //在这里解析分开方法和服务。 
	elems := strings.Split(in.Method, serviceMethodSeparator)
	
	if len(elems) != 2 {
		return nil, false, &methodNotFoundError{in.Method, ""}
	}

	// regular RPC call
	if len(in.Payload) == 0 {
		return []rpcRequest{{service: elems[0], method: elems[1], id: &in.Id}}, false, nil
	}

	return []rpcRequest{{service: elems[0], method: elems[1], id: &in.Id, params: in.Payload}}, false, nil
}
```

