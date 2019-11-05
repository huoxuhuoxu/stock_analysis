package main

import (
	"fmt"
	"sort"
	"stock_analysis/base"
	"strings"
	"time"
	"unicode/utf8"
)

type RealTime struct {
	*base.MainControl
	list                      base.GeneralDatas // 数据列表
	loopIntervalTime          time.Duration     // 更新频率
	filterLen                 int               // 过滤数量
	estimatedCloseCoefficient float64           // 预计收盘价系数(上涨1%, 目前固定值)
	estimatedLowCoefficient   float64           // 预计最低价(下跌1.5%)
	isDebug                   bool
}

func NewRealTime(mc *base.MainControl, isDebug bool) (*RealTime, error) {
	self := &RealTime{
		MainControl:               mc,
		loopIntervalTime:          time.Second * 10,
		filterLen:                 40,
		estimatedCloseCoefficient: 1.01,
		estimatedLowCoefficient:   0.995,
		isDebug:                   isDebug,
	}

	return self, nil
}

func (self *RealTime) Running() {
	self.loop()
}

// 轮询, 处理实时数据
func (self *RealTime) loop() {
	go func() {
		for {
			select {
			case <-self.Ctx.Done():
				self.Output.Log("Realtime, Ctx.Done")
				return
			default:
				// pass
			}

			// 计算得分
			gds := base.GetMoneyData()
			var list base.GeneralDatas
			for _, gd := range gds {
				gd.Score = gd.Changehands + gd.MainP

				// 限制最低交易额(5000w)
				if gd.Turnover < 0.5 {
					continue
				}
				// 限制最低换手率(0.4%)
				if gd.Changehands < 0.4 {
					continue
				}
				// 限制最大流通盘(100亿)
				if gd.Circulation > 100 {
					continue
				}

				list = append(list, gd)
			}
			sort.Sort(list)
			self.list = list[0:self.filterLen]
			self.FormatOutput(self.list)

			if isDebug {
				self.Output.Log("Debug-Model, ready exit")
				self.CtxCancel()
				return
			}
			time.Sleep(self.loopIntervalTime)
		}
	}()
}

// 格式化输出数据
func (self *RealTime) FormatOutput(list base.GeneralDatas) {
	fmt.Println(strings.Repeat("--", 80))
	fmt.Printf(strings.Repeat("|%-6s ", 16)+"|\r\n",
		"评分",
		"代码",
		"名称",
		"涨幅",
		"成交额",
		"主力占比",
		"换手率",
		"5min",
		"开盘价",
		"预收盘价",
		"现价",
		"开盘波动",
		"预低点",
		"低点",
		"振幅",
		"操作建议",
	)

	var (
		openWin, openFail, todayWin, todayFail float64
	)
	for _, v := range list {
		if utf8.RuneCountInString(v.StockName) < 4 {
			diff := 4 - utf8.RuneCountInString(v.StockName)
			v.StockName += strings.Repeat("补", diff)
		}
		var (
			score          = fmt.Sprintf("%.2f", v.Score)
			stockCode      = fmt.Sprintf("  %s", v.StockCode)
			stockName      = fmt.Sprintf("     %s", v.StockName)
			upAndDownRange = fmt.Sprintf("  %.2f%%", v.UpAndDownRange)
			turnover       = fmt.Sprintf("   %.2f亿", v.Turnover)
			mainP          = fmt.Sprintf("    %.2f%%", v.MainP)
			min5           = fmt.Sprintf("   %.2f%%", v.Min5)
			changehands    = fmt.Sprintf("   %.2f%%", v.Changehands)
			open           = fmt.Sprintf("   %.2f", v.Open)
			estimatedClose = fmt.Sprintf("      %.2f", v.Open*self.estimatedCloseCoefficient)
			currentPrice   = fmt.Sprintf("      %.2f", v.CurrentPrice)
			estimatedLow   = fmt.Sprintf("      %.2f", v.Open*self.estimatedLowCoefficient)
			low            = fmt.Sprintf("   %.2f", v.Low)
			amplitude      = fmt.Sprintf("     %.2f%%", v.Amplitude)
		)

		// 如果开盘涨幅 > 9%, 则认为收盘价 = 涨停
		if (v.Open-v.Close)/v.Close > 0.09 {
			estimatedClose = fmt.Sprintf("      %.2f", v.Close*1.1)
		}

		// 操作建议
		var action string = "        /"
		diffOc := (v.CurrentPrice - v.Open) / v.Open
		if diffOc < -0.02 {
			action = "买进"
		} else if diffOc < -0.01 {
			action = "关注"
		}
		diffOcStr := fmt.Sprintf("     %.2f%%", diffOc*100)

		// 开盘价胜率统计
		if v.CurrentPrice > v.Open {
			openWin += 1.0
		} else {
			openFail += 1.0
		}
		// 今日胜率统计
		if v.CurrentPrice > v.Close {
			todayWin += 1.0
		} else {
			todayFail += 1.0
		}

		// 不存在博弈空间
		if v.UpAndDownRange > 9.9 {
			continue
		}

		// Output
		fmt.Printf(strings.Repeat("%+6s ", 16)+"\r\n",
			score,
			stockCode,
			stockName,
			upAndDownRange,
			turnover,
			mainP,
			changehands,
			min5,
			open,
			estimatedClose,
			currentPrice,
			diffOcStr,
			estimatedLow,
			low,
			amplitude,
			action,
		)
	}

	fmt.Println(strings.Repeat("--", 80))
	fmt.Printf("开盘价胜率 %.2f%% \r\n", openWin/float64(self.filterLen)*100)
	fmt.Printf("开盘价败率 %.2f%% \r\n", openFail/float64(self.filterLen)*100)
	fmt.Printf("今日胜率 %.2f%% \r\n", todayWin/float64(self.filterLen)*100)
	fmt.Printf("今日败率 %.2f%% \r\n", todayFail/float64(self.filterLen)*100)
	fmt.Println(strings.Repeat("--", 80))
}
