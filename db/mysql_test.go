package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDB(t *testing.T) {
	db := NewMySQL()
	assert.NoError(t, db.Ping())
}
