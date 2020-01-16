package base

// raw-data -------------------------------------------------
type GeneralData struct {
	StockType      string  // 上市类型
	StockCode      string  // 代码
	StockName      string  // 名称
	CurrentPrice   float64 // 当前价格
	UpAndDownRange float64 // 涨跌幅度
	Turnover       float64 // 成交额
	High           float64 // 盘中最低
	Low            float64 // 盘中最高
	Open           float64 // 开盘价
	Close          float64 // 前个交易日收盘价
	Min5           float64 // 5分钟涨速
	Changehands    float64 // 换手率
	Circulation    float64 // 流通市值
	Amplitude      float64 // 振幅

	// 所有的流入与流出其实只是一个资金意向问题, 一笔交易有钱出去就意味着有钱进来, 盈亏只是这笔交易关联的上下文的价格的差价
	MainV   float64 // 主力净流入(由超大单与大单组成)
	CBigV   float64 // 超大单净流入
	BigV    float64 // 大单净流入
	MiddleV float64 // 中单净流入
	SmallV  float64 // 小单净流入

	MainP float64 // 主力流入占比

	Score float64 // 评分
}

type GeneralDatas []*GeneralData

func (self GeneralDatas) Len() int {
	return len(self)
}

func (self GeneralDatas) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self GeneralDatas) Less(i, j int) bool {
	return self[i].Score > self[j].Score
}
