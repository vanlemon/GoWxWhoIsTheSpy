package util

import (
	"crypto/md5"
	"fmt"
	"lmf.mortal.com/GoWxWhoIsTheSpy/config"
)

// 对请求进行签名
func GetSign(reqDataJsonString string) string {
	salt := config.GetSalt()                        // 获取签名盐值
	data := []byte(salt + reqDataJsonString + salt) // 构造签名数据
	hash := md5.Sum(data)                           // 计算 hash []byte
	sign := fmt.Sprintf("%x", hash)                 // []byte 转 16 进制
	return sign
}
