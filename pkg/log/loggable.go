package log

type Loggable struct {
	logger *Logger
}

func (l *Loggable) Log() *Logger {
	if l.logger == nil {
		l.logger = newLogger(2)
	}
	return l.logger
}
