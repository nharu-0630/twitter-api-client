package test

import (
	"testing"

	"github.com/nharu-0630/twitter-api-client/example"
	"github.com/nharu-0630/twitter-api-client/tools"
)

func TestRandomUsers(t *testing.T) {
	tools.LoadEnv()
	tools.LoadLogger()
	cmd := example.NewRandomUsersCmdHead("nharu_0630")
	cmd.Execute()
}
