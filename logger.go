package stun

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// LogLevel represents the logging level
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
)

// Logger wraps logrus.Logger with custom configuration and structured logging
type Logger struct {
	log *log.Logger
}

// LoggerConfig holds configuration for the logger
type LoggerConfig struct {
	Level      LogLevel
	Format     string // "text" or "json"
	Output     string // "stdout" or "stderr"
	ShowCaller bool
}

// NewLogger creates a new logger with the given configuration
func NewLogger(config LoggerConfig) *Logger {
	logger := log.New()

	// Set output
	switch config.Output {
	case "stderr":
		logger.SetOutput(os.Stderr)
	default:
		logger.SetOutput(os.Stdout)
	}

	// Set formatter
	switch config.Format {
	case "json":
		logger.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	default:
		logger.SetFormatter(&log.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
			ForceColors:     true,
		})
	}

	// Set log level
	switch config.Level {
	case DebugLevel:
		logger.SetLevel(log.DebugLevel)
	case InfoLevel:
		logger.SetLevel(log.InfoLevel)
	case WarnLevel:
		logger.SetLevel(log.WarnLevel)
	case ErrorLevel:
		logger.SetLevel(log.ErrorLevel)
	case FatalLevel:
		logger.SetLevel(log.FatalLevel)
	default:
		logger.SetLevel(log.InfoLevel)
	}

	// Enable caller information if requested
	if config.ShowCaller {
		logger.SetReportCaller(true)
	}

	return &Logger{
		log: logger,
	}
}

// NewDefaultLogger creates a logger with default configuration
func NewDefaultLogger() *Logger {
	return NewLogger(LoggerConfig{
		Level:      InfoLevel,
		Format:     "text",
		Output:     "stdout",
		ShowCaller: false,
	})
}

// Debug logs a message at debug level
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.log.WithFields(fields[0]).Debug(msg)
	} else {
		l.log.Debug(msg)
	}
}

// Info logs a message at info level
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.log.WithFields(fields[0]).Info(msg)
	} else {
		l.log.Info(msg)
	}
}

// Warn logs a message at warn level
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.log.WithFields(fields[0]).Warn(msg)
	} else {
		l.log.Warn(msg)
	}
}

// Error logs a message at error level
func (l *Logger) Error(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.log.WithFields(fields[0]).Error(msg)
	} else {
		l.log.Error(msg)
	}
}

// Fatal logs a message at fatal level and exits
func (l *Logger) Fatal(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.log.WithFields(fields[0]).Fatal(msg)
	} else {
		l.log.Fatal(msg)
	}
}

// LogRequest logs STUN request details
func (l *Logger) LogRequest(remoteAddr string, msgType MessageType, transactionID [12]byte) {
	l.Info("STUN request received", map[string]interface{}{
		"remote_addr":    remoteAddr,
		"message_type":   msgType.String(),
		"transaction_id": transactionID,
		"component":      "stun_server",
	})
}

// LogResponse logs STUN response details
func (l *Logger) LogResponse(remoteAddr string, msgType MessageType, transactionID [12]byte, xorAddr *XorMappedAddr) {
	fields := map[string]interface{}{
		"remote_addr":    remoteAddr,
		"message_type":   msgType.String(),
		"transaction_id": transactionID,
		"component":      "stun_server",
	}

	if xorAddr != nil {
		fields["xor_mapped_ip"] = xorAddr.IP.String()
		fields["xor_mapped_port"] = xorAddr.Port
	}

	l.Info("STUN response sent", fields)
}

// LogError logs error details with context
func (l *Logger) LogError(msg string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err.Error()
	fields["component"] = "stun_server"
	l.Error(msg, fields)
}

// LogClientRequest logs client request details
func (l *Logger) LogClientRequest(serverAddr string, msgType MessageType, transactionID [12]byte) {
	l.Debug("STUN client request", map[string]interface{}{
		"server_addr":    serverAddr,
		"message_type":   msgType.String(),
		"transaction_id": transactionID,
		"component":      "stun_client",
	})
}

// LogClientResponse logs client response details
func (l *Logger) LogClientResponse(serverAddr string, msgType MessageType, xorAddr *XorMappedAddr) {
	fields := map[string]interface{}{
		"server_addr":  serverAddr,
		"message_type": msgType.String(),
		"component":    "stun_client",
	}

	if xorAddr != nil {
		fields["xor_mapped_ip"] = xorAddr.IP.String()
		fields["xor_mapped_port"] = xorAddr.Port
	}

	l.Info("STUN client response received", fields)
}

// LogConnection logs connection details
func (l *Logger) LogConnection(localAddr, remoteAddr string, component string) {
	l.Info("Connection established", map[string]interface{}{
		"local_addr":  localAddr,
		"remote_addr": remoteAddr,
		"component":   component,
	})
}

// LogShutdown logs shutdown details
func (l *Logger) LogShutdown(component string, duration time.Duration) {
	l.Info("Component shutdown", map[string]interface{}{
		"component": component,
		"duration":  duration.String(),
	})
}
