package util

import (
	"encoding/json"
)

// 结构体转化为 Interface Map
func StructToInterfaceMap(v interface{}) map[string]interface{} {
	var ret map[string]interface{}
	bytes, _ := json.Marshal(v)     // TODO 忽略 json.Marshal 报错
	_ = json.Unmarshal(bytes, &ret) // TODO 忽略 json.Unmarshal 报错
	return ret
}

// 结构体转化为 Json String
func StructToJsonString(v interface{}) string {
	jsonByted, _ := json.Marshal(v) // TODO 忽略 json.Marshal 报错
	return string(jsonByted)
}

// Interface Map 转化为 Json String
func InterfaceMapToJsonString(m map[string]interface{}) string {
	mapString := make(map[string]string)
	for k, v := range m {
		if vString, ok := v.(string); ok {
			mapString[k] = vString
		} else {
			mapString[k] = StructToJsonString(v)
		}
	}
	return StructToJsonString(mapString)
}

// Json String 转化为 String Map
func JsonStringToStringMap(s string) map[string]string {
	m := make(map[string]string)
	_ = json.Unmarshal([]byte(s), &m) // TODO 忽略 json.Unmarshal 报错
	return m
}

// Json String 转化为 Interface Map
func JsonStringToInterfaceMap(s string) map[string]interface{} {
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(s), &m) // TODO 忽略 json.Unmarshal 报错
	return m
}
