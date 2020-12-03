package safeclose

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	deferFunc []func()
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
)

func init() {
	deferFunc = []func(){}
	ctx, cancel = context.WithCancel(context.Background())
	go signalClose()
}

// 异步执行 f 函数，未执行完之前 Wait 函数会一直阻塞。
func Do(f func()) {
	wg.Add(1)
	go func() {
		f()
		wg.Done()
	}()
}

// 异步执行 f 函数，未执行完之前 Wait 函数会一直阻塞。f 一般会是阻塞的任务，f 的参数可以接收信号用来终止 f 函数。
func DoContext(f func(context.Context)) {
	wg.Add(1)
	go func() {
		f(ctx)
		wg.Done()
	}()
}

// 这里的 f 函数不会立即执行，它会在 Do(f) 和 DoContext(f) 的 f 都结束之后执行。
func Defer(f func()) {
	deferFunc = append(deferFunc, f)
}

// 等待所有异步函数结束后，执行 Defer(f) 注册的函数。
func Wait() {
	wg.Wait()
	for _, f := range deferFunc {
		Do(f)
	}
	wg.Wait()
}

// 发送终止信号，DoContext(func(ctx context.Context){}) 的 ctx 可以接收到这个信号。
func Cancel() {
	cancel()
}

// 收到 SIGINT 或 SIGTERM 时，发送终止信号，DoContext(func(ctx context.Context){}) 的 ctx 可以接收到这个信号。
func signalClose() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT)
	<-s
	cancel()
}
