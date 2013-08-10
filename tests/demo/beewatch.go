package main

import (
	"time"

	"github.com/beego/beewatch"
)

func main() {
	beewatch.Start()

	appName := "Bee Watch"
	beewatch.Info().Display("App Name", appName).Break().
		Printf("Application name is %s.", appName)
	beewatch.Critical().Break()
	beewatch.Trace().Break()

	for i := 0; i < 3; i++ {
		go multipletest(i)
	}

	beewatch.Trace().Printf("Wait 3 seconds...")
	select {
	case <-time.After(3 * time.Second):
		beewatch.Trace().Printf("Done debug")
	}

	beewatch.Close()
}

func multipletest(num int) {
	beewatch.Critical().Break().Display("num", num)
}

// http://icalialabs.github.io/furatto/javascript.html
