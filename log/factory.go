package log

// Factory interface defines the interface that we have to provide
// in order to use this abstraction
type Factory interface {
	Create(map[string]interface{}) Logger
	CreateSub(Logger, map[string]interface{}) Logger
}
