package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRedis(t *testing.T) {
	db := NewRedis()
	assert.NoError(t, db.Ping())
}
