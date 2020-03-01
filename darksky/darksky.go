package darksky

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// 天気情報のAPIとしてDark Sky APIを採用
// 1日1000回までのリクエストは無料
// “Powered by Dark Sky” that links to https://darksky.net/poweredby/とメッセージを表記する必要あり

type darkSkyResult struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Currently struct {
		Time                 int64   `json:"time"`
		Summary              string  `json:"summary"`
		Icon                 string  `json:"icon"`
		NearestStormDistance int     `json:"nearestStormDistance"`
		NearestStormBearing  int     `json:"nearestStormBearing"`
		PrecipIntensity      float64 `json:"precipIntensity"`
		PrecipProbability    float64 `json:"precipProbability"`
		Temperature          float64 `json:"temperature"`
		ApparentTemperature  float64 `json:"apparentTemperature"`
		DewPoint             float64 `json:"dewPoint"`
		Humidity             float64 `json:"humidity"`
		Pressure             float64 `json:"pressure"`
		WindSpeed            float64 `json:"windSpeed"`
		WindGust             float64 `json:"windGust"`
		WindBearing          int     `json:"windBearing"`
		CloudCover           float64 `json:"cloudCover"`
		UvIndex              int     `json:"uvIndex"`
		Visibility           float64 `json:"visibility"`
		Ozone                float64 `json:"ozone"`
	} `json:"currently"`
	Minutely struct {
		Summary string `json:"summary"`
		Icon    string `json:"icon"`
		Data    []struct {
			Time              int64   `json:"time"`
			PrecipIntensity   float64 `json:"precipIntensity"`
			PrecipProbability float64 `json:"precipProbability"`
		} `json:"data"`
	} `json:"minutely"`
	Hourly struct {
		Summary string `json:"summary"`
		Icon    string `json:"icon"`
		Data    []struct {
			Time                int64   `json:"time"`                         // 時刻
			Summary             string  `json:"summary"`                      // 概要
			Icon                string  `json:"icon"`                         // アイコン
			PrecipIntensity     float64 `json:"precipIntensity"`              // 降水強度 ある一定時間に降った雨が1時間降り続いたと換算した量
			PrecipProbability   float64 `json:"precipProbability"`            // 降水確率
			PrecipType          string  `json:"precipType,omitempty"`         // 降水型
			PrecipAccumulation  float64 `json:"precipAccumulation,omitempty"` // 降水量 ある時間内で降った雨が地面等水平な地面に溜まった量
			Temperature         float64 `json:"temperature"`                  // 温度
			ApparentTemperature float64 `json:"apparentTemperature"`          // 体感温度
			DewPoint            float64 `json:"dewPoint"`                     // 露点
			Humidity            float64 `json:"humidity"`                     // 湿度
			Pressure            float64 `json:"pressure"`                     // 圧力
			WindSpeed           float64 `json:"windSpeed"`                    // 風速
			WindGust            float64 `json:"windGust"`                     // 突風
			WindBearing         float64 `json:"windBearing"`                  // 風向き？？
			CloudCover          float64 `json:"cloudCover"`                   // 雲量
			UvIndex             int     `json:"uvIndex"`                      // 紫外線の強さ
			Visibility          float64 `json:"visibility"`                   // 視界
			Ozone               float64 `json:"ozone"`                        // 光化学スモッグ
		} `json:"data"`
	} `json:"hourly"`
	Daily struct {
		Summary string `json:"summary"`
		Icon    string `json:"icon"`
		Data    []struct {
			Time                        int64   `json:"time"`
			Summary                     string  `json:"summary"`
			Icon                        string  `json:"icon"`
			SunriseTime                 int64   `json:"sunriseTime"`
			SunsetTime                  int64   `json:"sunsetTime"`
			MoonPhase                   float64 `json:"moonPhase"`
			PrecipIntensity             float64 `json:"precipIntensity"`
			PrecipIntensityMax          float64 `json:"precipIntensityMax"`
			PrecipIntensityMaxTime      int64   `json:"precipIntensityMaxTime"`
			PrecipProbability           float64 `json:"precipProbability"`
			PrecipType                  string  `json:"precipType"`
			PrecipAccumulation          float64 `json:"precipAccumulation,omitempty"`
			TemperatureHigh             float64 `json:"temperatureHigh"`
			TemperatureHighTime         int64   `json:"temperatureHighTime"`
			TemperatureLow              float64 `json:"temperatureLow"`
			TemperatureLowTime          int64   `json:"temperatureLowTime"`
			ApparentTemperatureHigh     float64 `json:"apparentTemperatureHigh"`
			ApparentTemperatureHighTime int64   `json:"apparentTemperatureHighTime"`
			ApparentTemperatureLow      float64 `json:"apparentTemperatureLow"`
			ApparentTemperatureLowTime  int64   `json:"apparentTemperatureLowTime"`
			DewPoint                    float64 `json:"dewPoint"`
			Humidity                    float64 `json:"humidity"`
			Pressure                    float64 `json:"pressure"`
			WindSpeed                   float64 `json:"windSpeed"`
			WindGust                    float64 `json:"windGust"`
			WindGustTime                int64   `json:"windGustTime"`
			WindBearing                 float64 `json:"windBearing"`
			CloudCover                  float64 `json:"cloudCover"`
			UvIndex                     int     `json:"uvIndex"`
			UvIndexTime                 int64   `json:"uvIndexTime"`
			Visibility                  float64 `json:"visibility"`
			Ozone                       float64 `json:"ozone"`
			TemperatureMin              float64 `json:"temperatureMin"`
			TemperatureMinTime          int     `json:"temperatureMinTime"`
			TemperatureMax              float64 `json:"temperatureMax"`
			TemperatureMaxTime          int64   `json:"temperatureMaxTime"`
			ApparentTemperatureMin      float64 `json:"apparentTemperatureMin"`
			ApparentTemperatureMinTime  int64   `json:"apparentTemperatureMinTime"`
			ApparentTemperatureMax      float64 `json:"apparentTemperatureMax"`
			ApparentTemperatureMaxTime  int64   `json:"apparentTemperatureMaxTime"`
		} `json:"data"`
	} `json:"daily"`
	Flags struct {
		Sources        []string `json:"sources"`
		NearestStation float64  `json:"nearest-station"`
		Units          string   `json:"units"`
	} `json:"flags"`
	Offset int `json:"offset"`
}

