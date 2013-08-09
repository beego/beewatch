Bee Watch
========

Bee Watch is an interactive debugger for the Go programming language.

## Installation

Bee Watch is a "go get" able Go project, you can execute the following command to auto-install:

	go get github.com/beego/beewatch

**Attention** This project can only be installed by source code now.

## Quick start

### Usage

	import (
		"github.com/beego/beewatch"
	)

	func main() {
		appName := "Bee Watch"
		// Suspends execution until hitting "Resume" in the browser.
		beewatch.Display("App Name", appName).Break()
	}

### Connect

Bee Watch debugger is automatically started on [http://localhost:23456](http://localhost:23456), you can change port and other configuration by editing `beewatch.json`(copy default from Bee Watch source folder).

You browser has to support WebSocket, it has been tested with Chrome, Safari and Firefox on Mac and Windows.

### Other code examples

## Credits

- [hopwatch](github.com/emicklei/hopwatch)

## API Documentation

[Go Walker](http://gowalker.org/github.com/beego/beewatch)


## License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).