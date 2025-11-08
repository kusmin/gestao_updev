package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	// setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// execute
	Success(c, http.StatusOK, gin.H{"foo": "bar"}, gin.H{"meta": "data"})

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"data":{"foo":"bar"},"meta":{"meta":"data"},"error":null}`, w.Body.String())
}

func TestError(t *testing.T) {
	// setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// execute
	Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request", gin.H{"details": "here"})

	// assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"data":null,"error":{"code":"BAD_REQUEST","message":"Invalid request","details":{"details":"here"}}}`, w.Body.String())
}
