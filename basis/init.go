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

	/*
		@desc
			basis < limit: 开仓
			basis > profit: 反向/平仓
			other: 等待
	*/
	actions = []string{"等待", "开仓", "反向/平仓"}
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

type Variety struct {
	Name            string  // 商品名称
	OriginDataUrl   string  // 数据源头
	Price           float64 // 当前价格
	Value           float64 // 涨跌值
	PricePrecision  uint8   // 价格精度
	Amount          int     // 套利持仓数
	Dash            float64 // 划点, 对价平仓需要付出的价钱
	DashCoefficient float64 // 每个单位点的价钱
}

type Group struct {
	Name        string    // 组合名称
	Combination [2]string // 包含合约, 持仓方向(0: 多, 1: 空)
	// 基差临界点数, 预计中可能会存在超跌值, 以此值为起点开始套利, 最关键的一个值
	Limit             float64
	Profit            float64 // 利润点数
	MarginConsumption string  // 组合需要消耗保证金
	Level             int     // 优先级, 99: 待验证待的逻辑, 100: 等待
	Describe          string  // 组合逻辑说明
}

// 反套组
var groups = []Group{
	// 套差与趋势跟随
	Group{
		// 需要追加基差回归的假设, 在主逻辑影响下
		// 为什么对赌是 多沥青 空燃油 的假设, 为什么不能倒过来, 需要解释
		// 增加相对锚定价差 2118~2090 : 1622
		Name:              "多沥青/空燃油",
		Combination:       [2]string{"bu2012", "fu2101"},
		Limit:             5,
		MarginConsumption: "5",
		Level:             1,
		Profit:            30,
		Describe: `
			在原油为主导因素情况下, 
			利用多空对冲对与原油具有高度相关性的衍生品进行对赌,
			赌 沥青副逻辑-基建 强于 燃油副逻辑-航运,
			不论涨或者跌, 在盘中体现, 沥青强于燃油, 
			日内会产生 20～40 点的套利机会
		`,
	},
	Group{
		Name:              "多黄金/空白银",
		Combination:       [2]string{"au2012", "ag2012"},
		Limit:             30,
		MarginConsumption: "50",
		Level:             999,
		Profit:            150,
		Describe: `
			反套逻辑弱, 不能做
			但是, 白银跟着黄金走, 可以做参考
			比如 04/27, 黄金盘中下跌, 白银暂时坚挺, 后跟跌
			可以用于做提前预判
		`,
	},
	Group{
		Name:              "多铁矿/空螺纹",
		Combination:       [2]string{"i2009", "rb2101"},
		Limit:             10,
		MarginConsumption: "20",
		Level:             999,
		Profit:            80,
		Describe: `
			反套逻辑弱, 不能做, 如果非要做, 需要很强的选择,
			但是, 这两个品种互相影响, 属于上下游关系, 可以做参考
			铁矿强, 螺纹也不会太弱, 反之亦然
		`,
	},
	Group{
		Name:              "多菜粕/空豆粕",
		Combination:       [2]string{"rm009", "m2101"},
		Limit:             5,
		MarginConsumption: "4",
		Level:             999,
		Profit:            20,
		Describe: `
			反套, 有套利空间, 大约20个点上下
			菜粕波动大于豆粕, 在排除极端行情的情况下, 两倍的关系
			盘中建仓可以排除开盘的极端行情, 等待基差出现
			负基差建仓最好, 菜粕比较妖, 可以拉回来
			对赌逻辑:
				1. 菜粕持仓小 更能够将基本面的变化体现在价格上
				2. 菜粕持仓小 更容易出现超跌或者超买的情况
				3. 在趋势偏弱, 利用豆粕对冲其大波动下跌时的损失, 赌两者的波动点数, 最终会回归相近的位置
				4. 在趋势偏强, 菜粕由于持仓小的原因, 更容易导致价格上涨, 从而产生与豆粕的基差点数, 从中套利
		`,
	},
}

// 合约集合
var varietys = map[string]*Variety{
	"bu2012": &Variety{
		Name:            "沥青",
		OriginDataUrl:   "113_bu2012_qt?callbackName=aa&cb=aa&_=1587532980754",
		PricePrecision:  0,
		Amount:          1,
		Dash:            2,
		DashCoefficient: 10,
	},
	"fu2101": &Variety{
		Name:            "燃油",
		OriginDataUrl:   "113_fu2101_qt?callbackName=aa&cb=aa&_=1587533000837",
		PricePrecision:  0,
		Amount:          1,
		Dash:            1,
		DashCoefficient: 10,
	},
	"au2012": &Variety{
		Name:            "沪金",
		OriginDataUrl:   "113_au2012_qt?callbackName=aa&cb=aa&_=1587566439644",
		PricePrecision:  2,
		Amount:          1,
		Dash:            0.02,
		DashCoefficient: 1000,
	},
	"ag2012": &Variety{
		Name:            "白银",
		OriginDataUrl:   "113_ag2012_qt?callbackName=aa&cb=aa&_=1587533190507",
		PricePrecision:  0,
		Amount:          3,
		Dash:            1,
		DashCoefficient: 15,
	},
	"rm009": &Variety{
		Name:            "菜粕",
		OriginDataUrl:   "115_RM009_qt?callbackName=aa&cb=aa&_=1587831952152",
		PricePrecision:  0,
		Amount:          1,
		Dash:            1,
		DashCoefficient: 10,
	},
	"m2101": &Variety{
		Name:            "豆粕",
		OriginDataUrl:   "114_m2101_qt?callbackName=aa&cb=aa&_=1585752611719",
		PricePrecision:  0,
		Amount:          1,
		Dash:            1,
		DashCoefficient: 10,
	},
	"rb2101": &Variety{
		Name:            "螺纹",
		OriginDataUrl:   "113_rb2101_qt?callbackName=aa&cb=aa&_=1587308954178",
		PricePrecision:  0,
		Amount:          3,
		Dash:            1,
		DashCoefficient: 10,
	},
	"i2009": &Variety{
		Name:            "铁矿石",
		OriginDataUrl:   "114_i2009_qt?callbackName=aa&cb=aa&_=1587396948222",
		PricePrecision:  1,
		Amount:          2,
		Dash:            0.5,
		DashCoefficient: 100,
	},
}
