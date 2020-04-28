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

/*
	摘除五个反套逻辑:
		黄金白银: 回报比低, 承担风险高
		螺纹铁矿, 豆粕菜粕: 其中一个品种走势强于另一个, 并且多空都是
		豆油棕榈, 橡胶棉花: 一个纯品种上关联性小, 分化比较厉害, 一个互相的带着跑无分化

	反套就留: 沥青/燃油
	单边: 原油及其同逻辑衍生品(沥青, 燃油)
	注意: 甲醇是化工, 走势不完全随原油, 不搞, 风险也不小的
*/

// 反套组
var groups = []Group{
	// 套差与趋势跟随
	Group{
		// 需要追加基差回归的假设, 在主逻辑影响下
		// 为什么对赌是 多沥青 空燃油 的假设, 为什么不能倒过来, 需要解释
		// 增加相对锚定价差 2118~2090 : 1622
		/*
			建仓反套，建仓反向反套，相对低点锚定相对基差距离，相对基差回归
			20~40个点, 如果不转跨日, 那就一定要平掉, 一天最多两次机会, 谨记

			单边机会还没出现, bu2012=2000, fu2101=1500
		*/
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
}
