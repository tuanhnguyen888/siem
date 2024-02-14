package rules

import (
	"fmt"
	"testing"
)

func TestIsMaliciousIP(t *testing.T) {
	b,alert, level := IsMaliciousIP("8.8.8.8")
	if b {
		fmt.Println("---------- ALERT -------------" , alert, level)
	}
}
