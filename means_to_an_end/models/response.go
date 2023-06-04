package models

type Response struct {
	Method byte
	Number *int32
	Error  error
}

func NewResponse(method byte, number *int32, error error) *Response {
	return &Response{
		Method: method,
		Number: number,
		Error:  error,
	}
}
