package async

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockProcessor_Process(t *testing.T) {
	assert := assert.New(t)
	m := MockProcessor{}
	assert.NoError(m.Process(context.TODO(), &Message{}))
}
