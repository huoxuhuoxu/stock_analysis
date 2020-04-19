package main

var (
	headers = map[string]string{
		"Accept":                    " text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Accept-Encoding":           "gzip, deflate",
		"Accept-Language":           "zh-CN,zh;q=0.9,en;q=0.8",
		"Cache-Control":             "no-cache",
		"Cookie":                    "qgqp_b_id=175ececc76779338ed237e0147093167; st_si=48998687825010; st_asi=delete; emshistory=%5B%22bai%20y%22%5D; HAList=a-sz-000061-%u519C%20%u4EA7%20%u54C1; em_hq_fls=js; st_pvi=61646748513010; st_sp=2020-03-29%2023%3A01%3A34; st_inirUrl=http%3A%2F%2Ffutures.eastmoney.com%2Fqihuo%2FM.html; st_sn=215; st_psi=20200401231508322-113200301325-3052892464",
		"Host":                      "futsse.eastmoney.com",
		"Pragma":                    "no-cache",
		"Proxy-Connection":          "keep-alive",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
	}
)

const (
	ORIGIN_URL = "http://futsse.eastmoney.com/static/"
)

// 合约
type Contract struct {
	code  string  // 商品编号
	open  float64 // 开盘价
	close float64 // 昨日收盘价
	high  float64 // 盘中最高价
	low   float64 // 盘中最低价
	price float64 // 现价
}

// 品种
type Variety struct {
	Name           string  // 商品名称
	OriginDataUrl  string  // 数据源头
	Price          float64 // 当前价格
	Value          float64 // 涨跌值
	SpotPrice      float64 // 现货价格
	Aims           string  // 目标范围
	Describe       string  // 说明
	PricePrecision uint8   // 价格精度
	Level          int     // 优先级, 1: 优先事项, 2: 逐步进行 3: 可试可不试, 4: 观察, 等待, 有极好的机会可试
	IsShow         bool    // 是否获取及输出
	// 还需要关联标的物价格, 比如 美原油/铁矿石/美豆/美黄金 ..., 之后加一个字段, 带上关联合约的信息, 辅助输出
}

// 请记住, 期货是给远期现货定价的, 而不是用现货现在的价格给期货定价, 谨记
var varietys = map[string]*Variety{
	"sc2009": &Variety{
		Name:          "原油",
		OriginDataUrl: "142_sc2009_qt?callbackName=aa&cb=aa&_=1587309656220",
		SpotPrice:     0,
		Aims:          "250~270",
		Describe: `
			假设五月原油价格是其实际供需进而产生的, 
			那么后续月份的现在价格需要加上时间周期的仓储费用, 由买家买单,
			但一个月后基本面不改善, 价格也会到五月的价格, 仓储费用就由卖家买单,
			我的观点, 下半年要么不会复苏, 要么稳健复苏, 并且远月保持高位的价格是仓储费的叠加,
			从时间周期的角度上看, 实际原油的价格只会降低, 最终买方需要支付的其实是卖方仓储的费用,
			暂看震荡下行
		`,
		Level:          1,
		PricePrecision: 1,
		IsShow:         true,
	},
	"au2012": &Variety{
		Name:          "黄金",
		OriginDataUrl: "113_au2012_qt?callbackName=aa&cb=aa&_=1587309117276",
		SpotPrice:     0,
		Aims:          "360~400",
		Describe: `
			暂时没有看法, 但是看技术面进入平台, 遇阻, 下跌/回踩 是大概率事件
		`,
		Level:          3,
		PricePrecision: 2,
		IsShow:         true,
	},
	"rb2101": &Variety{
		Name:          "螺纹",
		OriginDataUrl: "113_rb2101_qt?callbackName=aa&cb=aa&_=1587308954178",
		SpotPrice:     3427,
		Aims:          "3400",
		Describe: `
			供需强劲, 关联原料铁矿石也很强劲, 进入平台区, 一旦突破, 3400指日可待
		`,
		Level:          1,
		PricePrecision: 0,
		IsShow:         true,
	},
	"m2101": &Variety{
		Name:          "豆粕",
		OriginDataUrl: "114_m2101_qt?callbackName=aa&cb=aa&_=1585752611719",
		SpotPrice:     3167.5,
		Aims:          "2720",
		Describe: `
			受巴西大豆, 美大豆到港量影响, 美大豆价格不断近期新低
			进口国外猪肉打压国内猪肉价格影响
			短期豆粕低位横盘，甚至会在破新低, 目标01合约, 2720
		`,
		Level:          2,
		PricePrecision: 0,
		IsShow:         true,
	},
}
