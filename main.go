package main

import (
	"flag"

	"github.com/nharu-0630/twitter-api-client/cmd"
	"github.com/nharu-0630/twitter-api-client/tools"
	"go.uber.org/zap"
)

func main() {
	tools.LoadEnv()
	tools.SetZapGlobals()
	cmdType := flag.String("cmd", "", "クライアントコマンドの種類 (必須)")
	flag.Parse()
	if *cmdType == "" {
		zap.L().Fatal("Cmd type is required")
	}
	switch *cmdType {
	case "UserFollowers":
		cmd := cmd.UserFollowersCmdFromFlag()
		cmd.Execute()
	case "UserID":
		cmd := cmd.UserIDCmdFromFlag()
		cmd.Execute()
	case "UserFollowings":
		cmd := cmd.UserFollowingsCmdFromFlag()
		cmd.Execute()
	default:
		zap.L().Fatal("Cmd type is required")
	}
}
