package main

import (
	"context"
	"fmt"
)

type ExampleService struct{}

func (s *ExampleService) SayHello(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("Hello, %s!", name), nil
}

func main() {

}