const key string = "5a7d1f43ddeb56aabfec6ada3384968b"

// GetForecastToday darksky apiを使用して検索した時点のその日の天気情報を取得する
func GetForecastToday(latitude, longitude float64) {
	// 実行された日の1時間単位の天気予報を取得する
	// Hourlyから検索された日時の24時間分のデータを取得する
	now := time.Now().Unix()
	url := fmt.Sprintf("https://api.darksky.net/forecast/%s/%v,%v,%v?exclude=minutely,daily,alerts,flags&lang=ja", key, latitude, longitude, now)
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("DarkSky API failed : %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatalf("DarkSky API response failed : %v", response.StatusCode)
	}

	var data darkSkyResult
	// dataにJSONに加工したデータをセット
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		log.Fatalf("JSON decorde failed : %s", err)
	}

	//時間指定するには現在の日付0:00を取得してtimeオプションを付与する必要がある
	fmt.Println()
	fmt.Printf("%v年%v月%v日の天気：%v\n", data.Hourly.Summary)
	for index, v := range data.Hourly.Data {
		// 3時間毎の天気データを出力
		if index%3 == 0 {
			fmt.Printf("時間：%v ", time.Unix(v.Time, 0))
			fmt.Printf("天気：%v ", v.Summary)
			fmt.Printf("降水確率：%v ", v.PrecipProbability*100)
			fmt.Printf("気温：%v ", v.Temperature)
			fmt.Printf("湿度：%v\n", v.Humidity)
		}
	}
}
