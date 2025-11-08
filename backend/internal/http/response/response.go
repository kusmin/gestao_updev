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
	Meta  interface{} `json:"meta,omitempty"`
	Error interface{} `json:"error"`
}

// Success envia uma resposta com dados válidos.
func Success(c *gin.Context, status int, data interface{}, meta interface{}) {
	c.JSON(status, APIResponse{
		Data:  data,
		Meta:  meta,
		Error: nil,
	})
}

// Error envia uma resposta consistente de erro.
func Error(c *gin.Context, status int, code, message string, details interface{}) {
	c.JSON(status, APIResponse{
		Data: nil,
		Error: APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
