package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const prefixWidth = 15

var (
	loggerInstance *zap.Logger
	level          zap.AtomicLevel
)

func levelBadge(l logLevel) string {
	switch l {
	case levelDebug:
		return levelColorDebug
	case levelWarn:
		return levelColorWarn
	case levelError:
		return levelColorError
	case levelSuccess:
		return levelColorSuccess
	default:
		return levelColorInfo
	}
}

func formatPrefix(prefix string) string {
	if prefix == "" {
		return ""
	}

	padding := ""
	if len(prefix) < prefixWidth {
		padding = strings.Repeat(" ", prefixWidth-len(prefix))
	}

	return cyan + "[" + prefix + "]" + reset + padding
}

func colorMessage(l logLevel, msg string) string {
	switch l {
	case levelDebug:
		return messageColorDebug + msg + reset
	case levelWarn:
		return messageColorWarn + msg + reset
	case levelError:
		return messageColorError + msg + reset
	case levelSuccess:
		return messageColorSuccess + msg + reset
	default:
		return messageColorInfo + msg + reset
	}
}

func Init() {
	level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	encoderCfg := zapcore.EncoderConfig{
		MessageKey: "msg",
		LineEnding: "\n",
	}

	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	stdout := zapcore.AddSync(os.Stdout)
	stderr := zapcore.AddSync(os.Stderr)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, stdout, zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l < zapcore.WarnLevel && level.Enabled(l)
		})),
		zapcore.NewCore(encoder, stderr, zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= zapcore.WarnLevel && level.Enabled(l)
		})),
	)

	loggerInstance = zap.New(core)
}

func SetDebug(enabled bool) {
	if enabled {
		level.SetLevel(zapcore.DebugLevel)
	} else {
		level.SetLevel(zapcore.InfoLevel)
	}
}

func Debugf(prefix, format string, args ...any) {
	log(levelDebug, zapcore.DebugLevel, prefix, fmt.Sprintf(format, args...))
}

func Infof(prefix, format string, args ...any) {
	log(levelInfo, zapcore.InfoLevel, prefix, fmt.Sprintf(format, args...))
}

func Successf(prefix, format string, args ...any) {
	log(levelSuccess, zapcore.InfoLevel, prefix, fmt.Sprintf(format, args...))
}

func Warnf(prefix, format string, args ...any) {
	log(levelWarn, zapcore.WarnLevel, prefix, fmt.Sprintf(format, args...))
}

func Errorf(prefix, format string, args ...any) {
	log(levelError, zapcore.ErrorLevel, prefix, fmt.Sprintf(format, args...))
}

func Fatalf(prefix, format string, args ...any) {
	log(levelError, zapcore.ErrorLevel, prefix, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func log(label logLevel, zapLevel zapcore.Level, prefix string, msg any) {
	if loggerInstance == nil {
		panic("logger.Init() was not called")
	}

	message := fmt.Sprint(msg)
	colored := colorMessage(label, message)
	fullMessage := levelBadge(label) + " " + formatPrefix(prefix) + colored

	loggerInstance.Log(zapLevel, fullMessage)
}