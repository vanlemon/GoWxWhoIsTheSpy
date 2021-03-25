package util

import (
	"strconv"
	"strings"
	"time"
)

// 生成一个 logId 格式的 UUID，虽然保障了时间连续性，但是需要人工校检重复
// hexNum 为末位随机十六进制字符串的长度，20+hexNum 位
func GenUUID(hexNum int) string {
	t := time.Now()

	return strings.Join([]string{
		t.Format("20060102150405"),          // 14
		strconv.Itoa(t.Nanosecond() / 1000), // 6
		RandHexString(hexNum),
	}, "")
}
