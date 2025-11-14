package contextutil

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/kusmin/gestao_updev/backend/internal/middleware"
)

func TestTenantIDSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := testContext()
	expected := uuid.New()
	ctx.Set(middleware.ContextTenantIDKey, expected.String())

	actual, err := TenantID(ctx)

	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestTenantIDErrorsWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := testContext()

	tenant, err := TenantID(ctx)

	require.Error(t, err)
	require.Equal(t, uuid.Nil, tenant)
}

func TestTenantIDErrorsWhenTypeInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := testContext()
	ctx.Set(middleware.ContextTenantIDKey, 123)

	tenant, err := TenantID(ctx)

	require.Error(t, err)
	require.Equal(t, uuid.Nil, tenant)
}

func TestUserIDSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := testContext()
	expected := uuid.New()
	ctx.Set(middleware.ContextUserIDKey, expected.String())

	actual, err := UserID(ctx)

	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestUserIDErrorsWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := testContext()

	uid, err := UserID(ctx)

	require.Error(t, err)
	require.Equal(t, uuid.Nil, uid)
}

func TestUserIDErrorsWhenTypeInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := testContext()
	ctx.Set(middleware.ContextUserIDKey, true)

	uid, err := UserID(ctx)

	require.Error(t, err)
	require.Equal(t, uuid.Nil, uid)
}

func testContext() *gin.Context {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	return ctx
}
