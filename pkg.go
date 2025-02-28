package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/muesli/termenv"
)

var (
	// registry is a map of all registered lipgloss renderers.
	registry = sync.Map{}

	// defaultLogger is the default global logger instance.
	defaultLoggerOnce sync.Once
	defaultLogger     *Logger
)

// Default returns the default logger. The default logger comes with timestamp enabled.
func Default() *Logger {
	defaultLoggerOnce.Do(func() {
		if defaultLogger != nil {
			// already set via SetDefault.
			return
		}
		defaultLogger = NewWithOptions(os.Stderr, Options{ReportTimestamp: true})
	})
	return defaultLogger
}

// SetDefault sets the default global logger.
func SetDefault(logger *Logger) {
	defaultLogger = logger
}

// New returns a new logger with the default options.
func New(w io.Writer) *Logger {
	return NewWithOptions(w, Options{})
}

// NewWithOptions returns a new logger using the provided options.
func NewWithOptions(w io.Writer, o Options) *Logger {
	l := &Logger{
		b:               bytes.Buffer{},
		mu:              &sync.RWMutex{},
		helpers:         &sync.Map{},
		level:           int32(o.Level),
		reportTimestamp: o.ReportTimestamp,
		reportCaller:    o.ReportCaller,
		prefix:          o.Prefix,
		timeFunc:        o.TimeFunction,
		timeFormat:      o.TimeFormat,
		formatter:       o.Formatter,
		fields:          o.Fields,
		callerFormatter: o.CallerFormatter,
		callerOffset:    o.CallerOffset,
	}

	l.SetOutput(w)
	l.SetLevel(Level(l.level))
	l.SetStyles(DefaultStyles())

	if l.callerFormatter == nil {
		l.callerFormatter = ShortCallerFormatter
	}

	if l.timeFunc == nil {
		l.timeFunc = time.Now
	}

	if l.timeFormat == "" {
		l.timeFormat = DefaultTimeFormat
	}

	return l
}

// SetReportTimestamp sets whether to report timestamp for the default logger.
func SetReportTimestamp(report bool) {
	Default().SetReportTimestamp(report)
}

// SetReportCaller sets whether to report caller location for the default logger.
func SetReportCaller(report bool) {
	Default().SetReportCaller(report)
}

// SetLevel sets the level for the default logger.
func SetLevel(level Level) {
	Default().SetLevel(level)
}

// GetLevel returns the level for the default logger.
func GetLevel() Level {
	return Default().GetLevel()
}

// SetTimeFormat sets the time format for the default logger.
func SetTimeFormat(format string) {
	Default().SetTimeFormat(format)
}

// SetTimeFunction sets the time function for the default logger.
func SetTimeFunction(f TimeFunction) {
	Default().SetTimeFunction(f)
}

// SetOutput sets the output for the default logger.
func SetOutput(w io.Writer) {
	Default().SetOutput(w)
}

// SetFormatter sets the formatter for the default logger.
func SetFormatter(f Formatter) {
	Default().SetFormatter(f)
}

// SetCallerFormatter sets the caller formatter for the default logger.
func SetCallerFormatter(f CallerFormatter) {
	Default().SetCallerFormatter(f)
}

// SetCallerOffset sets the caller offset for the default logger.
func SetCallerOffset(offset int) {
	Default().SetCallerOffset(offset)
}

// SetPrefix sets the prefix for the default logger.
func SetPrefix(prefix string) {
	Default().SetPrefix(prefix)
}

// SetColorProfile force sets the underlying Lip Gloss renderer color profile
// for the TextFormatter.
func SetColorProfile(profile termenv.Profile) {
	Default().SetColorProfile(profile)
}

// SetStyles sets the logger styles for the TextFormatter.
func SetStyles(s *Styles) {
	Default().SetStyles(s)
}

// GetPrefix returns the prefix for the default logger.
func GetPrefix() string {
	return Default().GetPrefix()
}

// With returns a new logger with the given keyvals.
func With(keyvals ...interface{}) *Logger {
	return Default().With(keyvals...)
}

// WithPrefix returns a new logger with the given prefix.
func WithPrefix(prefix string) *Logger {
	return Default().WithPrefix(prefix)
}

// Helper marks the calling function as a helper
// and skips it for source location information.
// It's the equivalent of testing.TB.Helper().
func Helper() {
	Default().helper(1)
}

// Log logs a message with the given level.
func Log(level Level, msg interface{}, keyvals ...interface{}) {
	Default().Log(level, msg, keyvals...)
}

// Debug logs a debug message.
func Debug(msg interface{}, keyvals ...interface{}) {
	Default().Log(DebugLevel, msg, keyvals...)
}

// Info logs an info message.
func Info(msg interface{}, keyvals ...interface{}) {
	Default().Log(InfoLevel, msg, keyvals...)
}

// Warn logs a warning message.
func Warn(msg interface{}, keyvals ...interface{}) {
	Default().Log(WarnLevel, msg, keyvals...)
}

// Error logs an error message.
func Error(msg interface{}, keyvals ...interface{}) {
	Default().Log(ErrorLevel, msg, keyvals...)
}

// Fatal logs a fatal message and exit.
func Fatal(msg interface{}, keyvals ...interface{}) {
	Default().Log(FatalLevel, msg, keyvals...)
	os.Exit(1)
}

// Print logs a message with no level.
func Print(msg interface{}, keyvals ...interface{}) {
	Default().Log(noLevel, msg, keyvals...)
}

// Logf logs a message with formatting and level.
func Logf(level Level, format string, args ...interface{}) {
	Default().Logf(level, format, args...)
}

// Debugf logs a debug message with formatting.
func Debugf(format string, args ...interface{}) {
	Default().Log(DebugLevel, fmt.Sprintf(format, args...))
}

// Infof logs an info message with formatting.
func Infof(format string, args ...interface{}) {
	Default().Log(InfoLevel, fmt.Sprintf(format, args...))
}

// Warnf logs a warning message with formatting.
func Warnf(format string, args ...interface{}) {
	Default().Log(WarnLevel, fmt.Sprintf(format, args...))
}

// Errorf logs an error message with formatting.
func Errorf(format string, args ...interface{}) {
	Default().Log(ErrorLevel, fmt.Sprintf(format, args...))
}

// Fatalf logs a fatal message with formatting and exit.
func Fatalf(format string, args ...interface{}) {
	Default().Log(FatalLevel, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Printf logs a message with formatting and no level.
func Printf(format string, args ...interface{}) {
	Default().Log(noLevel, fmt.Sprintf(format, args...))
}

// StandardLog returns a standard logger from the default logger.
func StandardLog(opts ...StandardLogOptions) *log.Logger {
	return Default().StandardLog(opts...)
}
