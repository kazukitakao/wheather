package controllers

import (
	"fmt"
	"log"
	"net/http"
	"udemy/gotrading/config"

	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	chanelSecret = "5655997ac303fd6fa4d9887b4b1fbe97"
	channelToken = "QpilbXWCyEwgXkH/JXU6ZAE7xAE4isIe/l6bs9xtxUiyadt7jWC0KVTiERsM93+W9c4ER31S4bliFjvMHxJDfzBFhkbtzDzbEIDTCneEp2QEZwQKNqk4KMVW9Ub3GpPp1baNAikgvnpHuzViYb7wVQdB04t89/1O/w1cDnyilFU="
)

var bot *linebot.Client

func init() {
	client, err := linebot.New(
		chanelSecret,
		channelToken,
		// 今回はお試しなので直接定数をセット
		// os.Getenv("CHANNEL_SECRET"),
		// os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	} else {
		bot = client
	}
}

func lineEventHandler(w http.ResponseWriter, req *http.Request) {
	// ParseRequest内でヘッダーX-Line-Signatureを利用した署名の検証も行っている
	events, err := bot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		switch event.Type {
		// メッセージ送信イベント
		case linebot.EventTypeMessage:
			if message, ok := event.Message.(*linebot.TextMessage); ok {
				fmt.Println(message)
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
}

// StartWebServer webサーバ起動処理
func StartWebServer() error {
	http.HandleFunc("/callback", lineEventHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}
