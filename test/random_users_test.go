package test

import (
	"testing"

	"github.com/nharu-0630/twitter-api-client/example"
	"github.com/nharu-0630/twitter-api-client/tools"
)

func TestRandomUsers(t *testing.T) {
	tools.LoadEnv()
	cmd := example.NewRandomUsersCmdPoliticalParty("jimin_koho")
	cmd.Execute()
}
