package env

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("C", "D")
	var correctBuf bytes.Buffer
	correctBuf.WriteString("A=B")
	var invalidBuf bytes.Buffer
	invalidBuf.WriteString("AB")
	var correctExistBuf bytes.Buffer
	correctExistBuf.WriteString("C=E")
	var invalidKVBuf bytes.Buffer
	invalidKVBuf.WriteString("=B")

	tests := []struct {
		name     string
		r        io.Reader
		wantErr  bool
		key      string
		expected string
	}{
		{"success with no file", nil, false, "C", "D"},
		{"success with file", &correctBuf, false, "A", "B"},
		{"success with skipped env", &correctExistBuf, false, "C", "D"},
		{"failed with invalid file", &invalidBuf, true, "", ""},
		{"failed with invalid key,value", &invalidKVBuf, true, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.r)

			if tt.wantErr {
				assert.Error(err)
				assert.Nil(got)
			} else {
				assert.NoError(err)
				assert.NotNil(got)
				assert.Equal(tt.expected, os.Getenv(tt.key))
			}
		})
	}
}

func TestConfig_Set(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("C", "D")
	tests := []struct {
		name     string
		key      string
		value    interface{}
		expected string
		wantErr  bool
	}{
		{"success", "A", "B", "B", false},
		{"success", "C", "E", "E", false},
		{"failure missing key", "", "B", "", true},
		{"failure nil value", "A", nil, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Config{}.Set(tt.key, tt.value)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(tt.expected, os.Getenv(tt.key))
			}
		})
	}
}

func TestConfig_Get(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("A", "B")
	tests := []struct {
		name    string
		key     string
		want    interface{}
		wantErr bool
	}{
		{"success", "A", "B", false},
		{"failed not exists", "Z", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Config{}.Get(tt.key)
			if tt.wantErr {
				assert.Error(err)
				assert.Nil(got)
			} else {
				assert.NoError(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}

func TestConfig_GetBool(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("GETBOOL", "true")
	os.Setenv("GETBOOL_FAIL", "ABC")
	tests := []struct {
		name    string
		key     string
		want    bool
		wantErr bool
	}{
		{"success", "GETBOOL", true, false},
		{"failed not exists", "XXX", false, true},
		{"failed not bool", "GETBOOL_FAIL", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Config{}.GetBool(tt.key)
			if tt.wantErr {
				assert.Error(err)
				assert.False(got)
			} else {
				assert.NoError(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}

func TestConfig_GetInt64(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("GETINT64", "1")
	os.Setenv("GETINT64_FAIL", "ABC")
	tests := []struct {
		name    string
		key     string
		want    int64
		wantErr bool
	}{
		{"success", "GETINT64", 1, false},
		{"failed not exists", "XXX", 0, true},
		{"failed not bool", "GETINT64_FAIL", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Config{}.GetInt64(tt.key)
			if tt.wantErr {
				assert.Error(err)
				assert.Zero(got)
			} else {
				assert.NoError(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}

func TestConfig_GetString(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("GETSTRING", "STRING")
	tests := []struct {
		name    string
		key     string
		want    string
		wantErr bool
	}{
		{"success", "GETSTRING", "STRING", false},
		{"failed not exists", "XXX", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Config{}.GetString(tt.key)
			if tt.wantErr {
				assert.Error(err)
				assert.Empty(got)
			} else {
				assert.NoError(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}

func TestConfig_GetFloat64(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("GETFLOAT64", "1.99")
	os.Setenv("GETFLOAT64_FAIL", "ABC")
	tests := []struct {
		name    string
		key     string
		want    float64
		wantErr bool
	}{
		{"success", "GETFLOAT64", 1.99, false},
		{"failed not exists", "XXX", 0.0, true},
		{"failed not bool", "GETFLOAT64_FAIL", 0.0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Config{}.GetFloat64(tt.key)
			if tt.wantErr {
				assert.Error(err)
				assert.Zero(got)
			} else {
				assert.NoError(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}
