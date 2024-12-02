package model

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseAPI struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Code int `json:"code"`
	error
}

func ThrowError(code int, e error) *ErrorResponse {
	return &ErrorResponse{
		Code:  code,
		error: e,
	}
}
