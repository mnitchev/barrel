package main

import (
	"os"

	"github.com/mnitchev/barrel/runner"
)

func main() {
	command := os.Args[1]
	arguments := os.Args[2:]
	if err := runner.Run(command, arguments); err != nil {
		panic(err)
	}
}
