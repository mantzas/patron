package log

type Logger interface {
	Fields() map[string]interface{}
	Error(...interface{})
	Errorf(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Panic(...interface{})
	Panicf(string, ...interface{})
}
