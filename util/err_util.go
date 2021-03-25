package util

import (
	"errors"
	"fmt"
)

// 新建错误：字符串格式化
func NewErrf(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args...))
}
