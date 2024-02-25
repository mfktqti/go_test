package main

import (
	"context"
	"fmt"
	"time"
)

// func main() {
// 	// 创建一个具有1秒超时的context
// 	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
// 	defer cancel() // 确保在完成操作后取消context

// 	select {
// 	case <-ctx.Done():
// 		fmt.Println("操作超时，ctx.Err():", ctx.Err())
// 		return
// 	case <-time.After(2 * time.Second):
// 		fmt.Println("操作完成")
// 	}
// }

// func main() {
// 	// 创建一个空的Context
// 	ctx := context.Background()

// 	// 调用myFunc，并传入带有超时时间的Context
// 	err := myFunc(ctx)
// 	if err != nil {
// 		fmt.Println("myFunc failed:", err)
// 	}
// }

// func myFunc(ctx context.Context) error {
// 	// 创建一个带有超时时间的Context
// 	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
// 	defer cancel() // 在函数结束时调用cancel释放资源

// 	// 在goroutine中执行一些耗时操作，同时需要监控Context的状态
// 	go func() {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Println("myFunc is cancelled.")
// 		}
// 	}()

// 	// 模拟一个耗时的操作
// 	time.Sleep(time.Second * 10)
// 	fmt.Println("myFunc is com.")
// 	return nil
// }

func myFunc(ctx context.Context) error {
	for {
		select {
		default:
			// 模拟耗时操作
			time.Sleep(1 * time.Second)
			fmt.Println("Working...")
		case <-ctx.Done():
			fmt.Println("Cancelled.")
			return ctx.Err()
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(5 * time.Second)
		cancel() // 在5秒后取消操作
	}()

	err := myFunc(ctx)
	if err != nil {
		fmt.Println("myFunc failed:", err)
	}
}
