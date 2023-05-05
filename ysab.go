package main

import (
	"os"
	"runtime/pprof"

	"github.com/yunsonbai/ysab/worker"
)

func main() {
	f, _ := os.OpenFile("cpu.prof", os.O_RDWR|os.O_CREATE, 0644)
	pprof.StartCPUProfile(f)
	worker.StartWork()
	defer pprof.StopCPUProfile()
}
