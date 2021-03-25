package util

import (
	"fmt"
	"testing"
)

func BenchmarkRandOctString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandHexString(10)
	}
}

func TestRandOctString(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(RandHexString(10))
	}
}
