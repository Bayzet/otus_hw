package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		foo := `   foo
with new line`
		wait := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: false},
			"FOO":   EnvValue{Value: foo, NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		got, err := ReadDir("testdata/env")
		require.Nil(t, err)
		require.Equal(t, wait, got)
	})

	t.Run("dir is not exist", func(t *testing.T) {
		envs, err := ReadDir("/qwerty")
		require.Nil(t, envs)
		require.Truef(t, os.IsNotExist(err), "actual err - %v", err)
	})
}
