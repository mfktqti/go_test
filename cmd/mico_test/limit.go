package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func limit_test() {
	//说明

	/*
				type Limiter struct {
				 mu     sync.Mutex
				 limit  Limit
				 burst  int // 令牌桶的大小
				 tokens float64
				 last time.Time // 上次更新tokens的时间
				 lastEvent time.Time // 上次发生限速器事件的时间（通过或者限制都是限速器事件）
				}

				其主要字段的作用是：

		limit：limit字段表示往桶里放Token的速率，它的类型是Limit，是int64的类型别名。设置limit时既可以用数字指定每秒向桶中放多少个Token，也可以指定向桶中放Token的时间间隔，其实指定了每秒放Token的个数后就能计算出放每个Token的时间间隔了。
		burst: 令牌桶的大小。
		tokens: 桶中的令牌。
		last: 上次往桶中放 Token 的时间。
		lastEvent：上次发生限速器事件的时间（通过或者限制都是限速器事件）
		可以看到在 timer/rate 的限流器实现中，并没有单独维护一个 Timer 和队列去真的每隔一段时间向桶中放令牌，而是仅仅通过计数的方式表示桶中剩余的令牌。每次消费取 Token 之前会先根据上次更新令牌数的时间差更新桶中Token数。




	*/

	//使用
	/*
	   我们可以使用以下方法构造一个限流器对象：

	   limiter := rate.NewLimiter(10, 100);

	   这里有两个参数：

	   第一个参数是 r Limit，设置的是限流器Limiter的limit字段，代表每秒可以向 Token 桶中产生多少 token。Limit 实际上是 float64 的别名。
	   第二个参数是 b int，b 代表 Token 桶的容量大小，也就是设置的限流器 Limiter 的burst字段。
	   那么，对于以上例子来说，其构造出的限流器的令牌桶大小为 100, 以每秒 10 个 Token 的速率向桶中放置 Token。

	   除了给r Limit参数直接指定每秒产生的 Token 个数外，还可以用 Every 方法来指定向桶中放置 Token 的间隔，例如：

	   limit := rate.Every(100 * time.Millisecond);
	   limiter := rate.NewLimiter(limit, 100);
	   以上就表示每 100ms 往桶中放一个 Token。本质上也是一秒钟往桶里放 10 个。

	*/

	// 每秒钟最多允许一个请求
	limiter := rate.NewLimiter(1, 2)
	for i := 0; i < 5; i++ {
		if limiter.Allow() {
			fmt.Println("Allow request at", time.Now())
		} else {
			fmt.Println("Rate limited at", time.Now())
		}
		time.Sleep(time.Second / 2)
	}

	fmt.Println("Rate limited at----------------")
	limit := rate.Every(1000 * time.Millisecond)
	limiter = rate.NewLimiter(limit, 2)
	for i := 0; i < 10; i++ {
		if limiter.Allow() {
			fmt.Println("Allow request at", time.Now())
		} else {
			fmt.Println("Rate limited at", time.Now())
		}
		time.Sleep(time.Second / 2)
	}
	time.Sleep(10 * time.Second)

	// 一直等到获取到桶中的令牌
	err := limiter.WaitN(context.Background(), 10)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	//// 设置一秒的等待超时时间
	ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
	err = limiter.Wait(ctx)
	if err != nil {
		fmt.Println("Error: ", err)
	}

}
