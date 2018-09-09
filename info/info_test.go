package info

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	UpdateName("Name")
	UpdateVersion("1.2.3")
	UpdateHost("Host")
	UpsertMetric("Name", "Description", "Counter")
	UpsertConfig("Config", "Value")
	err := ImportDoc("testdata/test.md")
	assert.NoError(t, err)
	exp := info{
		Name:    "Name",
		Version: "1.2.3",
		Host:    "Host",
		Metrics: map[string]string{"Name": "[Counter] Description"},
		Doc:     "<h1>Markdown: Syntax</h1>\n\n<p>This is the first paragraph.</p>\n\n<h2>Overview</h2>\n\n<p>This is the second paragraph.</p>\n",
		Configs: map[string]string{"Config": "Value"},
	}
	assert.Equal(t, exp, serviceInfo)
	expected, err := json.Marshal(exp)
	assert.NoError(t, err)
	got, err := Marshal()
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(got))
}

func TestAddDoc(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
	}{
		{name: "no file", args: args{file: ""}, wantError: false},
		{name: "file not exists", args: args{file: "file_not_exists.md"}, wantError: false},
		{name: "success", args: args{file: "testdata/test.md"}, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ImportDoc(tt.args.file)
			if tt.wantError {
				assert.NoError(t, err)
				assert.NotEmpty(t, serviceInfo.Doc)
			} else {
				assert.Error(t, err)
				assert.Empty(t, serviceInfo.Doc)
			}
		})
	}
}
