package main

import (
	"github.com/vcsfrl/random-fit/cmd/shell"
)

//go:generate go run main.go exec generate-code
func main() {
	shell.New().Run()
}
