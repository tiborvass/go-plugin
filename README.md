# C Shared Libraries as Go plugins

Package plugin allows you to easily define a plugin for your Go application
and have it call out at runtime, to C shared libraries fully or partially
implementing the user-defined plugin.

The advantage of this is that the implementation of the plugin is language-agnostic.

Tested only on 64bit Linux.

## Installation

Go 1.5+ is needed

```
go get -u github.com/tiborvass/go-plugin
```

## Example 1

```Go
package main

import "github.com/tiborvass/go-plugin"

type MyPlugin struct {
	plugin.Plugin
	Hello   func()
	Goodbye func()
}

func main() {
	var p MyPlugin
	if err := plugin.Open(&p, "path/to/shared/lib/plugin"); err != nil {
		panic(err)
	}
	defer p.Close()

	p.Hello()
	p.Goodbye()
}
```

## Example 2

You'll need gcc installed for this example to work.

```
$ cd $GOPATH/src/github.com/tiborvass/go-plugin
$ make clean
$ make
go build -buildmode=c-shared -o plugin-go/plugin.so ./plugin-go
gcc -fPIC -shared -o plugin-c/plugin.so ./plugin-c/plugin.c
go test
Hello from main
Hello from Go plugin
This is a function implemented only in Go
Hello from C plugin!
This is a function implemented only in C
Goodbye from main
PASS
ok  	plugin	0.034s
```
