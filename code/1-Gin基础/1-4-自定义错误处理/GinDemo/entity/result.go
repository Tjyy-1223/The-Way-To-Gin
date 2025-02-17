package entity

// Result 定义 Result 结构体
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// 定义错误码
const (
	CODE_SUCCESS int = 1
	CODE_ERROR   int = -1
)

// SetCode 设置错误码
func (res *Result) SetCode(code int) *Result {
	res.Code = code
	return res
}

// SetMessage 设置错误信息
func (res *Result) SetMessage(msg string) *Result {
	res.Message = msg
	return res
}

// SetData 设置返回数据
func (res *Result) SetData(data interface{}) *Result {
	res.Data = data
	return res
}
