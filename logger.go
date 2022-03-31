package logger

import (
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

func Info(message string) {
	if !initialized {
		log.Print(`Initialize logger before calling "logger.Info()"`)
	}
	if initialized && options.SentryEnabled {
		sentry.CaptureMessage(message)
	}
	log.Printf("Info: %v", message)
}

func Error(err error) {
	if !initialized {
		log.Print(`Initialize logger before calling "logger.Error()"`)
	}
	if initialized && options.SentryEnabled {
		sentry.CaptureException(err)
	}
	log.Printf("Error: %v", err)
}
