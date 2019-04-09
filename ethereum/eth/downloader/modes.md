源码解读



从eth/handler中的这一行代码中，我们可以得到结论

```
if mode == downloader.FastSync && blockchain.CurrentBlock().NumberU64() > 0 {
		log.Warn("Blockchain not empty, fast sync disabled")
		mode = downloader.FullSync
	}
```

当同步的模式是快速同步但当前链的高度已经大于0时，也就是不仅仅是只存在创世区块了，这个时候同步模式会被强制设置为全节点同步模式。也即是说当模式还是快速同步时，它的高度就不可能大于0。

```
package downloader

import "fmt"

// SyncMode represents the synchronisation mode of the downloader.
type SyncMode int

const (
  //全节点同步
	FullSync  SyncMode = iota // Synchronise the entire blockchain history from full blocks
	//快速同步，先同步区块头，后同步区块
	FastSync                  // Quickly download the headers, full sync only at the chain head
	//只同步区块头
	LightSync                 // Download only the headers and terminate afterwards
)

func (mode SyncMode) IsValid() bool {
	return mode >= FullSync && mode <= LightSync
}

// String implements the stringer interface.
func (mode SyncMode) String() string {
	switch mode {
	case FullSync:
		return "full"
	case FastSync:
		return "fast"
	case LightSync:
		return "light"
	default:
		return "unknown"
	}
}

func (mode SyncMode) MarshalText() ([]byte, error) {
	switch mode {
	case FullSync:
		return []byte("full"), nil
	case FastSync:
		return []byte("fast"), nil
	case LightSync:
		return []byte("light"), nil
	default:
		return nil, fmt.Errorf("unknown sync mode %d", mode)
	}
}

func (mode *SyncMode) UnmarshalText(text []byte) error {
	switch string(text) {
	case "full":
		*mode = FullSync
	case "fast":
		*mode = FastSync
	case "light":
		*mode = LightSync
	default:
		return fmt.Errorf(`unknown sync mode %q, want "full", "fast" or "light"`, text)
	}
	return nil
}
```

