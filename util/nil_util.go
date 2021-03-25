package util

// 判断变量是否为空
func IsNil(i interface{}) bool {
	return i == nil || i == 0 || i == "" || i == false
}

// 判断列表是否有空，且列表本身不能为空
func HasNil(is ...interface{}) bool {
	if len(is) == 0 { // 列表本身不能为空
		return true
	}
	for _, i := range is { // 遍历列表判断为空
		if IsNil(i) {
			return true
		}
	}
	return false
}
