package log

// Factory interface for creating loggers.
type Factory interface {
	Create(map[string]interface{}) Logger
	CreateSub(Logger, map[string]interface{}) Logger
}
