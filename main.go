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
	data = map[string]float64{
		"19/11/07": 20000, "19/11/08": 19978.84,
		"19/11/11": 25148.45, "19/11/12": 25656.71, "19/11/13": 49737.71, "19/11/14": 48738.16, "19/11/15": 50693.73,
		"19/11/18": 51000.16, "19/11/19": 47732.97, "19/11/20": 47241.63, "19/11/21": 46544.75, "19/11/22": 46601.91,
		"19/11/25": 47246.89, "19/11/26": 46773.16, "19/11/27": 46713.56, "19/11/28": 46468.27, "19/11/29": 46755.27,
		"19/12/02": 47460.27, "19/12/03": 47540.03, "19/12/04": 47192.53, "19/12/05": 47014.43, "19/12/06": 47246.93,
		"19/12/09": 47871.96, "19/12/10": 47779.71, "19/12/11": 48842.48, "19/12/12": 49496.87, "19/12/13": 49031.63,
		"19/12/16": 50478.87, "19/12/17": 49050.97, "19/12/18": 49140.97, "19/12/19": 49194.50, "19/12/20": 47341.28,
		"19/12/23": 47180.52, "19/12/34": 45910.67, "19/12/35": 46125.84, "19/12/36": 46403.90, "19/12/27": 45907.13,
		"19/12/30": 46204.47, "19/12/31": 45594.42, "20/01/02": 45620.39, "20/01/03": 45204.59,
		"20/01/06": 45166.33, "20/01/07": 44610.73, "20/01/08": 41244.05, "20/01/09": 42062.69, "20/01/10": 43518.37,
		"20/01/13": 42247.73, "20/01/14": 43262.82, "20/01/15": 43536.69, "20/01/16": 64668.56,
	}
)

func handler(w http.ResponseWriter, _ *http.Request) {
	var (
		nameItems []string
		foodItems []float64
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
		if k >= "19/11/07" && k < "19/11/11" {
			foodItems = append(foodItems, data[k]-20000)
		} else if k >= "19/11/11" && k < "19/11/13" {
			foodItems = append(foodItems, data[k]-25000)
		} else if k >= "19/11/13" && k < "20/01/16" {
			foodItems = append(foodItems, data[k]-50000)
		} else if k >= "20/01/16" {
			foodItems = append(foodItems, data[k]-70000)
		}
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
