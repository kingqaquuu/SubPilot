package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, status int, code string, message string) {
	c.JSON(status, Body{
		Code:    status,
		Message: message,
		Data: gin.H{
			"error": code,
		},
	})
}
