package common

import (
	"fmt"
	"testing"
)

func TestGoDotEnvVariable(t *testing.T) {
	str  ,err:= GoDotEnvVariable("SERVER_MONGODB")
	fmt.Println(str,err)
}
