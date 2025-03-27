package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	assert.NotNil(t, NEO4J_DIVER)
}
