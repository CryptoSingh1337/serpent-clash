package utils

import (
	"fmt"
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
