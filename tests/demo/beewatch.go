package main

import (
	"github.com/beego/beewatch"
)

func main() {
	beewatch.Start()

	appName := "Bee Watch"
	beewatch.Info().Display("App Name", appName).Break().
		Printf("Application name is %s.", appName)
	beewatch.Critical().Break()
	beewatch.Trace().Break()
	beewatch.Close()
}
