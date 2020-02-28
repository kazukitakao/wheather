package main

import (
	"fmt"
	"wheather/geocoding"
)

func main() {

	// apiテスト
	lng, lat := geocoding.GetGeocoding("東京")
	fmt.Printf("緯度：%v 経度:%v", lng, lat)

}
