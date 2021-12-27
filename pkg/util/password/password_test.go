package password

import (
	"fmt"
	"testing"
)

func TestGenPwd(t *testing.T) {
	fmt.Println(Generate(10, false))
	fmt.Println(Generate(10, true))
}
