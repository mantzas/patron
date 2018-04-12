package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevelToOrder(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		lvl  Level
		want LevelOrder
	}{
		{"no", NoLevel, NoLevelOrder},
		{"debug", DebugLevel, DebugLevelOrder},
		{"info", InfoLevel, InfoLevelOrder},
		{"warn", WarnLevel, WarnLevelOrder},
		{"error", ErrorLevel, ErrorLevelOrder},
		{"panic", PanicLevel, PanicLevelOrder},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(tt.want, LevelToOrder(tt.lvl))
		})
	}
}

func TestOrderToLevel(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name  string
		order LevelOrder
		want  Level
	}{
		{"no", NoLevelOrder, NoLevel},
		{"debug", DebugLevelOrder, DebugLevel},
		{"info", InfoLevelOrder, InfoLevel},
		{"warn", WarnLevelOrder, WarnLevel},
		{"error", ErrorLevelOrder, ErrorLevel},
		{"panic", PanicLevelOrder, PanicLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(tt.want, OrderToLevel(tt.order))
		})
	}
}

func TestShouldLog(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		min  LevelOrder
		cur  LevelOrder
		want bool
	}{
		{"log with min warn and error", WarnLevelOrder, ErrorLevelOrder, true},
		{"log with min warn and warn", WarnLevelOrder, WarnLevelOrder, true},
		{"no log with min warn and info", WarnLevelOrder, InfoLevelOrder, false},
		{"no log with min no and debug", NoLevelOrder, DebugLevelOrder, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(tt.want, ShouldLog(tt.min, tt.cur))
		})
	}
}
