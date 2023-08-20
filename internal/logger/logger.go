package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	LogLevel string
	DevMode  bool
	Type     string
}

// Application logger
type appLogger struct {
	level       string
	devMode     bool
	encoding    string
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
}

func NewAppLogger(cfg *Config) *appLogger {
	return &appLogger{level: cfg.LogLevel, devMode: cfg.DevMode, encoding: cfg.Type}
}

func (l *appLogger) InitLogger() {
	logLevel, _ := zap.ParseAtomicLevel(l.level)

	logWriter := zapcore.AddSync(os.Stdout)

	var encoderCfg zapcore.EncoderConfig
	if l.devMode {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	var encoder zapcore.Encoder
	encoderCfg.NameKey = "[SERVICE]"
	encoderCfg.TimeKey = "[TIME]"
	encoderCfg.LevelKey = "[LEVEL]"
	encoderCfg.FunctionKey = "[CALLER]"
	encoderCfg.CallerKey = "[LINE]"
	encoderCfg.MessageKey = "[MESSAGE]"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoderCfg.EncodeName = zapcore.FullNameEncoder
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder

	if l.encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, logWriter, logLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.logger = logger
	l.sugarLogger = logger.Sugar()
}

type Logger interface {
	Fatal(args ...interface{})
	Info(args ...interface{})
	Infoln(args ...interface{})
	Debug(args ...interface{})
	Error(args ...interface{})
}

func (l *appLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *appLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l *appLogger) Infoln(args ...interface{}) {
	l.sugarLogger.Infoln(args...)
}
func (l *appLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l *appLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}
