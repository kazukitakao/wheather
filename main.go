package main

import (
	"fmt"
	"wheather/darksky"
	"wheather/geocoding"
)

func main() {

	// apiテスト
	longitude, latitude, location := geocoding.GetGeocoding("善行")
	fmt.Printf("検索された住所:%v 緯度：%v 経度:%v\n", location, longitude, latitude)
	darksky.GetForecastToday(latitude, longitude)
}
