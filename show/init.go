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
}

/*
	请记住:
		1. 期货是给远期现货定价的, 而不是用现货现在的价格给期货定价, 谨记
		2. 做短线就要有做短线的操守, 不持仓过夜, 日内 平仓/锁仓
		3. 做长线就要有承担压力的心理准备
		4. 被套了也不要恐慌, 短期的大起大落, 基于均值回归, 最终都会回到均线附近, 市场会给出解套的机会
		5. 左侧交易 承压与收益成正比, 右侧交易 顺势而行
		6. 长远来看任何一个商品最终都会回归一个被长时间考验的相对内在价值的合理价格
		7. 相信自己的判断, 坚信自己的决定, 坚定不移的执行定下的策略, 贯彻到底, 始终如一
*/
var varietys = map[string]*Variety{
	"cl00y": &Variety{
		Name:           "美油06",
		OriginDataUrl:  "102_CL00Y_qt?callbackName=aa&cb=aa&_=1587646825954",
		SpotPrice:      0,
		Aims:           "",
		Describe:       ``,
		Level:          1,
		PricePrecision: 2,
		IsShow:         true,
	},
	"cl20u": &Variety{
		Name:           "美油09",
		OriginDataUrl:  "102_CL20U_qt?callbackName=aa&cb=aa&_=1587646921066",
		SpotPrice:      0,
		Aims:           "",
		Describe:       ``,
		Level:          1,
		PricePrecision: 2,
		IsShow:         true,
	},
	"sc2006": &Variety{
		Name:           "原油06",
		OriginDataUrl:  "142_sc2006_qt?callbackName=aa&cb=aa&_=1587309656220",
		SpotPrice:      0,
		Aims:           "",
		Describe:       ``,
		Level:          2,
		PricePrecision: 1,
		IsShow:         true,
	},
	"sc2009": &Variety{
		Name:          "原油09",
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
		Level:          2,
		PricePrecision: 1,
		IsShow:         true,
	},
	"gc20z": &Variety{
		Name:           "黄金",
		OriginDataUrl:  "101_GC20Z_qt?callbackName=aa&cb=aa&_=1587647068923",
		SpotPrice:      0,
		Aims:           "",
		Describe:       ``,
		Level:          3,
		PricePrecision: 1,
		IsShow:         true,
	},
	"au2012": &Variety{
		Name:          "沪金",
		OriginDataUrl: "113_au2012_qt?callbackName=aa&cb=aa&_=1587566439644",
		SpotPrice:     0,
		Aims:          "360~390",
		Describe: `
			围绕 360 上下 30 个点震荡, 跟着外盘跑
			技术面上看在 380/390 受压制, 但
			风险较大, 伦敦黄金可能存在逼仓的可能, 导致外盘黄金突破新高
			进而引发内盘黄金突破平台区, 观察为主, 不动
		`,
		Level:          4,
		PricePrecision: 2,
		IsShow:         true,
	},
}

/*
	请严格执行
*/
var dailyOperation = `
	外盘美油 近月反弹多, 沅月反弹少
	国内原油 近月反弹少, 沅月反弹多, 主要原因在于:
		1. 美油近月跌多, 远月跌少, 在基本面可能改善的情况下, 近月反弹空间更大
		2. 内盘远近合约跌势相近, 并且在对未来有较好预期的情况下, 远月反弹预期高于近月
		3. 基本美油这波价格战中, 远月价格中枢依旧维持在 30～35 之间


	操作该要:
		原油,明日可交易情况下,
		不出意外明日开盘以涨为主, 如果平开或者09合约280以下开, 可短多

		黄金,
		盘中看趋势, 可短空, 目前处于压力位, 不管内外盘, 可以顺势空一波
`
