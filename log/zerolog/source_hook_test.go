package zerolog

import (
	"bytes"
	"testing"

	"github.com/mantzas/patron/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_sourceHook(t *testing.T) {
	assert := assert.New(t)
	h := sourceHook{}
	var b bytes.Buffer
	zerolog.LevelFieldName = "lvl"
	zerolog.MessageFieldName = "msg"
	zl := zerolog.New(&b).Hook(&h)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Debugf("testing %d", 1)
	assert.Equal("{\"lvl\":\"debug\",\"key\":\"value\",\"src\":\"zerolog/source_hook_test.go:20\",\"msg\":\"testing 1\"}\n", b.String())
}

func Test_getSource(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		wantSrc string
	}{
		{name: "success", args: args{file: "/home/root/code.go"}, wantSrc: "root/code.go:1"},
		{name: "success without path", args: args{file: "code.go"}, wantSrc: "code.go:1"},
		{name: "success without path", args: args{file: ""}, wantSrc: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(tt.wantSrc, getSource(tt.args.file, 1))
		})
	}
}
func Benchmark_SourceEnabled(b *testing.B) {

	benchmarks := []struct {
		name         string
		enableSource bool
	}{
		{"without source hook", false},
		{"with source hook", false},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			var bf bytes.Buffer
			var zl zerolog.Logger
			if bm.enableSource {
				zl = zerolog.New(&bf).Hook(&sourceHook{})
			} else {
				zl = zerolog.New(&bf)
			}
			l := NewLogger(&zl, log.DebugLevel, f)
			l.Debugf("testing %d", 1)
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				l.Debugf("testing %d", 1)
				t = i
			}
		})
	}
}
