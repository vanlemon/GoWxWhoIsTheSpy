package util

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// 合并 ip 地址和端口号
func CombineIpAndPort(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

// 解析 ip 地址和端口号
func ParseIpAndPort(ipAndPort string) (string, int, error) {
	strs := strings.Split(ipAndPort, ":")
	if len(strs) == 2 {
		address := net.ParseIP(strs[0]) // 校检 ip 地址合法性
		if address == nil {
			return "", 0, errors.New("Parse Ip:" + strs[0] + " Error.")
		}
		ip := address.String()

		port, err := strconv.Atoi(strs[1]) // 端口转 int，未检验合法性
		if err != nil {
			return "", 0, err
		}
		return ip, port, nil
	}
	return "", 0, errors.New("Parse Ip Port:" + ipAndPort + " Error.")
}
