package async

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockMesssageProcessor_Process(t *testing.T) {
	assert := assert.New(t)
	m := MockMesssageProcessor{}
	assert.NoError(m.Process(context.TODO(), []byte{}))
}
