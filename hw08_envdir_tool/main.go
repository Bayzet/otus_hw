package main

import (
	"flag"
	"fmt"
)

// import flag "github.com/spf13/pflag"

func main() {
	flag.Parse()

	cliArgs := flag.Args()
	envDir := cliArgs[0]
	cmd := cliArgs[1:]

	envs, err := ReadDir(envDir)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	RunCmd(cmd, envs)
}
