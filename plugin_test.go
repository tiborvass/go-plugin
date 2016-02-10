package plugin

import (
	"fmt"
	"testing"
)

func ExampleOpen() {
	type MyPlugin struct {
		Plugin
		HelloWorld func()
		OnlyInC    func()
		OnlyInGo   func()
	}

	fmt.Println("Hello from main")
	defer fmt.Println("Goodbye from main")

	for _, path := range []string{"./plugin-go/plugin", "./plugin-c/plugin"} {
		var myPlugin MyPlugin
		if err := Open(&myPlugin, path); err != nil {
			panic(err)
		}
		defer myPlugin.Close()
		myPlugin.HelloWorld()
		myPlugin.OnlyInGo()
		myPlugin.OnlyInC()
	}
}

func TestOpen(t *testing.T) {
	ExampleOpen()
}
