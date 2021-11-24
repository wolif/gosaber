package password

import (
	"fmt"
	"testing"
)

func TestGenPwd(t *testing.T) {
	fmt.Println(GenPwd(10, false))
	fmt.Println(GenPwd(10, true))
}
