# Logsaw

[![Go Report Card](https://goreportcard.com/badge/github.com/mertakinstd/logsaw)](https://goreportcard.com/report/github.com/mertakinstd/logsaw)

A simple, colorful, and JSON-capable logging library for Go applications.

## Features

- Colorful console output support
- JSON formatted logging
- Multiple log levels (DEBUG, INFO, WARNING, ERROR, FATAL, PANIC)
- Timestamp support
- Simple and easy-to-use API

## Installation

```bash
go get github.com/mertakinstd/logsaw
```

## Usage

```go
package main

import saw "github.com/mertakinstd/logsaw"

func main() {
    // Initialize the logger
    log := saw.Initialize()

    // Configure for colored output
    log.SetConfig(saw.SawConfig{
        Colors: true,
    })

    // Console logging
    log.Debug("Debug message")
    log.Info("Info message")
    log.Warning("Warning message")
    log.Error("Error message")
    log.Fatal("Fatal message") // Will exit the program
    log.Panic("Panic message") // Will panic

    // JSON logging
    jsonDebug := log.JSON.Debug("Debug message")
    jsonInfo := log.JSON.Info("Info message")
    // ... etc
}
```

## Example Output

### Console Output

```
----------------------------------------------------------------
Log level: INFO
Log message: Server started on port 8080
Log time: 15:04:05
----------------------------------------------------------------
```

### JSON Output

```json
{
	"Lvl": "INFO",
	"Msg": "Server started on port 8080",
	"Time": 1683720245
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
