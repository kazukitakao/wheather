package geocoding

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// baseAPIURL
// https://maps.googleapis.com/maps/api/geocode/json?parameters

// クエリ条件
// 必須 住所(単純に文字で打つか、郵便番号等で打つか選択可能)
// key APIkeyが必要

// 緯度 経度
// longitude latitude
// レスポンス geometry/location内の値を取得か？？

type geocodeResult struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID  string `json:"place_id"`
		PlusCode struct {
			CompoundCode string `json:"compound_code"`
			GlobalCode   string `json:"global_code"`
		} `json:"plus_code"`
		Types []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

// リクエスト機能の作成 keyはそのまま？md5とか必要？
// 送信して受け取ったレスポンスをjsonに変換する
// jsonから位置情報を取得
// APIコールの上限値をサーバ側で記録する必要あり

// API keyは後に環境変数に組み込む
const apiKey string = "AIzaSyDzrywihv-2Ii7xDM8UQD5OPlLEm_Xzs8c"

// GetGeocoding google APIを使用して指定した住所の座標情報を取得する
func GetGeocoding(address string) (longitude, latitude float64, location string) {
	// APIを呼び出して結果を取得
	url := "https://maps.googleapis.com/maps/api/geocode/json?address=" + address + "&key=" + apiKey
	fmt.Println(url)
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Google Geocode API failed : %s", err)
	}
	defer response.Body.Close()

	var data geocodeResult
	// dataにJSONに加工したデータをセット
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		log.Fatalf("JSON decorde failed : %s", err)
	}

	return data.Results[0].Geometry.Location.Lng, data.Results[0].Geometry.Location.Lat, data.Results[0].FormattedAddress
}
