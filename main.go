package main

import "github.com/vcsfrl/random-fit/cmd"

//go:generate go run main.go exec generate-code
func main() {
	cmd.Execute()
}
