package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestEnsureObjectDefinesGinHForNil(t *testing.T) {
	result := ensureObject(nil)
	require.Equal(t, gin.H{}, result)
}

func TestEnsureObjectReturnsValueAsIs(t *testing.T) {
	payload := gin.H{"foo": "bar"}

	result := ensureObject(payload)
	require.Equal(t, payload, result)
}

type httpResponse struct {
	Data  map[string]interface{} `json:"data"`
	Meta  map[string]interface{} `json:"meta"`
	Error struct {
		Code    string                 `json:"code"`
		Message string                 `json:"message"`
		Details map[string]interface{} `json:"details"`
	} `json:"error"`
}

func TestSuccessWritesExpectedStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	Success(ctx, http.StatusCreated, gin.H{"foo": "bar"}, gin.H{"page": 1})

	require.Equal(t, http.StatusCreated, recorder.Code)

	var resp httpResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.Equal(t, "bar", resp.Data["foo"])
	require.Equal(t, float64(1), resp.Meta["page"])
	require.Empty(t, resp.Error.Code)
	require.Empty(t, resp.Error.Message)
}

func TestErrorWritesErrorPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	Error(ctx, http.StatusBadRequest, "TEST_CODE", "mensagem", gin.H{"detail": "value"})

	require.Equal(t, http.StatusBadRequest, recorder.Code)

	var resp httpResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.Empty(t, resp.Data)
	require.Empty(t, resp.Meta)
	require.Equal(t, "TEST_CODE", resp.Error.Code)
	require.Equal(t, "mensagem", resp.Error.Message)
	require.Equal(t, "value", resp.Error.Details["detail"])
}
