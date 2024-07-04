package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/nharu-0630/twitter-api-client/cmd"
	"github.com/nharu-0630/twitter-api-client/tools"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func main() {
	tools.LoadEnv()
	tools.SetZapGlobals()
	flag.Parse()
	ymlFilePath := flag.Arg(0)
	absPath, err := filepath.Abs(ymlFilePath)
	if err != nil {
		zap.L().Fatal("Failed to get absolute path", zap.Error(err))
	}
	ymlFilePath = absPath

	configFile, err := os.ReadFile(ymlFilePath)
	if err != nil {
		zap.L().Fatal("Failed to read configuration file", zap.Error(err))
	}

	var config struct {
		CmdType string `yaml:"CmdType"`
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		zap.L().Fatal("Failed to unmarshal configuration", zap.Error(err))
	}

	switch config.CmdType {
	case "GroupUsers":
		var props cmd.GroupUsersProps
		err = yaml.Unmarshal(configFile, &props)
		if err != nil {
			zap.L().Fatal("Failed to unmarshal configuration", zap.Error(err))
		}
		cmd := cmd.GroupUsersCmd{
			Props: props,
		}
		cmd.Execute()
	case "UserFollowers":
		var props cmd.UserFollowersProps
		err = yaml.Unmarshal(configFile, &props)
		if err != nil {
			zap.L().Fatal("Failed to unmarshal configuration", zap.Error(err))
		}
		cmd := cmd.UserFollowersCmd{
			Props: props,
		}
		cmd.Execute()
	case "UserFollowings":
		var props cmd.UserFollowingsProps
		err = yaml.Unmarshal(configFile, &props)
		if err != nil {
			zap.L().Fatal("Failed to unmarshal configuration", zap.Error(err))
		}
		cmd := cmd.UserFollowingsCmd{
			Props: props,
		}
		cmd.Execute()
	case "UserID":
		var props cmd.UserIDProps
		err = yaml.Unmarshal(configFile, &props)
		if err != nil {
			zap.L().Fatal("Failed to unmarshal configuration", zap.Error(err))
		}
		cmd := cmd.UserIDCmd{
			Props: props,
		}
		cmd.Execute()
	default:
		zap.L().Fatal("Unknown CmdType", zap.String("CmdType", config.CmdType))
	}
}
