package models

type Request struct {
	Method  byte  `json:"method"`
	Number1 int32 `json:"number1"`
	Number2 int32 `json:"number2"`
}

func NewRequest(method byte, number1 int32, number2 int32) *Request {
	return &Request{
		Method:  method,
		Number1: number1,
		Number2: number2,
	}
}
