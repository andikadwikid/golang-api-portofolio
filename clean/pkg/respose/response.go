package response

import "github.com/gin-gonic/gin"

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string, err error) {

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	c.JSON(status, ErrorResponse{
		Status:  status,
		Message: message,
		Error:   errMsg,
	})
}
