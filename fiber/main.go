package main

import (
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	InitSlog()

	Serve("localhost:8421")
}
