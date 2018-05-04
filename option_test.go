package patron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

type testExporter struct {
}

func (e *testExporter) ExportView(vd *view.Data) {
}

func (e *testExporter) ExportSpan(vd *trace.SpanData) {
}

func TestMetric(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		e  view.Exporter
		rp time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{&testExporter{}, 10 * time.Second}, false},
		{"failed due to nil exporter", args{nil, 10 * time.Second}, true},
		{"failed due to min reporting interval", args{&testExporter{}, 100 * time.Millisecond}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Metric(tt.args.e, tt.args.rp)(&Server{})
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestTrace(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		e   trace.Exporter
		cfg trace.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{&testExporter{}, trace.Config{}}, false},
		{"failed due to nil exporter", args{nil, trace.Config{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Trace(tt.args.e, tt.args.cfg)(&Server{})
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
