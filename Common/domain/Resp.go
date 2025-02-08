package domain

import (
	"encoding/json"
)

// Resp 统一返回对象
type Resp struct {
	Code int32       `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// NewResp 创建一个新的 Resp 对象
func NewResp(code int32, data interface{}, msg string) *Resp {
	return &Resp{
		Code: code,
		Data: data,
		Msg:  msg,
	}
}

// Success 创建一个成功的响应
func SuccessWithCode(code int32, data interface{}) *Resp {
	return NewResp(code, data, "成功")
}
func Success(data interface{}) *Resp {
	return NewResp(200, data, "成功")
}

// Error 创建一个错误的响应
func Error(code int32, msg string) *Resp {
	return NewResp(code, nil, msg)
}
func ErrorWithData(code int32, data interface{}, msg string) *Resp {
	return NewResp(code, data, msg)
}

// ToJSON 将 Resp 转换为 JSON 字符串
func (r *Resp) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
