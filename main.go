package main

import (
	"github.com/nharu-0630/twitter-api-client/example"
	"github.com/nharu-0630/twitter-api-client/tools"
)

func main() {
	tools.LoadEnv()
	cmd := example.NewRandomUsersCmdPoliticalParty("jimin_koho")
	cmd.Execute()
}
