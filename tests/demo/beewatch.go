package main

import (
	"github.com/beego/beewatch"
)

func main() {
	beewatch.Start(beewatch.LevelTrace)

	appName := "Bee Watch"
	beewatch.Info().Display("App Name", appName)
	beewatch.Info().Break()
	beewatch.Critical().Break()
	beewatch.Trace().Break()
	beewatch.Close()
}
