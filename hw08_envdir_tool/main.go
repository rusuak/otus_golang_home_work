package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		os.Stdout.WriteString("Please pass the directory path and the command to execute\n")

		return
	}

	dir, err := ReadDir(args[1])
	if err != nil {
		os.Stdout.WriteString(fmt.Sprintln(err))

		return
	}

	RunCmd(args[2:], dir)
}
