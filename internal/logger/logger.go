package logger

type appLogger struct {
	// sugarLogger *zap.SugaredLogger
	// logger      *zap.Logger
}

type Logger interface {
	InitLogger()
}

// NewAppLogger App Logger constructor
func NewAppLogger() *appLogger {
	return &appLogger{}
}

// func InitLogger() {
func (l *appLogger) InitLogger() {

}
