# 基于golang实现的k8s container web terminal

```
TerminalSizeQueue 是一个

type TerminalSizeQueue interface {
	// Next returns the new terminal size after the terminal has been resized. It returns nil when
	// monitoring has been stopped.
	Next() *TerminalSize
}

接口
```