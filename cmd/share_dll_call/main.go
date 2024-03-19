package main

import (
	"fmt"
	"plugin"
)

// 插件接口
type PluginDemo interface {
	Initialize()
	Execute()
}

func main() {
	// // 动态加载插件
	// p, err := plugin.Open("libshare_dll.so")
	// if err != nil {
	// 	fmt.Println("Failed to open plugin:", err)
	// 	return
	// }

	// // 查找并调用插件中的函数
	// sym, err := p.Lookup("MyFunction")
	// if err != nil {
	// 	fmt.Println("Failed to lookup symbol:", err)
	// 	return
	// }

	// // 调用函数
	// startGame, ok := sym.(func())
	// if !ok {
	// 	fmt.Println("Unexpected type from module symbol")
	// 	return
	// }
	// startGame()
	callLocalPlugin()
}

func callLocalPlugin() {
	fmt.Printf("\"init\": %v\n", "init")
	p, err := plugin.Open("libshare_dll.so")
	fmt.Printf("err: %v\n", err)
	s, err := p.Lookup("P")
	if err != nil {
		fmt.Println("can't get plugin instance, err=", err)
		return
	}
	plugin := s.(PluginDemo)
	plugin.Initialize()
	plugin.Execute()
}
