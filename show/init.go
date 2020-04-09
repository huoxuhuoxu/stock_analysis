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
	Name          string  // 商品名称
	OriginDataUrl string  // 数据源头
	Price         float64 // 当前价格
	Value         float64 // 涨跌值
	SpotPrice     float64 // 现货价格
	Amount        uint8   // 预计建仓数
	Aims          float64 // 长期目标价位
	Trend         string  // 操作趋势
	Describe      string  // 说明
	Level         int     // 优先级, 1: 优先事项, 2: 逐步进行 3: 可试可不试, 4: 观察, 等待, 有极好的机会可试
	IsShow        bool    // 是否获取及输出
}

var varietys = map[string]*Variety{
	"a2009": &Variety{
		Name:          "一",
		OriginDataUrl: "114_a2009_qt?callbackName=aa&cb=aa&_=1585880873855",
		SpotPrice:     4500,
		Amount:        10,
		Aims:          4150,
		Trend:         "试空",
		Describe:      "短期承压, 等周线金叉",
		Level:         1,
		IsShow:        true,
	},
	"m2101": &Variety{
		Name:          "豆",
		OriginDataUrl: "114_m2101_qt?callbackName=aa&cb=aa&_=1585752611719",
		SpotPrice:     3207.5,
		Amount:        16,
		Aims:          3100,
		Trend:         "2830,多",
		Describe:      "周线看还未启动, 多",
		Level:         2,
		IsShow:        true,
	},
	"c2009": &Variety{
		Name:          "玉",
		OriginDataUrl: "114_c2009_qt?callbackName=aa&cb=aa&_=1585757187349",
		SpotPrice:     1881.43,
		Amount:        8,
		Aims:          2100,
		Trend:         "2030,多",
		Describe:      "05交割有基差修复需求, 短期看回踩, 长期看趋势不变, 多",
		Level:         2,
		IsShow:        true,
	},
	"CF009": &Variety{
		Name:          "棉",
		OriginDataUrl: "115_CF009_qt?callbackName=aa&cb=aa&_=1585757766042",
		SpotPrice:     11327,
		Amount:        1,
		Aims:          9000,
		Trend:         "12000,试空",
		Describe:      "6/7月, 蝗灾预期 偏多, 近两年服贸市场 偏空, 短多长空",
		Level:         4,
		IsShow:        true,
	},
	"Y2009": &Variety{
		Name:          "油",
		OriginDataUrl: "114_y2009_qt?callbackName=aa&cb=aa&_=1586162165978",
		SpotPrice:     5690,
		Amount:        2,
		Aims:          5800,
		Trend:         "试多",
		Describe:      "金叉, 看一波反弹, 被棕榈带着跑了",
		Level:         3,
		IsShow:        true,
	},
	"JD2009": &Variety{
		Name:          "蛋",
		OriginDataUrl: "114_jd2009_qt?callbackName=aa&cb=aa&_=1586162496980",
		SpotPrice:     2924,
		Amount:        0,
		Aims:          0,
		Trend:         "4340,空",
		Describe:      "金叉, 但是下影线两根, 预期需求不足, 高空",
		Level:         3,
		IsShow:        true,
	},
	"CS2009": &Variety{
		Name:          "淀",
		OriginDataUrl: "114_cs2009_qt?callbackName=aa&cb=aa&_=1586350881813",
		SpotPrice:     2406,
		Amount:        6,
		Aims:          2580,
		Trend:         "2360,多",
		Describe:      "周线金叉, 大趋势向上, 长多",
		Level:         1,
		IsShow:        true,
	},

	// 看不懂的, 今年不做的, 要专注
	"rb2101": &Variety{
		Name:          "螺纹2101",
		OriginDataUrl: "113_rb2101_qt?callbackName=aa&cb=aa&_=1585753542283",
		SpotPrice:     3386,
		Amount:        4,
		Aims:          0,
		Trend:         "空",
		Describe:      "..., 震荡向下, 3000筑波底反弹",
		IsShow:        false,
	},
	"fu2101": &Variety{
		Name:          "燃油2101",
		OriginDataUrl: "113_fu2101_qt?callbackName=aa&cb=aa&_=1585757527911",
		SpotPrice:     0,
		Amount:        10,
		Aims:          1500,
		IsShow:        false,
	},
	"I2009": &Variety{
		Name:          "铁矿石2009",
		OriginDataUrl: "114_i2009_qt?callbackName=aa&cb=aa&_=1586161835923",
		SpotPrice:     650,
		Amount:        0,
		Aims:          0,
		Trend:         "沽",
		Describe:      "看空今年黑色, 短期不见底",
		IsShow:        false,
	},
	"AU2006": &Variety{
		Name:          "黄金2006",
		OriginDataUrl: "113_au2006_qt?callbackName=aa&cb=aa&_=1586162013163",
		SpotPrice:     363.56,
		Amount:        0,
		Aims:          0,
		Trend:         "沽",
		Describe:      "存粹就是想空, 但是7/8月美国国债, 要翻多",
		IsShow:        false,
	},
}
