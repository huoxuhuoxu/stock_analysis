package main

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
	initUrl := "http://nufm.dfcfw.com/EM_Finance2014NumericApplication/JS.aspx?cb=jQuery112407378636087109598_1528454080519&type=CT&token=4f1862fc3b5e77c150a2b985b12db0fd&js=(%7Bdata%3A%5B(x)%5D%2CrecordsTotal%3A(tot)%2CrecordsFiltered%3A(tot)%7D)&cmd=C._A&sty=FCOIATC&st=(ChangePercent)&sr=-1&p=1&ps=10000&_=1528454080520"
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

	// 读取数据
	reader, err := gzip.NewReader(bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	tmpBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	// 解析数据
	str := string(tmpBytes)
	ret := strings.Split(strings.Split(str, "data:[\"")[1], "\"],")

	strData := ret[0]
	datas := strings.Split(strData, "\",\"")

	for _, data := range datas {
		arr := strings.Split(data, ",")

		// 排除长创业版和科创板
		if strings.HasPrefix(arr[1], "300") || strings.HasPrefix(arr[1], "688") {
			continue
		}
		// 排除停牌的和ST的
		if arr[4] == "-" || strings.HasPrefix(arr[2], "*ST") || strings.HasPrefix(arr[2], "ST") {
			continue
		}

		currentPrice, _ := strconv.ParseFloat(arr[3], 64)
		upAndDownRange, _ := strconv.ParseFloat(arr[5], 64)
		turnover, _ := strconv.ParseFloat(arr[7], 64)
		high, _ := strconv.ParseFloat(arr[9], 64)
		low, _ := strconv.ParseFloat(arr[10], 64)
		open, _ := strconv.ParseFloat(arr[11], 64)
		close, _ := strconv.ParseFloat(arr[12], 64)
		min5, _ := strconv.ParseFloat(arr[13], 64)
		changehands, _ := strconv.ParseFloat(arr[15], 64)
		circulation, _ := strconv.ParseFloat(arr[19], 64)

		gds[arr[1]] = &GeneralData{
			stockType:      arr[0],
			stockCode:      arr[1],
			stockName:      arr[2],
			currentPrice:   currentPrice,
			upAndDownRange: upAndDownRange,
			turnover:       turnover / 100000000,
			high:           high,
			low:            low,
			open:           open,
			close:          close,
			min5:           min5,
			changehands:    changehands,
			circulation:    circulation / 100000000,
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
			gd.mainV = g["f62"].(float64) / 100000000
			gd.cBigV = g["f66"].(float64)
			gd.bigV = g["f72"].(float64)
			gd.middleV = g["f78"].(float64)
			gd.smallV = g["f84"].(float64)
			gd.mainP = g["f69"].(float64)
		}
	}

	// 删除无需使用的股
	for k, gd := range gds {
		if gd.mainV == 0 {
			delete(gds, k)
		}
	}

	return gds
}
