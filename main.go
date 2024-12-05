package main

import (
	"flag"
	"os"
	"reflect"

	"github.com/nharu-0630/twitter-api-client/cmd"
	"github.com/nharu-0630/twitter-api-client/model"
	"github.com/nharu-0630/twitter-api-client/util"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var (
	version  string
	revision string
	build    string
)

func main() {
	util.LoadEnv()
	util.SetZapGlobals()

	var (
		configPath  = flag.String("c", "./default.yml", "設定ファイルのパス")
		showVersion = flag.Bool("v", false, "バージョンを表示")
		showHelp    = flag.Bool("h", false, "ヘルプを表示")
	)
	flag.Parse()

	if *showVersion {
		println("version:", version)
		println("revision:", revision)
		println("build:", build)
		os.Exit(0)
	}

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	data, err := os.ReadFile(*configPath)
	if err != nil {
		zap.L().Fatal("設定ファイルの読み込みに失敗しました", zap.Error(err))
	}
	var config model.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		zap.L().Fatal("設定ファイルのパースに失敗しました", zap.Error(err))
	}

	cmdMap := map[string]interface{}{
		"brute_politician": &cmd.BrutePolitician{},
		"brute_query":      &cmd.BruteQuery{},
		"politician":       &cmd.PoliticianCmd{},
		"user_favorite":    &cmd.UserFavoriteCmd{},
		"user_followers":   &cmd.UserFollowersCmd{},
		"user_followings":  &cmd.UserFollowingsCmd{},
		"user_ids":         &cmd.UserIDsCmd{},
	}

	instance, exists := cmdMap[config.Type]
	if !exists {
		zap.L().Fatal("指定されたコマンドは存在しません")
	}

	props := reflect.ValueOf(instance).Elem().FieldByName("Props").Addr().Interface()
	err = yaml.Unmarshal(data, props)
	if err != nil {
		zap.L().Fatal("設定ファイルのパースに失敗しました", zap.Error(err))
	}

	reflect.ValueOf(instance).MethodByName("Execute").Call(nil)
}
