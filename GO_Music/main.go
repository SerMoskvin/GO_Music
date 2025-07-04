package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func main() {
	v := validator.New()
	fmt.Println("Validator created:", v)
}
