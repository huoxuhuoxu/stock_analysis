package main

type GeneralData struct {
	stockType      string  // 上市类型
	stockCode      string  // 代码
	stockName      string  // 名称
	currentPrice   float64 // 当前价格
	upAndDownRange float64 // 涨跌幅度
	turnover       float64 // 成交额
	high           float64 // 盘中最低
	low            float64 // 盘中最高
	open           float64 // 开盘价
	close          float64 // 前个交易日收盘价
	min5           float64 // 5分钟涨速
	changehands    float64 // 换手率
	circulation    float64 // 流通市值

	// 所有的流入与流出其实只是一个资金意向问题, 一笔交易有钱出去就意味着有钱进来, 盈亏只是这笔交易关联的上下文的价格的差价
	mainV   float64 // 主力净流入(由超大单与大单组成)
	cBigV   float64 // 超大单净流入
	bigV    float64 // 大单净流入
	middleV float64 // 中单净流入
	smallV  float64 // 小单净流入

	mainP float64 // 主力流入占比

	score float64 // 评分
}

type GeneralDatas []*GeneralData

func (self GeneralDatas) Len() int {
	return len(self)
}

func (self GeneralDatas) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self GeneralDatas) Less(i, j int) bool {
	return self[i].score > self[j].score
}
