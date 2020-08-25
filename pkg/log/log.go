package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
)

type Level = zerolog.Level
type Mode int8

const (
	PrettyMode Mode = iota
	JSONMode
)

var currentMode = new(int32)

func SetMode(mode Mode) {
	atomic.StoreInt32(currentMode, int32(mode))
}

func SetLevel(lvl Level) {
	zerolog.SetGlobalLevel(lvl)
}

func SetModeString(mode string) error {
	switch strings.ToUpper(mode) {
	case "", "PRETTY":
		SetMode(PrettyMode)
	case "JSON":
		SetMode(JSONMode)
	default:
		return fmt.Errorf("%s is not correct mode: PRETTY or JSON", mode)
	}
	return nil
}

func SetLevelString(level string) error {
	if level == "" {
		level = "info"
	}
	lvl, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		return err
	}
	SetLevel(lvl)
	return nil
}

func NewLogger() *Logger {
	return newLogger(2)
}

func newLogger(deep int) *Logger {
	zerolog.TimestampFieldName = "@timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.ErrorFieldName = "stack_trace"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	programCounter, file, _, _ := runtime.Caller(deep)
	funcName := runtime.FuncForPC(programCounter).Name()
	_, fileName := path.Split(file)
	var writer io.Writer
	switch Mode(atomic.LoadInt32(currentMode)) {
	case JSONMode:
		writer = os.Stdout
	case PrettyMode:
		writer = zerolog.NewConsoleWriter()
	}
	return &Logger{zerolog.
		New(writer).
		With().
		Timestamp().
		Str("@version", "1").
		Str("logger_name", fmt.Sprintf("%s:%s", extractPackageName(funcName), fileName)).
		Logger()}
}

// "root.of.module.with.dots/root_package/package/subpackage.func.lambda"
// we need "root.of.module.with.dots/root_package/package/subpackage"
func extractPackageName(funcName string) string {
	lastSlash := strings.LastIndex(funcName, "/") + 1
	firstDotAfterPackageName := strings.Index(funcName[lastSlash:], ".")
	return funcName[:lastSlash+firstDotAfterPackageName]
}
