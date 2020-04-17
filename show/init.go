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
	/*
		受巴西大豆, 美大豆到港量影响, 美大豆价格不断近期新低
		进口国外猪肉打压国内猪肉价格影响
		短期豆粕低位横盘，甚至会在破新低, 目标01合约, 2720
	*/
	"m2101": &Variety{
		Name:          "豆",
		OriginDataUrl: "114_m2101_qt?callbackName=aa&cb=aa&_=1585752611719",
		SpotPrice:     3182.5,
		Amount:        16,
		Aims:          3100,
		Trend:         "2830,多",
		Describe:      "周线看还未启动, 多",
		Level:         2,
		IsShow:        false,
	},
	"MA009": &Variety{
		Name:          "醇",
		OriginDataUrl: "115_MA009_qt?callbackName=aa&cb=aa&_=1586842655155",
		SpotPrice:     1760,
		Amount:        0,
		Aims:          0,
		Trend:         "",
		Describe:      "原油带着化工崩盘了, ..., 砍仓",
		Level:         1,
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
		IsShow:        false,
	},
	"JD2009": &Variety{
		Name:          "蛋",
		OriginDataUrl: "114_jd2009_qt?callbackName=aa&cb=aa&_=1586162496980",
		SpotPrice:     3070,
		Amount:        0,
		Aims:          0,
		Trend:         "4340,空",
		Describe:      "金叉, 但是下影线两根, 预期需求不足, 高空",
		Level:         3,
		IsShow:        false,
	},
}
