package example

import (
	"context"
	"fmt"
)

type serv struct {
}

var Serv = new(serv)

func (s *serv) F1(ctx *context.Context, input *string, output *string) (string, error) {
	return "F1", nil
}

func (s *serv) F2(ctx *context.Context, input *string, output *string) (string, error) {
	return "F2", nil
}

func (s *serv) F3(ctx *context.Context, input *string, output *string) (string, error) {
	return "F3", nil
}

func (s *serv) F4(ctx *context.Context, input *string, output *string) (string, error) {
	fmt.Printf("get input: %s", *input)
	fmt.Println()
	return "F4", fmt.Errorf("f4 error")
}
