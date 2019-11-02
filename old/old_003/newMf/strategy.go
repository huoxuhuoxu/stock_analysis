package main

import (
	"fmt"
	"sort"
)

/*
	@DESC
		成交额/流通市值 = 换手率

	假设
		任何一个市场资金都是在流动的, 那么跟踪资金最好的方式就是成交额
		换手率表示单股, 热度/交易热情/追逐

	评分 = 换手 + 主力流入占比
*/

func strategy001(gds map[string]*GeneralData) {

	// 开始打分
	var list GeneralDatas
	for _, gd := range gds {
		gd.score = gd.changehands + gd.mainP
		list = append(list, gd)
	}
	sort.Sort(list)

	// 输出
	// for _, gd := range list {
	// fmt.Printf("%2.2f, %2.2f%% \r\n", gd.score, gd.upAndDownRange)
	// fmt.Printf("冷静分析: %s, %2.2f%%, %2.2f, %2.2f \r\n", gd.stockName, gd.changehands/gd.score, gd.changehands, gd.mainP)
	// }

	// 各分组情况
	// statistics001(list)

	// 第一分组详情
	// statistics002(list[0:100])

	// 开盘作为初始价格, 计算涨跌幅
	// statistics003(list[0:100])

	// 其他分组情况
	// statistics003(list[100:200])

	// 盘中最高与开盘价, 即 盘中操作空间
	// statistics004(list[0:100])

	// 盘中最高价与最低价, 即 盘中操作空间
	// statistics005(list[0:100])

	// 跌破开盘价, 即 盘中低吸机会
	// statistics006(list[0:100])

	// 最低价与收盘价, 即 盘中反弹情况
	// statistics007(list[0:100])

	// 前一个交易日收盘价与当日开盘价, 价差
	// statistics008(list[0:100])

	// 前一个交易日收盘价与当日收盘价, 价差
	statistics009(list[0:100])
}

/*
	@DESC
		经过统计验证, 第一组属于有效分组, 涨幅及大涨幅主要集中在第一分组中

		不同涨幅, 在各个分组内的占有率情况
*/
func statistics001(list GeneralDatas) {
	var (
		step               = 100 // 分割统计
		count              = 0   // 记录符合
		totalCount         = 0   // 一共符合
		condUp     float64 = 0   // 分割条件 - 涨跌幅
		records    []int         // 记录分组结果
	)

	for i, gd := range list {
		if i%step == 0 && i != 0 {
			records = append(records, count)
			totalCount += count
			count = 0
		}
		if gd.upAndDownRange < condUp {
			count++
		}
	}

	records = append(records, count)
	totalCount += count

	for _, record := range records {
		fmt.Printf("%d, %2.2f%% \r\n", record, float64(record)/float64(totalCount)*100)
	}
}

/*
	@DESC
		第一分组内的涨跌幅情况

*/
func statistics002(list GeneralDatas) {
	var (
		upCount   int
		downCount int

		condUp float64 = 0
	)

	for _, gd := range list {
		if gd.upAndDownRange > condUp {
			upCount++
		} else {
			downCount++
		}
	}

	fmt.Printf("胜率: %2.2f \r\n", float64(upCount)/float64(len(list))*100)
	fmt.Printf("败率: %2.2f \r\n", float64(downCount)/float64(len(list))*100)
}

/*
	@DESC
		从开盘到收盘, 涨跌情况
*/
func statistics003(list GeneralDatas) {
	var (
		upCount   int
		downCount int

		condUp float64 = 1
	)

	for _, gd := range list {
		if gd.currentPrice-gd.open > condUp {
			upCount++
		} else {
			downCount++
		}
	}

	fmt.Printf("胜率: %2.2f \r\n", float64(upCount)/float64(len(list))*100)
	fmt.Printf("败率: %2.2f \r\n", float64(downCount)/float64(len(list))*100)
}

/*
	@DESC
		从开盘到盘中最高, 涨跌情况
*/
func statistics004(list GeneralDatas) {
	var (
		upCount   int
		downCount int

		condUp float64 = 5
	)

	for _, gd := range list {
		if gd.high-gd.open > condUp {
			upCount++
		} else {
			downCount++
		}
	}

	fmt.Printf("胜率: %2.2f \r\n", float64(upCount)/float64(len(list))*100)
	fmt.Printf("败率: %2.2f \r\n", float64(downCount)/float64(len(list))*100)
}

/*
	@DESC
		最高价与最低价差值, 盘中振幅
*/
func statistics005(list GeneralDatas) {
	var (
		upCount   int
		downCount int

		condUp float64 = 10
	)

	for _, gd := range list {
		if (gd.high-gd.low)/gd.low*100 > condUp {
			upCount++
		} else {
			downCount++
		}
	}

	fmt.Printf("胜率: %2.2f \r\n", float64(upCount)/float64(len(list))*100)
	fmt.Printf("败率: %2.2f \r\n", float64(downCount)/float64(len(list))*100)
}

/*
	@DESC
		开盘价与最低价差值, 跌破开盘价, 低吸机会
*/
func statistics006(list GeneralDatas) {
	var (
		upCount   int
		downCount int

		condUp float64 = 4
	)

	for _, gd := range list {
		if (gd.open-gd.low)/gd.open*100 > condUp {
			upCount++
		} else {
			downCount++
		}
	}

	fmt.Printf("胜率: %2.2f \r\n", float64(upCount)/float64(len(list))*100)
	fmt.Printf("败率: %2.2f \r\n", float64(downCount)/float64(len(list))*100)
}

/*
	@DESC
		最低价与收盘价, 反弹回来的概率
*/
func statistics007(list GeneralDatas) {
	var (
		upCount   int
		downCount int

		condUp float64 = 10
	)

	for _, gd := range list {
		if (gd.currentPrice-gd.low)/gd.low*100 > condUp {
			upCount++
		} else {
			downCount++
		}
	}

	fmt.Printf("胜率: %2.2f \r\n", float64(upCount)/float64(len(list))*100)
	fmt.Printf("败率: %2.2f \r\n", float64(downCount)/float64(len(list))*100)
}

/*
	@DESC
		前交易日收盘与今日开盘, 价差
*/
func statistics008(list GeneralDatas) {
	var (
		upCount   int
		downCount int

		condUp float64 = 2
	)

	for _, gd := range list {
		if (gd.open-gd.close)/gd.close*100 > condUp {
			upCount++
		} else {
			downCount++
		}
	}

	fmt.Printf("胜率: %2.2f \r\n", float64(upCount)/float64(len(list))*100)
	fmt.Printf("败率: %2.2f \r\n", float64(downCount)/float64(len(list))*100)
}

/*
	@DESC
		前交易日收盘与今日收盘, 价差
*/
func statistics009(list GeneralDatas) {
	var (
		upCount   int
		downCount int

		condUp float64 = 7
	)

	for _, gd := range list {
		if (gd.currentPrice-gd.close)/gd.close*100 > condUp {
			upCount++
		} else {
			downCount++
		}
	}

	fmt.Printf("胜率: %2.2f \r\n", float64(upCount)/float64(len(list))*100)
	fmt.Printf("败率: %2.2f \r\n", float64(downCount)/float64(len(list))*100)
}
