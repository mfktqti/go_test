package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	fmt.Println("Hello World")
	fmt.Fprintln(os.Stdout, "Hello World") //直接输出

	var b bytes.Buffer
	fmt.Fprintln(&b, "Hello World") //将字符串写入b变量
	fmt.Println(b.String())

	//同时写给多个writer
	var foo, bar bytes.Buffer
	mw := io.MultiWriter(&foo, &bar) // 多写入器
	fmt.Fprintln(mw, "hello world")
	fmt.Println(foo.String())
	fmt.Println(bar.String())

	r := strings.NewReader("Hello world")
	b2, _ := io.ReadAll(r)
	fmt.Printf("b2: %v\n", string(b2))

	//一次性从多个 reader 上读取数据
	r2 := strings.NewReader("hello world")
	r3 := strings.NewReader("hello world")
	mr := io.MultiReader(r2, r3)
	b3, _ := io.ReadAll(mr)
	fmt.Printf("b3: %v\n", string(b3))

	//reader 将数据推送给 writer
	// Create a reader
	rr := strings.NewReader("Hello World")
	// Create a writer
	var b4 bytes.Buffer
	// Push data
	rr.WriteTo(&b4) // Don't forget &
	// Optional: verify data
	fmt.Println(b4.String())

	//writer 从 reader 中拉出数据
	r4 := strings.NewReader("hello world")
	var b5 bytes.Buffer
	b5.ReadFrom(r4)
	fmt.Println("b5", b5.String())

	//使用 io.Copy
	r5 := strings.NewReader("hello world")
	var b6 bytes.Buffer
	_, _ = io.Copy(&b6, r5)
	fmt.Printf("b6.String(): %v\n", b6.String())

	io_pipe()
	io_pipe_test()
}

func io_pipe() {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		fmt.Fprintln(pw, "Hello world")
	}()

	b, _ := io.ReadAll(pr)
	fmt.Printf("string(b): %v\n", string(b))
}

// 用io.Pipe、io.Copy和io.MultiWriter捕捉函数的stdout到一个变量中。
func io_pipe_test() {
	fmt.Println("io_pipe_test-------")
	pr, pw := io.Pipe()
	go func(w *io.PipeWriter) {
		defer w.Close()
		fmt.Fprintln(w, "Hello world")
	}(pw)

	var b bytes.Buffer

	mw := io.MultiWriter(os.Stdout, &b)
	_, _ = io.Copy(mw, pr)
	fmt.Printf("b.String(): %v\n", b.String())

}
