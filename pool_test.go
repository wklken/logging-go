package logging

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogEntryPool(t *testing.T) {
	pool := newLogEntryPool()

	logger := &logrus.Logger{}

	// get one
	e := pool.Get(logger)
	assert.NotNil(t, e)

	e.Data["hello"] = "world"

	// put back
	pool.Put(e)

	e1 := pool.Get(logger)
	assert.NotNil(t, e1)
	// the new entry is empty
	assert.Empty(t, e1.Data)
}
