package util

import (
	"context"
	"fmt"
)

// gin 上下文中保存 ctx_key 和大数据 data_key
// ctx 上下文中保存小数据

const (
	MethodKey  = "method_key"  // ctx 中的请求接口唯一标识
	RequestKey = "request_key" // ctx 中的请求参数唯一标识
)

func CtxWithMethodKey(ctx context.Context, methodKey string) context.Context {
	return context.WithValue(ctx, MethodKey, methodKey)
}

func CtxWithRequestKey(ctx context.Context, requestKey string) context.Context {
	return context.WithValue(ctx, RequestKey, requestKey)
}

func CtxGetString(ctx context.Context, key string) string {
	value := ctx.Value(key)
	if value == nil {
		panic(fmt.Sprintf("ctx no value of key %s", key))
	}
	valueString, ok := value.(string)
	if !ok {
		panic(fmt.Sprintf("ctx no string value of key %s, value %#v", key, value))
	}
	return valueString
}
