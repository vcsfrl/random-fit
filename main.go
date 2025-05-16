package main

import "github.com/vcsfrl/random-fit/cmd"
import _ "github.com/mkevac/debugcharts"

//go:generate go run main.go code generate
func main() {
	cmd.Execute()
}
