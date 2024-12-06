package utils

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

type CustomLogger struct {
	zerolog.Logger
}

var Logger CustomLogger

func NewLogger() CustomLogger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i any) string {
		return strings.ToUpper(fmt.Sprintf("|%-6s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatErrFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s: ", i)
	}
	logger := zerolog.New(output).With().Caller().Timestamp().Logger()
	Logger = CustomLogger{logger}
	return Logger
}

func (l *CustomLogger) LogInfo() *zerolog.Event {
	return l.Logger.Info()
}

func (l *CustomLogger) LogError() *zerolog.Event {
	return l.Logger.Error()
}

func (l *CustomLogger) LogDebug() *zerolog.Event {
	return l.Logger.Debug()
}

func (l *CustomLogger) LogWarn() *zerolog.Event {
	return l.Logger.Warn()
}

func (l *CustomLogger) LogFatal() *zerolog.Event {
	return l.Logger.Fatal()
}

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// log the request
		start := time.Now()
		err := next(c)
		stop := time.Now()
		Logger.LogInfo().Fields(map[string]interface{}{
			"method":  c.Request().Method,
			"uri":     c.Request().URL.Path,
			"query":   c.Request().URL.RawQuery,
			"latency": stop.Sub(start),
		}).Msg("Request")
		if err != nil {
			Logger.LogError().Fields(map[string]interface{}{
				"error": err.Error(),
			}).Msg("Response")
			return err
		}
		return nil
	}
}
