package response

import "github.com/gin-gonic/gin"

// APIError representa erros padronizados.
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// APIResponse padroniza todas as respostas JSON.
type APIResponse struct {
	Data  interface{} `json:"data"`
	Meta  interface{} `json:"meta"`
	Error interface{} `json:"error"`
}

func ensureObject(value interface{}) interface{} {
	if value == nil {
		return gin.H{}
	}
	return value
}

// Success envia uma resposta com dados válidos.
func Success(c *gin.Context, status int, data interface{}, meta interface{}) {
	c.JSON(status, APIResponse{
		Data:  ensureObject(data),
		Meta:  ensureObject(meta),
		Error: gin.H{},
	})
}

// Error envia uma resposta consistente de erro.
func Error(c *gin.Context, status int, code, message string, details interface{}) {
	c.JSON(status, APIResponse{
		Data: gin.H{},
		Meta: gin.H{},
		Error: APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
