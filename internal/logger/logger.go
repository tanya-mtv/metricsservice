package logger

type appLogger struct {
}

type Logger interface {
	InitLogger()
}

func NewAppLogger() *appLogger {
	return &appLogger{}
}

func (l *appLogger) InitLogger() {

}
