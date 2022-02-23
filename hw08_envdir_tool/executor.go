package main

import (
	// "fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmdExec := exec.Command(cmd[0], cmd[1:]...)

	for k, e := range env {
		if e.NeedRemove {
			os.Unsetenv(k)
			break
		}

		if _, ok := os.LookupEnv(k); ok {
			os.Unsetenv(k)
		}
		// fmt.Println("=== === ENV: ", k+"="+e.Value)
		os.Setenv(k, e.Value)

		// cmdExec.Env = append(os.Environ(), k+"="+e.Value)
	}

	cmdExec.Stdout = os.Stdout
	cmdExec.Stdin = os.Stdin
	cmdExec.Stderr = os.Stderr

	if err := cmdExec.Run(); err != nil {
		returnCode = -1
	}

	return
}
