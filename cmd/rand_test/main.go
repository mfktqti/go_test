package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	for i := 0; i < 5; i++ {
		fmt.Printf("r.Intn(100): %v\n", r.Intn(100))
	}
	fmt.Printf("\"xxx\": %v\n", "xxx")
}
