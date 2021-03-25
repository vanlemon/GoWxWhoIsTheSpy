package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Http Get 工具方法
// - urlString：待拼接 url
// - params：拼接参数
func HttpGet(urlString string, params ...interface{}) (string, error) {
	getUrl := fmt.Sprintf(urlString, params...) // 拼接 url

	resp, err := http.Get(getUrl) // 获取请求结果
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close() // TODO 忽略错误
	}()

	body, err := ioutil.ReadAll(resp.Body) // 读取请求结果
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Http Get 工具方法
// - urlString：url
// - data：数据
// - sign：签名
func HttpPost(urlString, reqDataJsonString, reqSign string) (string, error) {
	resp, err := http.PostForm(urlString, url.Values{"data": {reqDataJsonString}, "sign": {reqSign}})
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close() // TODO 忽略错误
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
