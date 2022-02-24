package main

import (
	"flag"
	"log"
	"os"
)

// import flag "github.com/spf13/pflag"

func main() {
	flag.Parse()

	cliArgs := flag.Args()
	envDir := cliArgs[0]
	cmd := cliArgs[1:]

	envs, err := ReadDir(envDir)
	if err != nil {
		log.Fatal(err)
	}

	exitCode := RunCmd(cmd, envs)

	os.Exit(exitCode)
}
