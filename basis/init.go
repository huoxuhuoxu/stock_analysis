package main

// 命名 - 上帝之手

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
			basis > BOUNDARY: 反向
			other: 等待
			Profit: 自己把握
	*/
	actions = []string{"等待/平仓", "开仓", "反向"}
)

const (
	ORIGIN_URL = "http://futsse.eastmoney.com/static/"
	BOUNDARY   = 30 // 反向边界
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
	Matching    [2]int    // 配比
	// 基差临界点数, 预计中可能会存在超跌值, 以此值为起点开始套利, 最关键的一个值
	Limit             float64
	Profit            float64    // 利润点数, 不包含划点, 至少 4% 回报, 实际平仓获利需要叠加划点损失
	MarginConsumption string     // 组合需要消耗保证金
	Level             int        // 优先级, 99: 待验证待的逻辑, 100: 等待
	Describe          string     // 组合逻辑说明
	IsAll             bool       // 是否输出其余因子
	ReasonablePrice   [2]float64 // 相对点位
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

// 真要做黄金/白银的话, 也需要加入反向系统, 也可能无法做套, 只能做跟随的逻辑
// 原油组, 套, 取决于基本面强弱, 并且要看开盘后能否走出基差点位, 有待观察
// 反套组
var groups = []Group{
	// 套差与趋势跟随
	Group{
		/*
			20~40个点, 一个时间段一次机会, 把握住, 每天能够把握住一次就够了
			不跨日, 谨记
			单边机会还没出现, bu2012=2000, fu2101=1500
		*/
		Name:              "多沥青/空燃油",
		Combination:       [2]string{"bu2012", "fu2101"},
		Matching:          [2]int{1, 1},
		Limit:             10,
		MarginConsumption: "5",
		Level:             1,
		Profit:            20,
		Describe: `
			在原油为主导因素情况下, 
			利用多空对冲对与原油具有高度相关性的衍生品进行对赌,
			赌 沥青副逻辑-基建 强于 燃油副逻辑-航运,
			不论涨或者跌, 在盘中体现, 沥青强于燃油, 
			日内会产生 20～40 点的套利机会
			但是在连续的沥青强于燃油的情况下, 导致这两个品种的基差过于放大,
			那么接下来就会被修复, 基差回归, 定义 锚定相对价格, 推 基差回归 的相对点位
		`,
		IsAll:           true,
		ReasonablePrice: [2]float64{2118, 1622},
	},
	/*
		06/09 没有走出真正意义上的盘中分化行情, 很同步
		06/12, 09/12 走出了分化行情, 12盘中强于 06/09 盘中
		不管开盘价格, 不管开盘基差, 只看盘中能够分化的基差, 这部分才是我要挣的
		所以, 如果要做原油的 远月/近月 对冲套利的话, 12合约是需要多的, 至于空 06还是09 各有利弊, 再看看
	*/
	Group{
		Name:              "多远月/空近月",
		Combination:       [2]string{"sc2009", "sc2006"},
		Matching:          [2]int{1, 1},
		Limit:             100,
		MarginConsumption: "60",
		Level:             2,
		Profit:            240,
		Describe:          ``,
	},
	Group{
		Name:              "多远月/空近月",
		Combination:       [2]string{"sc2012", "sc2006"},
		Matching:          [2]int{1, 1},
		Limit:             100,
		MarginConsumption: "60",
		Level:             2,
		Profit:            240,
		Describe:          ``,
	},
	Group{
		Name:              "多远月/空近月",
		Combination:       [2]string{"sc2012", "sc2009"},
		Matching:          [2]int{1, 1},
		Limit:             100,
		MarginConsumption: "60",
		Level:             2,
		Profit:            240,
		Describe:          ``,
	},
	Group{
		Name:              "多黄金/空白银",
		Combination:       [2]string{"au2012", "ag2012"},
		Matching:          [2]int{1, 3},
		Limit:             50,
		MarginConsumption: "50",
		Level:             999,
		Profit:            200,
		Describe: `
			反套逻辑弱, 不能做
			但是, 白银跟着黄金走, 可以做参考
			比如 04/27, 黄金盘中下跌, 白银暂时坚挺, 后跟跌
			可以用于做提前预判
		`,
	},
	Group{
		Name:              "多黄金/空白银",
		Combination:       [2]string{"au2012", "ag2012"},
		Matching:          [2]int{1, 5},
		Limit:             50,
		MarginConsumption: "60",
		Level:             999,
		Profit:            240,
		Describe:          ``,
	},
}

// 合约集合
var varietys = map[string]*Variety{
	"bu2012": &Variety{
		Name:            "沥青",
		OriginDataUrl:   "113_bu2012_qt?callbackName=aa&cb=aa&_=1587532980754",
		PricePrecision:  0,
		Dash:            2,
		DashCoefficient: 10,
	},
	"fu2101": &Variety{
		Name:            "燃油",
		OriginDataUrl:   "113_fu2101_qt?callbackName=aa&cb=aa&_=1587533000837",
		PricePrecision:  0,
		Dash:            1,
		DashCoefficient: 10,
	},
	"au2012": &Variety{
		Name:            "沪金",
		OriginDataUrl:   "113_au2012_qt?callbackName=aa&cb=aa&_=1587566439644",
		PricePrecision:  2,
		Dash:            0.02,
		DashCoefficient: 1000,
	},
	"ag2012": &Variety{
		Name:            "白银",
		OriginDataUrl:   "113_ag2012_qt?callbackName=aa&cb=aa&_=1587533190507",
		PricePrecision:  0,
		Dash:            1,
		DashCoefficient: 15,
	},
	"sc2006": &Variety{
		Name:            "原油",
		OriginDataUrl:   "142_sc2006_qt?callbackName=aa&cb=aa&_=1587309656220",
		PricePrecision:  1,
		Dash:            0.1,
		DashCoefficient: 1000,
	},
	"sc2009": &Variety{
		Name:            "原油",
		OriginDataUrl:   "142_sc2009_qt?callbackName=aa&cb=aa&_=1587309656220",
		PricePrecision:  1,
		Dash:            0.1,
		DashCoefficient: 1000,
	},
	"sc2012": &Variety{
		Name:            "原油",
		OriginDataUrl:   "142_sc2012_qt?callbackName=aa&cb=aa&_=1587309656220",
		PricePrecision:  1,
		Dash:            0.1,
		DashCoefficient: 1000,
	},
}
