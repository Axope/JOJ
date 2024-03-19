package response

const (
	ERROR   = 7
	SUCCESS = 0
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SuccessMsg(data interface{}) *Response {
	msg := &Response{
		Code: SUCCESS,
		Msg:  "SUCCESS",
		Data: data,
	}
	return msg
}

func FailMsg(message string) *Response {
	msgObj := &Response{
		Code: ERROR,
		Msg:  message,
		Data: nil,
	}
	return msgObj
}
