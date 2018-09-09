package info

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	AddName("Name")
	AddVersion("1.2.3")
	AddMetric("Name", "Description")
	AddDoc("testdata/test.md")
	exp := info{
		Name:    "Name",
		Version: "1.2.3",
		Metrics: []metric{metric{Name: "Name", Description: "Description"}},
		Doc:     "<h1>Markdown: Syntax</h1>\n\n<p>This is the first paragraph.</p>\n\n<h2>Overview</h2>\n\n<p>This is the second paragraph.</p>\n",
	}
	assert.Equal(t, exp, serviceInfo)
	expected, err := json.Marshal(exp)
	assert.NoError(t, err)
	got, err := Marshal()
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestAddDoc(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name  string
		args  args
		empty bool
	}{
		{name: "no file", args: args{file: ""}, empty: true},
		{name: "file not exists", args: args{file: "file_not_exists.md"}, empty: true},
		{name: "success", args: args{file: "testdata/test.md"}, empty: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddDoc(tt.args.file)
			if tt.empty {
				assert.Empty(t, serviceInfo.Doc)
			} else {
				assert.NotEmpty(t, serviceInfo.Doc)
			}

		})
	}
}
