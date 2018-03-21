package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStdFactory(t *testing.T) {
	assert := assert.New(t)
	f := NewStdFactory(nil)
	assert.NotNil(f)
}

func TestStdFactory_Create(t *testing.T) {
	assert := assert.New(t)
	l := NewStdFactory(nil).Create(nil)
	assert.NotNil(l)
}

func TestStdFactory_CreateWithFields(t *testing.T) {
	assert := assert.New(t)
	f := make(map[string]interface{})
	f["test"] = "test"
	l := NewStdFactory(nil).Create(f)
	assert.NotNil(l)
	assert.EqualValues(f, l.Fields())
}

func TestStdFactory_CreateSub(t *testing.T) {
	assert := assert.New(t)
	f := make(map[string]interface{})
	f["test"] = "test"
	fct := NewStdFactory(nil)
	l := fct.Create(f)
	f1 := make(map[string]interface{})
	f1["test1"] = "test1"
	sl := fct.CreateSub(l, f1)
	assert.NotNil(sl)
	assert.EqualValues(f, sl.Fields())
}

func TestStdFactory_CreateSub_WithoutFields(t *testing.T) {
	assert := assert.New(t)
	f := make(map[string]interface{})
	f["test"] = "test"
	fct := NewStdFactory(nil)
	l := fct.Create(f)
	sl := fct.CreateSub(l, nil)
	assert.NotNil(sl)
	assert.EqualValues(f, sl.Fields())
}
