package main

import (
	"github.com/0xSplits/specta/cmd"
	"github.com/xh3b4sd/tracer"
)

func main() {
	err := cmd.New().Execute()
	if err != nil {
		tracer.Panic(tracer.Mask(err))
	}
}
