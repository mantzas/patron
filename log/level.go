package log

// The Level type definition.
type Level string

// LevelOrder defines the order for a level.
type LevelOrder int

const (
	// DebugLevel level.
	DebugLevel Level = "debug"
	// InfoLevel level.
	InfoLevel Level = "info"
	// WarnLevel level.
	WarnLevel Level = "warn"
	// ErrorLevel level.
	ErrorLevel Level = "error"
	// FatalLevel level.
	FatalLevel Level = "fatal"
	// PanicLevel level.
	PanicLevel Level = "panic"
	// NoLevel level.
	NoLevel Level = ""

	// DebugLevelOrder level.
	DebugLevelOrder LevelOrder = 0
	// InfoLevelOrder level.
	InfoLevelOrder LevelOrder = 1
	// WarnLevelOrder level.
	WarnLevelOrder LevelOrder = 2
	// ErrorLevelOrder level.
	ErrorLevelOrder LevelOrder = 3
	// FatalLevelOrder level.
	FatalLevelOrder LevelOrder = 4
	// PanicLevelOrder level.
	PanicLevelOrder LevelOrder = 5
	// NoLevelOrder level.
	NoLevelOrder LevelOrder = 6
)

var levelToOrder map[Level]LevelOrder
var orderToLevel map[LevelOrder]Level

func init() {
	levelToOrder = map[Level]LevelOrder{
		NoLevel:    NoLevelOrder,
		DebugLevel: DebugLevelOrder,
		InfoLevel:  InfoLevelOrder,
		WarnLevel:  WarnLevelOrder,
		ErrorLevel: ErrorLevelOrder,
		FatalLevel: FatalLevelOrder,
		PanicLevel: PanicLevelOrder,
	}
	orderToLevel = map[LevelOrder]Level{
		NoLevelOrder:    NoLevel,
		DebugLevelOrder: DebugLevel,
		InfoLevelOrder:  InfoLevel,
		WarnLevelOrder:  WarnLevel,
		ErrorLevelOrder: ErrorLevel,
		FatalLevelOrder: FatalLevel,
		PanicLevelOrder: PanicLevel,
	}
}

// LevelToOrder returns the order of the level.
func LevelToOrder(lvl Level) LevelOrder {
	return levelToOrder[lvl]
}

// OrderToLevel return the level of a order.
func OrderToLevel(order LevelOrder) Level {
	return orderToLevel[order]
}

// ShouldLog return true if we should log else false based on log level orders.
func ShouldLog(min, cur LevelOrder) bool {
	if min > cur {
		return false
	}
	return true
}
