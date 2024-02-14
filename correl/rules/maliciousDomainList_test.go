package rules

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestIsMaliciousURL(t *testing.T) {

	a,b,c :=  IsMaliciousURL("https://www.programiz.com/golang/online-compiler/")
	if a {
		logrus.Println( a,b,c)
	}
}