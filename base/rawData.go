package base

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// ------------------------------------ 获取实时基础数据

func GetDatas() map[string]*GeneralData {
	gds := make(map[string]*GeneralData)

	// 获取原始数据
	initUrl := "http://15.push2.eastmoney.com/api/qt/clist/get?cb=jQuery112407436525408153816_1606727467774&pn=1&pz=4000&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&fid=f3&fs=m:0+t:6,m:0+t:13,m:0+t:80,m:1+t:2,m:1+t:23&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152&_=1606727467775"
	header := make(map[string]string)
	header["Accept"] = "*/*"
	header["Accept-Encoding"] = "gzip, deflate"
	header["Accept-Language"] = "zh-CN,zh;q=0.9,en;q=0.8"
	header["Cache-Control"] = "no-cache"
	header["Connection"] = "keep-alive"
	header["Pragma"] = "no-cache"
	header["Host"] = "15.push2.eastmoney.com"
	header["Referer"] = "http://quote.eastmoney.com/"
	header["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36"
	header["Cookie"] = "m_hq_fls=js; intellpositionL=864px; HAList=a-sh-601012-%u9686%u57FA%u80A1%u4EFD%2Ca-sz-002475-%u7ACB%u8BAF%u7CBE%u5BC6%2Ca-sz-300136-%u4FE1%u7EF4%u901A%u4FE1%2Ca-sh-600736-%u82CF%u5DDE%u9AD8%u65B0%2Ca-sz-000998-%u9686%u5E73%u9AD8%u79D1%2Ca-sz-002961-%u745E%u8FBE%u671F%u8D27%2Ca-sz-300568-%u661F%u6E90%u6750%u8D28%2Ca-sz-002511-%u4E2D%u987A%u6D01%u67D4%2Ca-sz-002046-%u8F74%u7814%u79D1%u6280%2Ca-sz-000768-%u4E2D%u822A%u98DE%u673A%2Ca-sz-002400-%u7701%u5E7F%u96C6%u56E2%2Ca-sh-600516-%u65B9%u5927%u70AD%u7D20; cowminicookie=true; qgqp_b_id=ae50a0d1bf559adfe49a7af4fc928135; intellpositionT=978px; st_si=37791226608855; st_asi=delete; st_pvi=77027473738396; st_sp=2018-09-19%2022%3A43%3A31; st_inirUrl=https%3A%2F%2Fwww.baidu.com%2Flink; st_sn=11; st_psi=20201130164956448-113200301321-0792700827"

	client := &http.Client{}
	req, err := http.NewRequest("GET", initUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	for key, v := range header {
		req.Header.Add(key, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 解压缩
	reader, err := gzip.NewReader(bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	tmpBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	// 读取数据
	rawDataStr := string(tmpBytes)
	originDataStr := strings.Split(strings.Split(rawDataStr, "(")[1], ")")[0]

	// 分析数据
	var data = make(map[string]interface{})
	err = json.Unmarshal([]byte(originDataStr), &data)
	if err != nil {
		log.Fatal(err)
	}

	tmpData := data["data"].(map[string]interface{})
	list := tmpData["diff"].([]interface{})
	for _, v := range list {
		g := v.(map[string]interface{})

		stockCode := g["f12"].(string)
		stockName := g["f14"].(string)

		// 排除长创业版和科创板
		if strings.HasPrefix(stockCode, "300") || strings.HasPrefix(stockCode, "688") {
			continue
		}
		// 排除停牌的和ST的, 暂时包含停牌的 ....
		if strings.HasPrefix(stockName, "*ST") || strings.HasPrefix(stockName, "ST") {
			continue
		}

		currentPrice, _ := g["f2"].(float64)
		upAndDownRange, _ := g["f3"].(float64)
		turnover, _ := g["f6"].(float64)
		high, _ := g["f15"].(float64)
		low, _ := g["f16"].(float64)
		open, _ := g["f17"].(float64)
		close, _ := g["f18"].(float64)
		min5, _ := g["f11"].(float64)
		changehands, _ := g["f8"].(float64)
		circulation, _ := g["f21"].(float64)
		amplitude, _ := g["f7"].(float64)

		gds[stockCode] = &GeneralData{
			StockType:      strconv.FormatFloat(g["f1"].(float64), 'f', -1, 64),
			StockCode:      stockCode,
			StockName:      stockName,
			CurrentPrice:   currentPrice,
			UpAndDownRange: upAndDownRange,
			Turnover:       turnover / 100000000,
			High:           high,
			Low:            low,
			Open:           open,
			Close:          close,
			Min5:           min5,
			Changehands:    changehands,
			Circulation:    circulation / 100000000,
			Amplitude:      amplitude,
		}
	}

	return gds
}

// ------------------------------------- 获取实时资金数据
func GetMoneyData() map[string]*GeneralData {
	gds := GetDatas()

	// 获取原始数据
	header := make(map[string]string)
	header["Accept"] = "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
	header["Accept-Encoding"] = "gzip, deflate"
	header["Accept-Language"] = "zh-CN,zh;q=0.9,en;q=0.8"
	header["Cache-Control"] = "n-cache"
	header["Connection"] = "keep-alive"
	header["Host"] = "nufm.dfcfw.com"
	header["Pragma"] = "no-cache"
	header["Upgrade-Insecure-Requests"] = "1"
	header["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"

	client := &http.Client{}
	// 只抓取主力资金流入的前3000支票
	req, err := http.NewRequest("GET", "http://push2.eastmoney.com/api/qt/clist/get?pn=1&pz=3000&po=1&np=1&ut=b2884a393a59ad64002292a3e90d46a5&fltt=2&invt=2&fid0=f4001&fid=f62&fs=m:0+t:6+f:!2,m:0+t:13+f:!2,m:0+t:80+f:!2,m:1+t:2+f:!2,m:1+t:23+f:!2,m:0+t:7+f:!2,m:1+t:3+f:!2&stat=1&fields=f12,f14,f2,f3,f62,f184,f66,f69,f72,f75,f78,f81,f84,f87,f204,f205,f124&rt=52415130&cb=jQuery183019370468044060885_1572453915798&_=1572453916366", nil)
	if err != nil {
		log.Fatal(err)
	}
	for key, v := range header {
		req.Header.Add(key, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	reader, err := gzip.NewReader(bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	tmpBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	// 读取数据
	rawDataStr := string(tmpBytes)
	originDataStr := strings.Split(strings.Split(rawDataStr, "(")[1], ")")[0]

	// 解析数据
	var data = make(map[string]interface{})
	err = json.Unmarshal([]byte(originDataStr), &data)
	if err != nil {
		log.Fatal(err)
	}

	tmpData := data["data"].(map[string]interface{})
	list := tmpData["diff"].([]interface{})

	// 重组数据
	for _, v := range list {
		g := v.(map[string]interface{})
		if gd, ok := gds[g["f12"].(string)]; ok {
			gd.MainV = g["f62"].(float64) / 100000000
			gd.CBigV = g["f66"].(float64)
			gd.BigV = g["f72"].(float64)
			gd.MiddleV = g["f78"].(float64)
			gd.SmallV = g["f84"].(float64)
			// gd.MainP = g["f69"].(float64) // 超大单流入
			gd.MainP = g["f184"].(float64) // 主力流入
		}
	}

	// 删除无需使用的股
	for k, gd := range gds {
		if gd.MainV == 0 {
			delete(gds, k)
		}
	}

	return gds
}
