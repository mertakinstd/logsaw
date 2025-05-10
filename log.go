package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	colorReset   = "\033[0m"
	colorDebug   = "\033[0;36m"
	colorInfo    = "\033[0;32m"
	colorWarning = "\033[0;33m"
	colorError   = "\033[0;31m"
	colorFatal   = "\033[0;35m"
	colorPanic   = "\033[0;41m"
)

// Saw is a logger that supports both console and JSON formatted logging
type Saw struct {
	JSON   *JsonLog
	config SawConfig
}

// SawConfig holds the configuration for the logger
type SawConfig struct {
	Colors bool
}

// Log represents a log message with its level, message, time, and context
type Log struct {
	Lvl  string
	Msg  string
	Time int64
}

// JsonLog is a logger that formats log messages as JSON
type JsonLog struct {
	saw     *Saw
	encoder *json.Encoder
	buffer  *bytes.Buffer
}

// builderPool is a sync.Pool for reusing strings.Builder instances
var builderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// logPool is a sync.Pool for reusing Log instances
var logPool = sync.Pool{
	New: func() interface{} {
		return &Log{}
	},
}

// Initialize creates and returns a new Saw instance with initialized headContext
func Initialize() *Saw {
	saw := &Saw{}
	buffer := &bytes.Buffer{}
	jsonLog := &JsonLog{
		saw:     saw,
		encoder: json.NewEncoder(buffer),
		buffer:  buffer,
	}

	jsonLog.encoder.SetEscapeHTML(false)
	saw.JSON = jsonLog

	return saw
}

// SetConfig updates the logger configuration and returns the Saw instance
func (s *Saw) SetConfig(config SawConfig) *Saw {
	s.config = config
	return s
}

// newLog creates a new log message with the specified level and message
func (s *Saw) newLog(level string, msg string) {
	log := logPool.Get().(*Log)
	defer logPool.Put(log)

	log.Lvl = level
	log.Msg = msg
	log.Time = time.Now().Unix()

	s.printLog(*log)
}

// printLog formats and prints the log message to the console
func (s *Saw) printLog(log Log) {
	color := colorReset

	if s.config.Colors {
		switch log.Lvl {
		case "DEBUG":
			color = colorDebug
		case "INFO":
			color = colorInfo
		case "WARNING":
			color = colorWarning
		case "ERROR":
			color = colorError
		case "FATAL":
			color = colorFatal
		case "PANIC":
			color = colorPanic
		default:
			color = colorReset
		}
	}

	sb := builderPool.Get().(*strings.Builder)
	defer func() {
		sb.Reset()
		builderPool.Put(sb)
	}()

	sb.WriteString(fmt.Sprintf("%s----------------------------------------------------------------%s\n", color, colorReset))
	sb.WriteString(fmt.Sprintf("%sLog level: %s%s\n", color, log.Lvl, colorReset))
	sb.WriteString(fmt.Sprintf("%sLog message: %s%s\n", color, log.Msg, colorReset))
	sb.WriteString(fmt.Sprintf("%sLog time: %s%s\n", color, time.Unix(log.Time, 0).Format(time.TimeOnly), colorReset))
	sb.WriteString(fmt.Sprintf("%s----------------------------------------------------------------%s\n", color, colorReset))

	output := sb.String()
	switch log.Lvl {
	case "PANIC":
		panic(output)
	case "FATAL":
		fmt.Print(output)
		os.Exit(1)
	default:
		fmt.Print(output)
	}
}

// Debug prints a debug level log message
func (s *Saw) Debug(msg string) {
	s.newLog("DEBUG", msg)
}

// Info prints an info level log message
func (s *Saw) Info(msg string) {
	s.newLog("INFO", msg)
}

// Warning prints a warning level log message
func (s *Saw) Warning(msg string) {
	s.newLog("WARNING", msg)
}

// Error prints an error level log message
func (s *Saw) Error(msg string) {
	s.newLog("ERROR", msg)
}

// Fatal prints a fatal level log message and exits the program
func (s *Saw) Fatal(msg string) {
	s.newLog("FATAL", msg)
}

// Panic prints a panic level log message and panics
func (s *Saw) Panic(msg string) {
	s.newLog("PANIC", msg)
}

// newJSONLog creates a new JSON formatted log message
func (j *JsonLog) newJSONLog(level string, msg string) []byte {
	j.buffer.Reset()

	log := logPool.Get().(*Log)
	defer logPool.Put(log)

	log.Lvl = level
	log.Msg = msg
	log.Time = time.Now().Unix()

	err := j.encoder.Encode(log)
	if err != nil {
		j.buffer.Reset()
		j.buffer.WriteString(fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
		return j.buffer.Bytes()
	}

	return j.buffer.Bytes()
}

// Debug returns a JSON formatted debug level log message
func (j *JsonLog) Debug(msg string) []byte {
	return j.newJSONLog("DEBUG", msg)
}

// Info returns a JSON formatted info level log message
func (j *JsonLog) Info(msg string) []byte {
	return j.newJSONLog("INFO", msg)
}

// Warning returns a JSON formatted warning level log message
func (j *JsonLog) Warning(msg string) []byte {
	return j.newJSONLog("WARNING", msg)
}

// Error returns a JSON formatted error level log message
func (j *JsonLog) Error(msg string) []byte {
	return j.newJSONLog("ERROR", msg)
}

// Fatal returns a JSON formatted fatal level log message
func (j *JsonLog) Fatal(msg string) []byte {
	return j.newJSONLog("FATAL", msg)
}

// Panic returns a JSON formatted panic level log message
func (j *JsonLog) Panic(msg string) []byte {
	return j.newJSONLog("PANIC", msg)
}
