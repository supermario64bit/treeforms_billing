package application_types

import "github.com/gin-gonic/gin"

type ApplicationError struct {
	httpStatus  int
	isSuccess   bool
	httpMessage string
	err         error
}

func NewApplicationError(isSuccess bool, httpStatus int, httpMessage string, err error) *ApplicationError {
	return &ApplicationError{
		httpStatus:  httpStatus,
		httpMessage: httpMessage,
		isSuccess:   isSuccess,
		err:         err,
	}
}

func (appErr *ApplicationError) GetErrorMessage() string {
	return appErr.err.Error()
}

func (appErr *ApplicationError) WriteHTTPResponse(c *gin.Context) {
	responseBody := gin.H{}

	if appErr.isSuccess {
		responseBody["status"] = "success"
	} else {
		responseBody["status"] = "failed"

	}
	responseBody["message"] = appErr.httpMessage
	responseBody["result"] = gin.H{"error": appErr.GetErrorMessage()}

	c.JSON(appErr.httpStatus, responseBody)
}

func (appErr *ApplicationError) GetError() error {
	return appErr.err
}
