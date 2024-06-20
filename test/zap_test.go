package test

import (
	"testing"

	"github.com/nharu-0630/twitter-api-client/tools"
	"go.uber.org/zap"
)

func TestZap(t *testing.T) {
	tools.LoadEnv()
	tools.LoadLogger()
	zap.L().Warn("Get users", zap.String("UserCount", "3"))
}
