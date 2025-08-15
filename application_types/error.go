package application_types

type ApplicationError struct {
	httpStatus  int
	isSuccess   bool
	httpMessage string
	err         error
}

func NewApplicationError(isSuccess bool, httpStatus int, err error) *ApplicationError {
	return &ApplicationError{
		httpStatus: httpStatus,
		isSuccess:  isSuccess,
		err:        err,
	}
}

func (appErr *ApplicationError) GetErrorMessage() string {
	return appErr.err.Error()
}
