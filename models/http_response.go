package models

type ErrorInfo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type HttpResponse struct {
	Result bool        `json:"result"`
	Data   interface{} `json:"data,omitempty"`
	Error  *ErrorInfo  `json:"error,omitempty"`
}

func NewSuccessResonse(data interface{}) HttpResponse {
	return HttpResponse{true, data, nil}
}

func NewErrorResonse(code int, msg string) HttpResponse {
	return HttpResponse{true, nil, &ErrorInfo{code, msg}}
}

func NewErrorResonseWithInfo(errInfo *ErrorInfo) HttpResponse {
	return HttpResponse{true, nil, errInfo}
}