package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmdExec := exec.Command(cmd[0], cmd[1:]...)

	cmdExec.Stdout = os.Stdout
	cmdExec.Stdin = os.Stdin
	cmdExec.Stderr = os.Stderr

	for k, e := range env {
		if e.NeedRemove {
			os.Unsetenv(k)
			break
		}

		if _, ok := os.LookupEnv(k); ok {
			os.Unsetenv(k)
		}

		os.Setenv(k, e.Value)
	}

	_ = cmdExec.Run()
	returnCode = cmdExec.ProcessState.ExitCode()

	return
}
