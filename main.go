package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/go-echarts/go-echarts/charts"
)

/*
	目前账号期货公司手续费
		螺纹: 万2
		豆粕, 棕榈, 燃油, 白糖, 橡胶: 万1
*/

const (
	MONEY        = 120000.0
	PROFIT_MONEY = 4000.0
)

func handler(w http.ResponseWriter, _ *http.Request) {
	var bufferW bytes.Buffer

	buf, _ := ioutil.ReadFile("./data.json")
	var arrData []interface{}
	json.Unmarshal(buf, &arrData)
	data := arrData[0].(map[string]interface{})
	handlingFee := arrData[1].(map[string]interface{})

	line, fee := income(data, handlingFee)
	bar, moneyWin, str := income2(data, handlingFee)

	str2 := `
	<div style="text-align:center;">
		<p>收益率, 剥离手续费后: %.2f%%</p>
	</div>
	`
	str2 = fmt.Sprintf(str2, fee+moneyWin)

	bufferW.Write([]byte(str))
	bufferW.Write([]byte(str2))
	line.Render(&bufferW)
	bar.Render(&bufferW)

	// 替换标题
	htmlV := bufferW.String()
	htmlV = strings.Replace(htmlV, "Awesome go-echarts", "逐日结算", 1)

	// 输出html
	w.Write([]byte(htmlV))
}

// 收益曲线 与 剔除手续费 曲线
func income(data, handlingFee map[string]interface{}) (*charts.Line, float32) {
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
		var initV = getMoney(k)

		foodItems = append(foodItems, int(data[k].(float64)-initV))

		fee += handlingFee[k].(float64)
		foodItems2 = append(foodItems2, int(data[k].(float64)-initV+fee))
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

	// 收益率, 剥离手续费后
	diffRate := float32(foodItems2[len(foodItems2)-1]-foodItems[len(foodItems)-1]) / MONEY * 100

	return line, diffRate
}

// 手续费 与 收益曲线 折线图, 综合胜率
func income2(data, handlingFee map[string]interface{}) (*charts.Bar, float32, string) {
	var (
		nameItems  []string
		foodItems  []int // 收益
		foodItems2 []int // 手续费
		previous   int   // 上一个偏差值
		winCount   int   // 胜率总数
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
		var initV = getMoney(k)

		// 补偿获利取现的资金
		vv := data[k].(float64)
		if k >= "20/03/27" {
			vv += PROFIT_MONEY
		}

		var v int
		if i != 0 {
			v = int(vv-initV) - previous
		} else {
			v = int(vv - initV)
		}

		if v > 0 {
			winCount++
		}

		previous = int(vv - initV)
		foodItems = append(foodItems, v)
		foodItems2 = append(foodItems2, int(-handlingFee[k].(float64)))
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

	// 综合
	str := `
	<div style="text-align:center;">
		<p>胜率: %.2f%%</p>
		<p>收益率: %.2f%%</p>
	</div>
	`
	rateOfReturn := (float32(data[keySlice[len(keySlice)-1]].(float64)) + PROFIT_MONEY - MONEY) / MONEY * 100
	str = fmt.Sprintf(str, float32(winCount)/float32(len(foodItems))*100, rateOfReturn)

	return bar, rateOfReturn, str
}

// 时间与本金关系
func getMoney(k string) float64 {
	var initV float64
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
	} else if k >= "20/01/20" && k < "20/02/03" {
		initV = 100000
	} else if k >= "20/02/03" && k < "20/02/12" {
		initV = 110000
	} else {
		initV = MONEY
	}
	return initV
}

func main() {
	log.Println("Start !")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
