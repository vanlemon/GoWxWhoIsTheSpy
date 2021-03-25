package util

import (
	"fmt"
	"testing"
)

func TestGetFilePath(t *testing.T) {
	fmt.Println(GetFilePath())
}

func TestGetExecPath(t *testing.T) {
	fmt.Println(GetExecPath())
}
