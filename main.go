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
	var bufferW bytes.Buffer

	line1 := income()
	line1.Render(&bufferW)

	line3 := income2()
	line3.Render(&bufferW)

	// 替换标题
	htmlV := bufferW.String()
	htmlV = strings.Replace(htmlV, "Awesome go-echarts", "逐日结算", 1)

	// 输出html
	w.Write([]byte(htmlV))
}

// 收益曲线 与 剔除手续费 曲线
func income() *charts.Line {
	var (
		nameItems  []string
		foodItems  []int // 总结算
		foodItems2 []int // 去除手续费
		keySlice   sort.StringSlice
		fee        float64 // 手续费累积
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
		} else if k >= "20/01/17" && k < "20/01/20" {
			initV = 80000
		} else {
			initV = 100000
		}

		foodItems = append(foodItems, int(data[k]-initV))

		fee += handlingFee[k]
		foodItems2 = append(foodItems2, int(data[k]-initV+fee))
	}

	// 画图表
	line := charts.NewLine()
	line.SetGlobalOptions(charts.TitleOpts{Title: "收益曲线"})
	line.AddXAxis(nameItems).AddYAxis("盈亏", foodItems,
		// charts.LabelTextOpts{Show: true},
		charts.AreaStyleOpts{Opacity: 0.2},
		charts.LineOpts{Smooth: true},
	)
	line.AddXAxis(nameItems).AddYAxis("盈亏(去除手续费)", foodItems2,
		charts.AreaStyleOpts{Opacity: 0.2},
		charts.LineOpts{Smooth: true},
	)
	return line
}

// 手续费 与 收益曲线 折线图
func income2() *charts.Bar {
	var (
		nameItems  []string
		foodItems  []int // 收益
		foodItems2 []int // 手续费
		previous   int   // 上一个偏差值
		keySlice   sort.StringSlice
	)

	// 日期排序
	for k := range data {
		keySlice = append(keySlice, k)
	}
	sort.Sort(keySlice)

	// 计算盈亏
	for i, k := range keySlice {
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
		} else if k >= "20/01/17" && k < "20/01/20" {
			initV = 80000
		} else {
			initV = 100000
		}

		if i != 0 {
			foodItems = append(foodItems, int(data[k]-initV)-previous)
		} else {
			foodItems = append(foodItems, int(data[k]-initV))
		}

		previous = int(data[k] - initV)
		foodItems2 = append(foodItems2, int(-handlingFee[k]))
	}

	// 画图表
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: "收益与手续费"},
		charts.ToolboxOpts{Show: true},
		charts.DataZoomOpts{XAxisIndex: []int{0}, Start: 50, End: 100},
	)
	bar.AddXAxis(nameItems).
		AddYAxis("收益", foodItems, charts.BarOpts{Stack: "stack"}).
		AddYAxis("手续费", foodItems2, charts.BarOpts{Stack: "stack"})

	return bar
}

func main() {
	log.Println("Start !")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
