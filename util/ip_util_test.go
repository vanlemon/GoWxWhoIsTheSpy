package util

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestCombineIpAndPort(t *testing.T) {
	ip := "127.0.0.1"
	port := 80
	ip_port := CombineIpAndPort(ip, port)
	assert.Equal(t, ip_port, "127.0.0.1:80")
}

func TestParseIpAndPort(t *testing.T) {
	ip := "127.0.0.1"
	port := 80
	ip_t, port_t, err := ParseIpAndPort("127.0.0.1:80")
	assert.Equal(t, ip, ip_t)
	assert.Equal(t, port, port_t)
	assert.Equal(t, err, nil)
}

func TestParseIpAndPort_IpError(t *testing.T) {
	_, _, err := ParseIpAndPort("127.0.0.a:80")
	assert.NotEqual(t, err, nil)
}

func TestParseIpAndPort_PortError(t *testing.T) {
	_, _, err := ParseIpAndPort("127.0.0.1:ab")
	assert.NotEqual(t, err, nil)
}
