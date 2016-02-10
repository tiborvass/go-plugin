package main

import "fmt"

import "C"

//export HelloWorld
func HelloWorld() {
	fmt.Println("Hello from Go plugin")
}

//export OnlyInGo
func OnlyInGo() {
	fmt.Println("This is a function implemented only in Go")
}

func main() {
}
