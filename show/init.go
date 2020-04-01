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
	Name                     string  // 商品名称
	Code                     string  // 商品编号
	OriginDataUrl            string  // 数据源头
	Price                    float64 // 当前价格
	Value                    float64 // 涨跌值
	SpotPrice                float64 // 现货价格
	Amount                   uint8   // 预计建仓数
	Aims                     float64 // 目标价位
	TmpHigh                  float64 // 暂时预计高位
	TmpLow                   float64 // 暂时预计低位
	VolatilityValue          float64 // 当前波动值(high-low)
	RemainingVolatilityValue float64 // 剩余波动空间(近期每日波动 * 交割前剩余交易日)
	IsShow                   bool    // 是否获取及输出
}

var varietys = map[string]*Variety{
	"m2005": &Variety{
		Name:                     "豆粕2005",
		Code:                     "m2005",
		OriginDataUrl:            "114_m2005_qt?callbackName=aa&cb=aa&_=1585752611719",
		SpotPrice:                3327,
		Amount:                   20,
		Aims:                     3100,
		TmpHigh:                  2950,
		TmpLow:                   2800,
		RemainingVolatilityValue: 60 * 30,
		IsShow:                   true,
	},
	"rb2101": &Variety{
		Name:                     "螺纹2101",
		Code:                     "rb2101",
		OriginDataUrl:            "113_rb2101_qt?callbackName=aa&cb=aa&_=1585753542283",
		SpotPrice:                3388,
		Amount:                   0,
		Aims:                     3400,
		TmpHigh:                  0,
		TmpLow:                   2900,
		RemainingVolatilityValue: 0,
		IsShow:                   false,
	},
	"c2009": &Variety{
		Name:                     "玉米2009",
		Code:                     "c2009",
		OriginDataUrl:            "114_c2009_qt?callbackName=aa&cb=aa&_=1585757187349",
		SpotPrice:                1870,
		Amount:                   10,
		Aims:                     2100,
		TmpHigh:                  0,
		TmpLow:                   0,
		RemainingVolatilityValue: 0,
		IsShow:                   true,
	},
	"fu2101": &Variety{
		Name:                     "燃油2101",
		Code:                     "fu2101",
		OriginDataUrl:            "113_fu2101_qt?callbackName=aa&cb=aa&_=1585757527911",
		SpotPrice:                0,
		Amount:                   10,
		Aims:                     1500,
		TmpHigh:                  1800,
		TmpLow:                   1600,
		RemainingVolatilityValue: 0,
		IsShow:                   true,
	},
	"CF009": &Variety{
		Name:                     "棉花2009",
		Code:                     "CF009",
		OriginDataUrl:            "115_CF009_qt?callbackName=aa&cb=aa&_=1585757766042",
		SpotPrice:                11219,
		Amount:                   0,
		Aims:                     9000,
		TmpHigh:                  11000,
		TmpLow:                   10000,
		RemainingVolatilityValue: 0,
		IsShow:                   true,
	},
}
