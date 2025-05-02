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
	output.FormatTimestamp = func(i any) string {
		return "\x1b[m" + i.(string) + "\x1b[0m"
	}
	output.FormatLevel = func(i any) string {
		return strings.ToUpper(fmt.Sprintf("|%-5s|", i))
	}
	output.FormatFieldName = func(i any) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i any) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatErrFieldName = func(i any) string {
		return fmt.Sprintf("%s: ", i)
	}

	logger := zerolog.New(output).With().Caller().Timestamp().Logger()
	Logger = CustomLogger{logger}
	return Logger
}

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		stop := time.Now()
		fields := map[string]interface{}{
			"method":  c.Request().Method,
			"uri":     c.Request().URL.Path,
			"latency": stop.Sub(start),
		}
		if c.Request().URL.RawQuery != "" {
			fields["query"] = c.Request().URL.RawQuery
		}
		Logger.Info().Fields(fields).Msg("Request")
		if err != nil {
			Logger.Error().Fields(map[string]interface{}{
				"error": err.Error(),
			}).Msg("Response")
			return err
		}
		return nil
	}
}
