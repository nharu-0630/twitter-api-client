package main

import (
	"flag"

	"github.com/nharu-0630/twitter-api-client/cmd"
	"github.com/nharu-0630/twitter-api-client/tools"
	"go.uber.org/zap"
)

func main() {
	tools.LoadEnv()
	tools.LoadLogger()
	// switch cmd with flag
	cmdType := flag.String("cmd", "", "クライアントコマンドの種類 (必須)")
	flag.Parse()
	if *cmdType == "" {
		zap.L().Fatal("Cmd type is required")
	}
	switch *cmdType {
	case "UsersFollowings":
		cmd := cmd.UsersFollowingsCmdFromFlag()
		cmd.Execute()
	case "UsersFollowers":
		cmd := cmd.UsersFollowersCmdFromFlag()
		cmd.Execute()
	}
}
