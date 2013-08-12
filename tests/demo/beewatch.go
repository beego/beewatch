package main

import (
	"time"

	"github.com/beego/beewatch"
)

func main() {
	// Start with default level: Trace,
	// or use `beewatch.Start(beewatch.Info())` to set the level.
	beewatch.Start()

	// Some variables.
	appName := "Bee Watch"
	boolean := true
	number := 3862
	floatNum := 3.1415926
	test := "fail to watch"

	// Add variables to be watched, must be even sized.
	// Note that you have to pass the pointer.
	// For now it only support basic types.
	// In this case, 'test' will not be watched.
	beewatch.AddWatchVars("test", test, "App Name", &appName,
		"bool", &boolean, "number", &number, "float", &floatNum)

	// Display something.
	beewatch.Info().Display("App Name", appName).Break().
		Printf("boolean value is %v; number is %d", boolean, number)

	beewatch.Critical().Break().Display("float", floatNum)

	// Change some values.
	appName = "Bee Watch2"
	number = 250
	// Here you will see something interesting.
	beewatch.Trace().Break()

	// Multiple goroutines test.
	for i := 0; i < 3; i++ {
		go multipletest(i)
	}

	beewatch.Trace().Printf("Wait 3 seconds...")
	select {
	case <-time.After(3 * time.Second):
		beewatch.Trace().Printf("Done debug")
	}

	// Close debugger,
	// it's not useful when you debug a web server that may not have an end.
	beewatch.Close()
}

func multipletest(num int) {
	beewatch.Critical().Break().Display("num", num)
}
