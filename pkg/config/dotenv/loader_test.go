package dotenv

import (
	"os"
	"testing"
)

func TestLoader(t *testing.T) {
	err := Load(".env")
	if err != nil {
		t.Error(err)
	}
	t.Log(os.Getenv("ZZ"))
	t.Log(os.Getenv("A_B"))
	t.Log(os.Getenv("B_C"))
}