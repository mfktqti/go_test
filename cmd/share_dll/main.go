package main

import "fmt"

//执行下面的命令
// CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ../share_dll_call/libshare_dll.so main.go

// type PluginDemo interface {
// 	InitiaiExecute()
// 	Execute()
// }

type SimplePlugin struct{}

func (s *SimplePlugin) Initialize() {
	fmt.Println("initialize now..")
}

func (s *SimplePlugin) Execute() {
	fmt.Println("Execute Plugin now...")
}

func init() {
	fmt.Println("Plugin init..")
}

var P SimplePlugin
