package test

import (
	"testing"
)

func TestGCPressure(t *testing.T) {
	for i := 0; i < 20; i++ {
		buf := make([][]byte, 100000)

		for j := 0; j < 100000; j++ {
			buf[j] = make([]byte, 1024)
		}

		_ = buf
	}
}
