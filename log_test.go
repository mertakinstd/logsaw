package log

import (
	"encoding/json"
	"testing"
)

func TestInitialize(t *testing.T) {
	saw := Initialize()
	if saw == nil {
		t.Error("Initialize() returned nil")
		return
	}
	if saw.JSON == nil {
		t.Error("JSON logger not initialized")
	}
}

func TestSetConfig(t *testing.T) {
	saw := Initialize()
	config := SawConfig{Colors: true}
	saw.SetConfig(config)
	if !saw.config.Colors {
		t.Error("Config not set correctly")
	}
}

func TestJsonLogOutput(t *testing.T) {
	saw := Initialize()
	jsonLog := saw.JSON

	tests := []struct {
		level string
		msg   string
		fn    func(string) []byte
	}{
		{"DEBUG", "debug message", jsonLog.Debug},
		{"INFO", "info message", jsonLog.Info},
		{"WARNING", "warning message", jsonLog.Warning},
		{"ERROR", "error message", jsonLog.Error},
		{"FATAL", "fatal message", jsonLog.Fatal},
		{"PANIC", "panic message", jsonLog.Panic},
	}

	for _, tt := range tests {
		output := tt.fn(tt.msg)
		var log Log
		err := json.Unmarshal(output, &log)
		if err != nil {
			t.Errorf("Failed to unmarshal JSON for level %s: %v", tt.level, err)
		}
		if log.Lvl != tt.level {
			t.Errorf("Expected level %s, got %s", tt.level, log.Lvl)
		}
		if log.Msg != tt.msg {
			t.Errorf("Expected message %s, got %s", tt.msg, log.Msg)
		}
		if log.Time == 0 {
			t.Error("Time not set in JSON log")
		}
	}
}

func TestLogLevels(t *testing.T) {
	saw := Initialize()

	// Test non-fatal logs
	saw.Debug("debug test")
	saw.Info("info test")
	saw.Warning("warning test")
	saw.Error("error test")

	// Test panic with recovery
	defer func() {
		if r := recover(); r == nil {
			t.Error("Panic log did not panic")
		}
	}()
	saw.Panic("panic test")
}
