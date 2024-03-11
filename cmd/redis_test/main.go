package main

import (
	"fmt"
	"go-test/internal/config"
	"go-test/internal/db/redis"
	"go-test/internal/log"
	"time"
)

func main() {

	c, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log.InitLogger(c.Log)

	if err := redis.ConnectToRedis(c.Redis); err != nil {
		panic(err)
	}

	//string

	// fmt.Println(redis.Get("t"))
	// fmt.Printf("redis.Set(\"t\", time.Now().Format(time.DateOnly)): %v\n", redis.Set("t", time.Now().Format(time.DateOnly)))

	// fmt.Println(redis.Get("t"))
	// fmt.Printf("redis.Exists(\"t\"): %v\n", redis.Exists("t1"))
	// fmt.Printf("redis.Del(\"t\"): %v\n", redis.Del("t"))
	// // fmt.Printf("redis.Set(\"incr-1\", \"lakjsdl\"): %v\n", redis.Set("incr-1", "q"))
	// fmt.Println(redis.Incr("incr-1"))
	// fmt.Println(redis.Incr("incr-1"))
	// fmt.Println(redis.Expire("incr-1", 2*time.Second))
	// time.Sleep(1 * time.Second)
	// fmt.Println(redis.Incr("incr-1"))
	// fmt.Println(redis.Incr("incr-1"))
	// time.Sleep(1 * time.Second)
	// fmt.Println(redis.Incr("incr-1"))
	// fmt.Println(redis.SetNX("incr-11", 10, time.Second))

	//hash

	// fmt.Printf("redis.HSet(\"h-1\", \"h-1-1\", \"h-1-1-1\"): %v\n", redis.HSet("h-1", "h-1-1", "h-1-1-1"))
	// fmt.Printf("redis.HSet(\"h-1\", \"h-2-1\", \"h-2-1-1\"): %v\n", redis.HSet("h-1", "h-2-1", "h-2-1-1"))
	// fmt.Println(redis.HGet("h-1", "h-1-1"))
	// fmt.Println(redis.HGetAll("h-1"))
	// fmt.Println(redis.HKeys("h-1"))

	//list
	fmt.Printf("redis.LPush(\"list1\", 1, 2, 3): %v\n", redis.LPush("list1", 1, 2, 3))
	fmt.Printf("redis.LPush(\"list1\", 1, 2, 3): %v\n", redis.LPush("list1", 1, 2, 3))
	fmt.Println(redis.LLen("list1"))
	fmt.Println(redis.LRange("list1", 1, 2))

	fmt.Println(redis.ZAdd("zset1", "111", 12))
	fmt.Println(redis.ZAdd("zset1", "112", 13))
	fmt.Println(redis.ZAdd("zset1", "110", 10))
	fmt.Println(redis.ZRangeByScore("zset1", 12, 13))
	fmt.Println(redis.ZRevRangeWithScores("zset1", 1, 2))
	fmt.Println(redis.ZScore("zset1", "110"))
	fmt.Println(redis.ZRank("zset1", "111"))
	fmt.Printf("\"-----------\": %v\n", "-----------")
	// for i := 1; i <= 10; i++ {
	// 	go func(i int) {
	redisLock := redis.NewDistributedLock()
	err = redisLock.Lock("my-lock")
	if err != nil {
		fmt.Printf("Lock err: %v\n", err)
		return
	}
	fmt.Printf("\"lock\": %v\n", "go")
	time.Sleep(15 * time.Second)
	err = redisLock.Unlock()
	if err != nil {
		fmt.Printf("Unlock err: %v\n", err)
		return
	}
	// fmt.Printf("i: %v\n", i)

	// 	}(i)
	// }

	// time.Sleep(20 * time.Second)
	fmt.Printf("\"go\": %v\n", "go")

}
