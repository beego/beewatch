Bee Watch
========

[![Build Status](https://drone.io/github.com/beego/beewatch/status.png)](https://drone.io/github.com/beego/beewatch/latest)

Bee Watch is an interactive debugger for the Go programming language.

## Features

- Use `Critical`, `Info` and `Trace` three levels to change debugger behavior.
- `Display()` variable values or `Printf()` with customized format.
- User-friendly Web UI for **WebSocket** mode or easy-control **command-line** mode.
- Call `AddWatchVars()` to monitor variables and show their information when the program calls `Break()`.
- Configuration file with your customized settings(`beewatch.json`).

## Installation

Bee Watch is a "go get" able Go project, you can execute the following command to auto-install:

	go get github.com/beego/beewatch

**Attention** This project can only be installed by source code now.

## Quick start

[beego.me/docs/Reference_BeeWatch](http://beego.me/docs/Reference_BeeWatch)

## Credits

- [emicklei/hopwatch](https://github.com/emicklei/hopwatch)
- [realint/dbgutil](https://github.com/realint/dbgutil)

## Examples and API documentation

[Go Walker](http://gowalker.org/github.com/beego/beewatch)


## License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).