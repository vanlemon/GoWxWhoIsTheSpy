package util

import (
	"context"
	"fmt"
	logs "lmf.mortal.com/GoLogs"
)

const (
	MessageSuccess      = "Success"      // 请求成功
	MessageError        = "Error"        // 请求失败
	MessageDuplicate    = "Duplicate"    // 请求重复
	MessageOverload     = "Overload"     // 请求超载
	MessageNoSign       = "NoSign"       // 请求无签名
	MessageSignError    = "SignError"    // 请求签名错误
	MessageNoData       = "NoData"       // 请求无参数
	MessageUserNotLogin = "UserNotLogin" // 请求用户未登录
	MessageRoomInvalid  = "RoomInvalid"  // 请求用户未关联有效房间

)

// 封装返回值对象
type HTTPResponse struct {
	Success bool        `json:"Success"` // 请求是否成功
	Message string      `json:"Message"` // 返回信息，上述枚举值的返回信息可作为错误码
	Data    interface{} `json:"Data"`    // 返回数据
	LogId   string      `json:"LogId"`  // 日志 logid
}

func SuccessWithData(ctx context.Context, data interface{}) map[string]interface{} {
	return StructToInterfaceMap(HTTPResponse{true, MessageSuccess, data, logs.CtxGetLogId(ctx)})
}

func Success(ctx context.Context) map[string]interface{} {
	return StructToInterfaceMap(HTTPResponse{true, MessageSuccess, nil, logs.CtxGetLogId(ctx)})
}

func ErrorWithMessage(ctx context.Context, format string, v ...interface{}) map[string]interface{} {
	return StructToInterfaceMap(HTTPResponse{false, fmt.Sprintf(format, v...), nil, logs.CtxGetLogId(ctx)})
}

func Error(ctx context.Context) map[string]interface{} {
	return StructToInterfaceMap(HTTPResponse{false, MessageError, nil, logs.CtxGetLogId(ctx)})
}

func Duplicate(ctx context.Context) map[string]interface{} {
	return StructToInterfaceMap(HTTPResponse{false, MessageDuplicate, nil, logs.CtxGetLogId(ctx)})
}

func Overload(ctx context.Context) map[string]interface{} {
	return StructToInterfaceMap(HTTPResponse{false, MessageOverload, nil, logs.CtxGetLogId(ctx)})
}

func NoSign(ctx context.Context) map[string]interface{} {
	return StructToInterfaceMap(HTTPResponse{false, MessageNoSign, nil, logs.CtxGetLogId(ctx)})
}

func SignError(ctx context.Context) map[string]interface{} {
	return StructToInterfaceMap(HTTPResponse{false, MessageSignError, nil, logs.CtxGetLogId(ctx)})
}

func NoData(ctx context.Context) map[string]interface{} {
	return StructToInterfaceMap(HTTPResponse{false, MessageNoData, nil, logs.CtxGetLogId(ctx)})
}

func UserNotLogin(ctx context.Context) map[string]interface{} {
	return StructToInterfaceMap(HTTPResponse{false, MessageUserNotLogin, nil, logs.CtxGetLogId(ctx)})
}

func RoomInvalid(ctx context.Context) map[string]interface{} {
	return StructToInterfaceMap(HTTPResponse{false, MessageRoomInvalid, nil, logs.CtxGetLogId(ctx)})
}
