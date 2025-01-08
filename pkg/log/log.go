package log

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogInfo map[string]interface{}

var logger zerolog.Logger

func GetLogger() *zerolog.Logger {
	return &logger
}

func init() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}

	logFileName := fmt.Sprintf("./data/logs/app-%s.log", time.Now().Format("2006-01-02"))
	fileWriter := &lumberjack.Logger{
		Filename:  logFileName,
		LocalTime: true,
		Compress:  true,
	}

	multi := zerolog.MultiLevelWriter(consoleWriter, fileWriter)

	logger = zerolog.New(multi).With().Timestamp().Logger()
}

func UpdateContext(key, value string) {
	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(key, value)
	})
}

func Trace(fields LogInfo, msg string) {
	var convFields map[string]interface{} = fields

	event := logger.Trace()
	for key, value := range convFields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func Debug(fields LogInfo, msg string) {
	var convFields map[string]interface{} = fields

	event := logger.Debug()
	for key, value := range convFields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func Info(fields LogInfo, msg string) {
	var convFields map[string]interface{} = fields

	event := logger.Info()
	for key, value := range convFields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func Warn(fields LogInfo, msg string) {
	var convFields map[string]interface{} = fields

	event := logger.Warn()
	for key, value := range convFields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func Error(fields LogInfo, msg string) {
	var convFields map[string]interface{} = fields

	event := logger.Error()
	for key, value := range convFields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func Fatal(fields LogInfo, msg string) {
	var convFields map[string]interface{} = fields

	event := logger.Fatal()
	for key, value := range convFields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func Panic(fields LogInfo, msg string) {
	var convFields map[string]interface{} = fields

	event := logger.Panic()
	for key, value := range convFields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}
