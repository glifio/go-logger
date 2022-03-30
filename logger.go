package logger

import (
	"log"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

type LoggerOptions struct {
	ModuleName    string
	SentryEnabled bool
	SentryDsn     string
	SentryEnv     string
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
	return options.SentryEnabled
}

func GetSentryGin() gin.HandlerFunc {
	if !initialized || !options.SentryEnabled {
		log.Fatal("Initialize logger with Sentry enabled before getting Sentry Gin middleware")
	}
	return sentrygin.New(sentrygin.Options{
		Repanic: true,
	})
}

func Info(info string) {
	if !initialized {
		log.Printf("Logger not initialized, failed to log info properly: %v", info)
		return
	}
	if options.SentryEnabled {
		sentry.CaptureMessage(info)
	}
	log.Printf("Info: %v", info)
}

func Error(err error) {
	if !initialized {
		log.Printf("Logger not initialized, failed to log error properly: %v", err)
		return
	}
	if options.SentryEnabled {
		sentry.CaptureException(err)
	}
	log.Printf("Error: %v", err)
}
