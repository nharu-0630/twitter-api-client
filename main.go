package main

import (
	"flag"
	"math"

	"github.com/nharu-0630/twitter-api-client/cmd"
	"github.com/nharu-0630/twitter-api-client/tools"
)

func main() {
	tools.LoadEnv()
	tools.LoadLogger()
	seedScreenName := flag.String("from", "", "シードとなるユーザーのscreen_name (必須)")
	maxFollowersRequest := flag.Int("req", math.MaxInt, "1ユーザあたりの最大フォロワーリクエスト数 指定しない場合は全てのフォロワーを取得")
	maxChildRequest := flag.Int("depth", 1, "シードとなるユーザからの最大深度 指定しない場合は1(シードとなるユーザのフォロワーのみ取得)")
	maxUserLimit := flag.Int("limit", math.MaxInt, "取得するユーザー数の上限 指定しない場合は全てのユーザを取得")
	retryOnGuestFail := flag.Bool("retry", false, "ゲストトークンでのリクエスト失敗時に認証済みトークンでリトライする")
	statusUpdateSec := flag.Int("watch", 600, "ステータスを更新する間隔(秒) 指定しない場合は10分ごとに更新")
	flag.Parse()
	props := cmd.RandomUsersProps{
		SeedScreenName:      *seedScreenName,
		MaxFollowersRequest: *maxFollowersRequest,
		MaxChildRequest:     *maxChildRequest,
		MaxUserLimit:        *maxUserLimit,
		RetryOnGuestFail:    *retryOnGuestFail,
		StatusUpdateSec:     *statusUpdateSec,
	}
	cmd := cmd.RandomUsersCmd{
		Props: props,
	}
	cmd.Execute()
}
