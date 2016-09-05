package conf

import (
	"errors"
	"os"
	"testing"

	"github.com/takawang/sugar"
)

var (
	ErrFilterNotFound = errors.New("Filter not found")
)

func TestLogger(t *testing.T) {
	s := sugar.New(t)

	s.Assert("Test setLogger", func(_ sugar.Log) bool {
		os.Setenv(envBackendEndpoint, "123")
		base.initConfig()
		SetDefault(keyLogEnableDebug, defaultLogEnableDebug)
		Set(keyLogToJSONFormat, true)
		Set(keyLogEnableDebug, false)
		base.setLogger()
		Set(keyLogFileName, "/tmp/abc")
		Set(keyLogToFile, true)
		base.setLogger()
		return true
	})

	s.Assert("Test Init logger", func(_ sugar.Log) bool {
		os.Setenv(envConfWeb, "")
		os.Setenv(envBackendEndpoint, "")
		base.initConfig()
		os.Setenv(envConfWeb, "a")
		base.initConfig()
		os.Setenv(envConfWeb, "a")
		os.Setenv(envBackendEndpoint, "b")
		base.initConfig()
		os.Setenv(envConfWeb, "a")
		os.Setenv(envBackendEndpoint, "")
		base.initConfig()
		return true
	})

	s.Assert("Test Fail cases", func(_ sugar.Log) bool {
		os.Setenv(envBackendEndpoint, "123")
		base.initConfig()
		SetDefault(keyLogEnableDebug, defaultLogEnableDebug)
		Set(keyLogToJSONFormat, true)
		Set(keyLogEnableDebug, false)
		base.setLogger()
		Set(keyLogFileName, "/proc/111")
		Set(keyLogToFile, true)
		base.setLogger()
		return true
	})
}
