all: plugin-go/plugin.so plugin-c/plugin.so
	go test

clean:
	rm -f plugin-go/plugin.so plugin-go/plugin.h
	rm -f plugin-c/plugin.so

plugin-go/plugin.so:
	go build -buildmode=c-shared -o $@ ./plugin-go

plugin-c/plugin.so:
	gcc -fPIC -shared -o $@ ./plugin-c/plugin.c

.PHONY: clean all
