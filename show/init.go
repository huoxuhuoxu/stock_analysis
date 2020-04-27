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

// 模式, 配比单位, 基差临界线, 价格结算单位, 划点单位, 建仓平仓提示
// 可能不能基于现在这个来做, 需要在拓展, 还要增加 比例保证金消耗

// 期货是给远期现货定价的, 而不是用现货现在的价格给期货定价, 谨记
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
		Name:          "原油",
		OriginDataUrl: "142_sc2006_qt?callbackName=aa&cb=aa&_=1587309656220",
		SpotPrice:     0,
		Aims:          "",
		Describe: `
		`,
		Level:          4,
		PricePrecision: 1,
		IsShow:         true,
	},
	"sc2009": &Variety{
		Name:          "原油09",
		OriginDataUrl: "142_sc2009_qt?callbackName=aa&cb=aa&_=1587309656220",
		SpotPrice:     0,
		Aims:          "",
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
	"sc2012": &Variety{
		Name:          "原油",
		OriginDataUrl: "142_sc2012_qt?callbackName=aa&cb=aa&_=1587309656220",
		SpotPrice:     0,
		Aims:          "",
		Describe: `
		`,
		Level:          4,
		PricePrecision: 1,
		IsShow:         true,
	},
	"rb2101": &Variety{
		Name:          "螺纹",
		OriginDataUrl: "113_rb2101_qt?callbackName=aa&cb=aa&_=1587308954178",
		SpotPrice:     3405.5,
		Aims:          "",
		Describe: `
			供需强劲, 但价格上涨受制于高库存
			关联原料铁矿石也很强劲, 进入平台区, 一旦突破, 3400指日可待
			反之, 3100也很容易回踩, 3250为中线, 上下各150点宽幅震荡
		`,
		Level:          3,
		PricePrecision: 0,
		IsShow:         true,
	},
	"i2009": &Variety{
		Name:          "铁矿石",
		OriginDataUrl: "114_i2009_qt?callbackName=aa&cb=aa&_=1587396948222",
		SpotPrice:     663.89,
		Aims:          "",
		Describe: `
			减产, 但是四五月份到港量高压制价格继续上行
			可以与螺纹联做, 日内 多铁矿空螺纹
			目前 铁矿 强于 螺纹, 1:2 或 2:3 建仓, 不跨日持仓
		`,
		Level:          3,
		PricePrecision: 1,
		IsShow:         true,
	},
	"m2101": &Variety{
		Name:          "豆粕",
		OriginDataUrl: "114_m2101_qt?callbackName=aa&cb=aa&_=1585752611719",
		SpotPrice:     3116.25,
		Aims:          "2720",
		Describe: `
			受巴西大豆, 美大豆到港量影响, 美大豆价格不断近期新低
			进口国外猪肉打压国内猪肉价格影响养殖猪产生的利润
			短期豆粕低位横盘，甚至会在破新低, 目标01合约, 2720
			近期出现低点2750+, 开始反转的样子
			五月合约交割价, 预估3000以上, 目前200+基点, 或有补基差行情
		`,
		Level:          4,
		PricePrecision: 0,
		IsShow:         true,
	},
	"bu2012": &Variety{
		Name:          "沥青",
		OriginDataUrl: "113_bu2012_qt?callbackName=aa&cb=aa&_=1587532980754",
		SpotPrice:     2252.5,
		Aims:          "",
		Describe: `
			从一个足够长的周期来看, 价格会回归相对的内在价值与正常供需所在的价格
			长期看涨, 但短期严重受原油行情的压制
			不做单边观点, 与燃油做反套
			多沥青空燃油, 只做日内, 预期一个交易日有一到两次机会产生20～30个点的基差
			不跨日持仓, 日内平仓或者锁定仓位, 不赚跨日的钱, 防止大的风险, 只赚日内价差
		`,
		Level:          2,
		PricePrecision: 0,
		IsShow:         true,
	},
	"fu2101": &Variety{
		Name:          "燃油",
		OriginDataUrl: "113_fu2101_qt?callbackName=aa&cb=aa&_=1587533000837",
		SpotPrice:     0,
		Aims:          "",
		Describe: `
			航运逻辑受疫情影响较弱, 又受到原油行情的压制, 偏弱运行
			不做单边观点, 与沥青联合做
			逻辑, 同为原油相关性极高的衍生品 沥青的副逻辑基建 强于 燃油的副逻辑航运
		`,
		Level:          2,
		PricePrecision: 0,
		IsShow:         true,
	},
	"ag2012": &Variety{
		Name:          "白银",
		OriginDataUrl: "113_ag2012_qt?callbackName=aa&cb=aa&_=1587533190507",
		SpotPrice:     0,
		Aims:          "",
		Describe: `
			技术面受压制, 3650 可空, 目标3380
			但受黄金影响, 容易出黑天鹅, 观察为主, 不动
		`,
		Level:          3,
		PricePrecision: 0,
		IsShow:         true,
	},
	"au2012": &Variety{
		Name:          "沪金",
		OriginDataUrl: "113_au2012_qt?callbackName=aa&cb=aa&_=1587566439644",
		SpotPrice:     0,
		Aims:          "",
		Describe: `
			围绕 360 上下 30 个点震荡, 跟着外盘跑
			技术面上看在 380/390 受压制, 但
			风险较大, 伦敦黄金可能存在逼仓的可能, 导致外盘黄金突破新高
			进而引发内盘黄金突破平台区, 观察为主, 不动
		`,
		Level:          3,
		PricePrecision: 2,
		IsShow:         true,
	},
	"jd2009": &Variety{
		Name:          "鸡蛋",
		OriginDataUrl: "114_jd2009_qt?callbackName=aa&cb=aa&_=1587566503238",
		SpotPrice:     3142,
		Aims:          "4100~4300",
		Describe: `
			看震荡行情
			4100以下可多, 4300以上可空
		`,
		Level:          4,
		PricePrecision: 0,
		IsShow:         true,
	},
	"a2009": &Variety{
		Name:          "豆一",
		OriginDataUrl: "114_a2009_qt?callbackName=aa&cb=aa&_=1587566549527",
		SpotPrice:     4666.67,
		Aims:          "4800~5000",
		Describe: `
			妖的一逼, 现货价格强劲, 没有进行扩产, 交易所交割标准上升
			如上原因构建一个强大的豆一多头行情, 并且大波动大回踩, 但上升趋势不改
			不做, 纯粹就是验证想法, 为我近期最大的亏损, 把它盯完
			其实吧, 很可能五月多头移仓九月, 继续逼仓
		`,
		Level:          4,
		PricePrecision: 0,
		IsShow:         true,
	},
}

/*
	请严格执行
*/
var dailyOperation = `
	
`

/*
	找时间盯一下 豆油/菜油/棕榈油
	看有没有机会在里面做反套
	反套不锁仓, 只平仓, 尽量不要跨夜,
	其实逻辑错的话跨夜, 第二天开盘 风险对冲 也还是平的

	每次出手的目标 200~400利润/手, 不要贪, 贪了会反转, 反而被套里面, 不值得
	观察每天走势, 防止强弱反转引发基差错位
*/
