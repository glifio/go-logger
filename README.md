# Glif.io Logger

Logger to be used in Glif.io Go modules, supports logging to [Sentry](https://sentry.io/).

## Usage

Always make sure to call `logger.Init()` at the start of the application. Failing to do so may cause the application to exit with an error when calling certain setup methods of the `logger` package.

Calling `logger.Info()` or `logger.Error()` is safe before `logger.Init()`, to prevent applications exiting unexpectedly in production, but will log an additional warning message.

```go
package main

import "github.com/glifio/go-logger"

func main() {
	// Initialize logger
	err := logger.Init(logger.LoggerOptions{
		ModuleName:    "verifier",
		SentryEnabled: true,
		SentryDsn:     "https://abc123.ingest.sentry.io/1234567",
		SentryEnv:     "Development",
		SentryTraces:  0,
	})
	if err != nil {
		log.Panic(err)
	}

	// Add sentry gin middleware (optional)
	router := gin.Default()
	if logger.IsSentryEnabled() {
		router.Use(logger.GetSentryGin())
	}

	// Log an info message
	logger.Info("It works!")

	// Log an error message
	logger.error(errors.New("Oh no, it doesn't.."))
}
```
