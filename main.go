package main

import (
	"bytes"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/go-echarts/go-echarts/charts"
)

func handler(w http.ResponseWriter, _ *http.Request) {
	var (
		nameItems []string
		foodItems []int
		keySlice  sort.StringSlice
	)

	// 日期排序
	for k := range data {
		keySlice = append(keySlice, k)
	}
	sort.Sort(keySlice)

	// 计算盈亏
	for _, k := range keySlice {
		nameItems = append(nameItems, k)
		var initV = 0.0
		if k >= "19/11/07" && k < "19/11/11" {
			initV = 20000
		} else if k >= "19/11/11" && k < "19/11/13" {
			initV = 25000
		} else if k >= "19/11/13" && k < "20/01/16" {
			initV = 50000
		} else if k >= "20/01/16" && k < "20/01/17" {
			initV = 70000
		} else {
			initV = 80000
		}

		foodItems = append(foodItems, int(data[k]-initV))
	}

	// 画图表
	line := charts.NewLine()
	line.SetGlobalOptions(charts.TitleOpts{Title: "商品期货"})
	line.AddXAxis(nameItems).AddYAxis("盈亏图", foodItems,
		// charts.LabelTextOpts{Show: true},
		charts.AreaStyleOpts{Opacity: 0.2},
		charts.LineOpts{Smooth: true},
	)

	// 替换标题
	var bufferW bytes.Buffer
	line.Render(&bufferW)
	htmlV := bufferW.String()
	htmlV = strings.Replace(htmlV, "Awesome go-echarts", "逐日结算", 1)

	// 输出html
	w.Write([]byte(htmlV))
}

func main() {
	log.Println("Start !")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
