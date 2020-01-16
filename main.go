package main

import (
	"bytes"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/go-echarts/go-echarts/charts"
)

var (
	initV float64 = 50000

	data = map[string]float64{
		"19/11/01": 50000,
		"19/11/02": 51000,
		"19/11/03": 48000,
		"19/11/04": 46000,
		"19/11/05": 43000,
		"19/11/06": 44000,
	}
)

func handler(w http.ResponseWriter, _ *http.Request) {
	var (
		nameItems []string
		foodItems []float64
		keySlice  sort.StringSlice
	)

	for k := range data {
		keySlice = append(keySlice, k)
	}
	sort.Sort(keySlice)
	for _, k := range keySlice {
		nameItems = append(nameItems, k)
		foodItems = append(foodItems, data[k]-initV)
	}

	line := charts.NewLine()
	line.SetGlobalOptions(charts.TitleOpts{Title: "商品期货"})
	line.AddXAxis(nameItems).AddYAxis("盈亏图", foodItems,
		charts.LabelTextOpts{Show: true},
		charts.AreaStyleOpts{Opacity: 0.2},
		charts.LineOpts{Smooth: true},
	)

	var bufferW bytes.Buffer
	line.Render(&bufferW)
	htmlV := bufferW.String()
	htmlV = strings.Replace(htmlV, "Awesome go-echarts", "逐日结算", 1)

	w.Write([]byte(htmlV))
}

func main() {
	log.Println("Start !")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
