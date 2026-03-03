package logger

type logLevel string

const (
	levelDebug   logLevel = "debug"
	levelInfo    logLevel = "info"
	levelWarn    logLevel = "warn"
	levelError   logLevel = "error"
	levelSuccess logLevel = "success"
)

const (
	reset = "\033[0m"
	cyan  = "\033[36m"

	levelColorInfo    = "\033[34mINFO   \033[0m"
	levelColorWarn    = "\033[33mWARN   \033[0m"
	levelColorError   = "\033[31mERROR  \033[0m"
	levelColorDebug   = "\033[35mDEBUG  \033[0m"
	levelColorSuccess = "\033[32mSUCCESS\033[0m"

	messageColorInfo    = "\033[97m"
	messageColorWarn    = "\033[33m"
	messageColorError   = "\033[31m"
	messageColorDebug   = "\033[90m"
	messageColorSuccess = "\033[32m"
)