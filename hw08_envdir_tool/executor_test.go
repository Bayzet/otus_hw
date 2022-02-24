package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("exit code -1", func(t *testing.T) {
		var env Environment
		cmd := []string{
			"qwerty",
			"-aaa",
		}

		code := RunCmd(cmd, env)
		require.Equal(t, -1, code)
	})

	t.Run("exit code 0", func(t *testing.T) {
		var env Environment
		cmd := []string{
			"ls",
			"-a",
		}

		code := RunCmd(cmd, env)
		require.Equal(t, 0, code)
	})

	t.Run("exit code 1", func(t *testing.T) {
		var env Environment
		cmd := []string{
			"cp",
		}

		code := RunCmd(cmd, env)
		require.Equal(t, 1, code)
	})
}
