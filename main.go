package main

import (
	"github.com/nharu-0630/twitter-api-client/example"
	"github.com/nharu-0630/twitter-api-client/tools"
)

func main() {
	tools.LoadEnv()
	tools.SetZapGlobals()
	// cmdType := flag.String("cmd", "", "クライアントコマンドの種類 (必須)")
	// flag.Parse()
	// if *cmdType == "" {
	// 	zap.L().Fatal("Cmd type is required")
	// }
	// switch *cmdType {
	// case "UserFollowers":
	// 	cmd := cmd.UserFollowersCmdFromFlag()
	// 	cmd.Execute()
	// case "UserID":
	// 	cmd := cmd.UserIDCmdFromFlag()
	// 	cmd.Execute()
	// case "UserFollowings":
	// 	cmd := cmd.UserFollowingsCmdFromFlag()
	// 	cmd.Execute()
	// default:
	// 	zap.L().Fatal("Cmd type is required")
	// }
	cmd := example.NewGroupUsersCmd()
	cmd.Execute()
}
