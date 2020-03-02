package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/line/line-bot-sdk-go/linebot/httphandler"
)

const (
	chanelSecret = "5655997ac303fd6fa4d9887b4b1fbe97"
	channelToken = "QpilbXWCyEwgXkH/JXU6ZAE7xAE4isIe/l6bs9xtxUiyadt7jWC0KVTiERsM93+W9c4ER31S4bliFjvMHxJDfzBFhkbtzDzbEIDTCneEp2QEZwQKNqk4KMVW9Ub3GpPp1baNAikgvnpHuzViYb7wVQdB04t89/1O/w1cDnyilFU="
)

func main() {

	handler, err := httphandler.New(
		chanelSecret,
		channelToken,
		// 今回はお試しなので直接定数をセット
		// os.Getenv("CHANNEL_SECRET"),
		// os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		bot, err := handler.NewClient()
		if err != nil {
			log.Print(err)
			return
		}

		for _, event := range events {
			switch event.Type {
			// メッセージ送信イベント
			case linebot.EventTypeMessage:
				if message, ok := event.Message.(*linebot.TextMessage); ok {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
					fmt.Println("オウム返ししました")
				}
			// フォローイベント
			case linebot.EventTypeFollow:
				fmt.Println("新しくフォローされました")
			// フォロー解除イベント
			case linebot.EventTypeUnfollow:
				fmt.Println("フォロー解除されました")
			default:
				fmt.Printf("そのイベントには対応していません:%v", event.Type)
			}
		}
	})

	http.Handle("/callback", handler)
	// For actually use, you must support HTTPS by using `ListenAndServeTLS`, reverse proxy or etc.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	// 署名検証処理
	// LINEから送られ処理するイベント メッセージイベント(テキストのみ)、フォローイベント、フォロー解除イベント
	// フォロー、解除の時は取得したユーザ情報をDBに登録
	// メッセージの時は取得した文字列を天気APIに渡し、情報を検索してクライアントに返す
}
