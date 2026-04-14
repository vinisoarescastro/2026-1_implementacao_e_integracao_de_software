package main

import (
	"fmt"
	"runtime"
)

// version é injetada em tempo de build via ldflags.
var version = "dev"

func main() {
	fmt.Printf("simulador %s %s/%s — em construção\n", version, runtime.GOOS, runtime.GOARCH)
}
