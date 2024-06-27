package test

import (
	"testing"

	"github.com/nharu-0630/twitter-api-client/tools"
	"go.uber.org/zap"
)

func TestZap(t *testing.T) {
	tools.LoadEnv()
	tools.SetZapGlobals()
	zap.L().Info("test zap", zap.String("test", "test"))
}
