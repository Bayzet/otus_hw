package logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/consts"

	"github.com/stretchr/testify/require"
)

func TestLogger_New(t *testing.T) {
	tests := []struct {
		level    string
		filePath string
	}{
		{
			consts.LogLevelInfo,
			"log.txt",
		},
	}

	for _, tt := range tests {
		logger, err := New(tt.level, tt.filePath)
		require.NoError(t, err)
		require.Equal(t, logger.level, LogLevel(tt.level))
		require.FileExists(t, tt.filePath)
		_ = os.Remove(tt.filePath)
	}

}

func TestLogger_Logging(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2022, time.April, 28, 0, 0, 0, 0, time.UTC)
	})

	filePath := "log.txt"

	tests := []struct {
		level string
		exp   string
	}{
		{
			consts.LogLevelInfo,
			`{"time":"2022-04-28T00:00:00Z","level":"info","message":"some text"}
{"time":"2022-04-28T00:00:00Z","level":"warning","message":"some text"}
{"time":"2022-04-28T00:00:00Z","level":"error","message":"some text"}
`,
		},
		{
			consts.LogLevelWarn,
			`{"time":"2022-04-28T00:00:00Z","level":"info","message":"some text"}
{"time":"2022-04-28T00:00:00Z","level":"warning","message":"some text"}
{"time":"2022-04-28T00:00:00Z","level":"error","message":"some text"}
`,
		},
		{
			consts.LogLevelError,
			`{"time":"2022-04-28T00:00:00Z","level":"error","message":"some text"}
`,
		},
		{
			consts.LogLevelDebug,
			`{"time":"2022-04-28T00:00:00Z","level":"info","message":"some text"}
{"time":"2022-04-28T00:00:00Z","level":"warning","message":"some text"}
{"time":"2022-04-28T00:00:00Z","level":"error","message":"some text"}
{"time":"2022-04-28T00:00:00Z","level":"debug","message":"some text"}
`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
			logger, err := New(tt.level, filePath)
			require.NoError(t, err)

			logger.Info("some text")
			logger.Warn("some text")
			logger.Error("some text")
			logger.Debug("some text")
			logger.Close()

			b, _ := ioutil.ReadFile(filePath)
			require.Equal(t, tt.exp, string(b))

			_ = os.Remove(filePath)
		})
	}

	defer func() {
		_ = os.Remove(filePath)
	}()
}
