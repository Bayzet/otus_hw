package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {

	fromPath := "testdata/input.txt"
	toPath := "out.txt"
	t.Run("limit be negative", func(t *testing.T) {
		err := Copy(fromPath, toPath, 0, -1)
		require.Truef(t, errors.Is(err, ErrLimitCannotBeNegative), "actual err - %v", err)
	})

	t.Run("unsupported file", func(t *testing.T) {
		err := Copy("/dev/urandom", toPath, 0, 0)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	})

	t.Run("no such file or directory", func(t *testing.T) {
		err := Copy("", toPath, 0, 0)
		require.Truef(t, os.IsNotExist(err), "actual err - %v", err)

		err = os.Remove(toPath)
		require.Nil(t, err)
	})

	t.Run("offset exceeds fileSize", func(t *testing.T) {
		err := Copy(fromPath, toPath, 7000, 0)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)

		err = os.Remove(toPath)

		require.Nil(t, err)
	})
}
