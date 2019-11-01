package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	minTurnover              float64 = 2    // 最少成交额 2 亿
	minChangehands           float64 = 2    // 最低换手 2%
	minTurnoverToCirculation float64 = 0.02 // 成交占比流通最少 万2
	minMainP                 float64 = 5    // 主力流入占比最低 5%
	mainMainV                float64 = 0.1  // 主力流入最少 0.1 亿

	// 其实换手率就可以体现 成交额占流通市值的比例
)

func init() {
	flag.Float64Var(&minTurnover, "mt", 2, "最少成交额")
	flag.Float64Var(&minChangehands, "mc", 2, "最低换手")
	flag.Float64Var(&minTurnoverToCirculation, "mttc", 0.02, "成交占比流通")
	flag.Float64Var(&minMainP, "mp", 5, "主力流入占比")
	flag.Float64Var(&mainMainV, "mv", 0.1, "主力流入")
	flag.Parse()
}

func main() {
	gds := GetMoneyData()
	log.Println(len(gds))

	showGds := make([]*GeneralData, 0)
	for _, gd := range gds {
		gd.turnover = gd.turnover / 100000000
		gd.mainV = gd.mainV / 100000000

		if gd.mainV < mainMainV {
			continue
		}
		if gd.mainP < minMainP {
			continue
		}
		if gd.turnoverToCirculation < minTurnoverToCirculation {
			continue
		}
		if gd.changehands < minChangehands {
			continue
		}
		if gd.turnover < minTurnover {
			continue
		}
		showGds = append(showGds, gd)
	}

	fmt.Printf("%-6s, %-4s, %-4s, %-4s, %-4s, %-4s, %-4s, %-4s \r\n",
		"代码", "名称", "涨跌幅", "成/流", "主力流入", "占比", "5min", "换手")

	var (
		total, fCount, zCount float64
	)
	for _, gd := range showGds {
		fmt.Printf("%6s, %4s, %6.2f%%, %6.4f%%, %6.2f亿, %6.2f%%, %6.2f%%, %2.2f%% \r\n",
			gd.stockCode,
			gd.stockName,
			gd.upAndDownRange,
			gd.turnoverToCirculation,
			gd.mainV,
			gd.mainP,
			gd.min5,
			gd.changehands,
		)

		total++
		if gd.upAndDownRange > 0 {
			zCount++
		} else {
			fCount++
		}
	}

	fmt.Printf("\r\n涨幅胜率: %.2f%%\r\n", zCount/total*100)
	fmt.Printf("跌幅胜率: %.2f%%", fCount/total*100)
}
