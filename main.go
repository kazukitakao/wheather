package main

import (
	"wheather/controllers"
)

// TextMessage LINEから送信されたテキストメッセージ
type TextMessage struct {
	ReplyToken string `json:"replyToken"`
	Type       string `json:"type"`
	Mode       string `json:"mode"`
	Timestamp  int64  `json:"timestamp"`
	Source     struct {
		Type   string `json:"type"`
		UserID string `json:"userId"`
	} `json:"source"`
	Message struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"message"`
}

func main() {

	controllers.StartWebServer()
	// フォロー、解除の時は取得したユーザ情報をDBに登録
	// メッセージの時は取得した文字列を天気APIに渡し、情報を検索してクライアントに返す
}
