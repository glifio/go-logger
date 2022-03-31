package logger

import (
	"errors"
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

type LoggerOptions struct {
	ModuleName    string
	SentryEnabled bool
	SentryDsn     string
	SentryEnv     string
	SentryLevel   LogLevel
	SentryTraces  float64
}

var initialized bool
var options LoggerOptions

func init() {
	initialized = false
}

func Init(loggerOptions LoggerOptions) error {
	options = loggerOptions

	if options.SentryEnabled {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              options.SentryDsn,
			Environment:      options.SentryEnv,
			Release:          options.ModuleName,
			TracesSampleRate: options.SentryTraces,
			AttachStacktrace: true,
			Debug:            false,
		})
		if err != nil {
			return err
		}
	}

	initialized = true
	return nil
}

func IsSentryEnabled() bool {
	if !initialized {
		log.Fatal(`Initialize logger before calling "logger.IsSentryEnabled()"`)
	}
	return options.SentryEnabled
}

func GetSentryGin() gin.HandlerFunc {
	if !initialized || !options.SentryEnabled {
		log.Fatal(`Initialize logger with Sentry enabled before calling "logger.GetSentryGin()"`)
	}
	return sentrygin.New(sentrygin.Options{
		Repanic: true,
	})
}

func Debugf(format string, v ...any) {
	Debug(fmt.Sprintf(format, v...))
}

func Infof(format string, v ...any) {
	Info(fmt.Sprintf(format, v...))
}

func Warningf(format string, v ...any) {
	Warning(fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...any) {
	Error(errors.New(fmt.Sprintf(format, v...)))
}

func Debug(message string) {
	if !initialized {
		log.Print(`Initialize logger before calling "logger.Debug()"`)
	}
	formatted := fmt.Sprintf("Debug: %v", message)
	sendMessageToSentry(LogLevelDebug, formatted)
	log.Printf(formatted)
}

func Info(message string) {
	if !initialized {
		log.Print(`Initialize logger before calling "logger.Info()"`)
	}
	formatted := fmt.Sprintf("Info: %v", message)
	sendMessageToSentry(LogLevelInfo, formatted)
	log.Printf(formatted)
}

func Warning(message string) {
	if !initialized {
		log.Print(`Initialize logger before calling "logger.Warning()"`)
	}
	formatted := fmt.Sprintf("Warning: %v", message)
	sendMessageToSentry(LogLevelWarning, formatted)
	log.Printf(formatted)
}

func Error(err error) {
	if !initialized {
		log.Print(`Initialize logger before calling "logger.Error()"`)
	}
	sendErrorToSentry(err)
	log.Printf("Error: %v", err)
}

func sendMessageToSentry(level LogLevel, message string) {
	if initialized && options.SentryEnabled && level >= options.SentryLevel {
		sentry.CaptureMessage(message)
	}
}

func sendErrorToSentry(err error) {
	if initialized && options.SentryEnabled {
		sentry.CaptureException(err)
	}
}
