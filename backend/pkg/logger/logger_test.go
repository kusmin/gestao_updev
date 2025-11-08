package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	// Test with a valid log level
	logger, err := New("debug")
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.True(t, logger.Core().Enabled(zap.DebugLevel))

	// Test with an invalid log level
	logger, err = New("invalid")
	assert.Error(t, err)
	assert.Nil(t, logger)

	// Test with an empty log level (should default to "info")
	logger, err = New("")
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.True(t, logger.Core().Enabled(zap.InfoLevel))
	assert.False(t, logger.Core().Enabled(zap.DebugLevel))
}
