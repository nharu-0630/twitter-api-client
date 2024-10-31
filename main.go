package main

import (
	"flag"
	"os"
	"path/filepath"
	"reflect"

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

	cmdMap := map[string]interface{}{
		"Politician":     &cmd.PoliticianCmd{},
		"UserFavorite":   &cmd.UserFavoriteCmd{},
		"UserFollowers":  &cmd.UserFollowersCmd{},
		"UserFollowings": &cmd.UserFollowingsCmd{},
		"UserIDs":        &cmd.UserIDsCmd{},
	}

	cmdInstance, exists := cmdMap[config.CmdType]
	if !exists {
		zap.L().Fatal("Unknown CmdType", zap.String("CmdType", config.CmdType))
	}

	propsField := reflect.ValueOf(cmdInstance).Elem().FieldByName("Props").Addr().Interface()
	err = yaml.Unmarshal(configFile, propsField)
	if err != nil {
		zap.L().Fatal("Failed to unmarshal configuration", zap.Error(err))
	}

	reflect.ValueOf(cmdInstance).MethodByName("Execute").Call(nil)
}
