package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	CodeSuccess         = 0
	CodeAuthFailed      = 1001
	CodeSessionExpired  = 1002
	CodeResourceNotFound = 2001
	CodeUpdateFailed    = 2002
	CodeYAMLError       = 2003
	CodeK8sConnectError = 3001
	CodePermissionDenied = 3002
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, CodeUpdateFailed, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, CodeAuthFailed, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, CodePermissionDenied, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, CodeResourceNotFound, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, CodeK8sConnectError, message)
}
