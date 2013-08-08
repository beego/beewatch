package main

import (
	"github.com/beego/beewatch"
)

func main() {
	i := 1
	beewatch.Start(beewatch.Trace)
	beewatch.Display(beewatch.Info, "i", i)
	beewatch.Close()
}
