package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	MONEY_FLOW_URL = "http://push2.eastmoney.com/api/qt/clist/get?pn=1&pz=50&po=1&np=1&ut=b2884a393a59ad64002292a3e90d46a5&fltt=2&invt=2&fid0=f4001&fid=f62&fs=m:0+t:6+f:!2,m:0+t:13+f:!2,m:0+t:80+f:!2,m:1+t:2+f:!2,m:1+t:23+f:!2,m:0+t:7+f:!2,m:1+t:3+f:!2&stat=1&fields=f12,f14,f2,f3,f62,f184,f66,f69,f72,f75,f78,f81,f84,f87,f204,f205,f124&rt=52415130&cb=jQuery183019370468044060885_1572453915798&_=1572453916366"
)

var (
	cje  float64 // 成交额
	ltzb float64 // 流通占比
)

func init() {
	flag.Float64Var(&cje, "cje", 3.0, "成交额")
	flag.Float64Var(&ltzb, "ltzb", 0.004, "流通占比")
	flag.Parse()
	log.SetFlags(0)
}

func main() {
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
	req, err := http.NewRequest("GET", MONEY_FLOW_URL, nil)
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

	rawDataStr := string(tmpBytes)
	originDataStr := strings.Split(strings.Split(rawDataStr, "(")[1], ")")[0]

	var data = make(map[string]interface{})
	err = json.Unmarshal([]byte(originDataStr), &data)
	if err != nil {
		log.Fatal(err)
	}

	tmpData := data["data"].(map[string]interface{})
	list := tmpData["diff"].([]interface{})

	lll := getL()

	fmt.Printf("%4s, %4s, %4s, %4s, %4s, %4s, %4s, %4s \r\n",
		"代码", "名称", "涨跌幅", "净流入", "成交占比", "成交额", "流通市值", "流通占比")
	for _, v := range list {
		g := v.(map[string]interface{})

		// 股票代码
		code := g["f12"].(string)
		if strings.HasPrefix(code, "300") || strings.HasPrefix(code, "688") {
			continue
		}

		infos := lll[code]

		// 资金流入意向
		money := (g["f62"].(float64)) / 100000000

		// 流通市值
		tmpl, _ := strconv.ParseFloat(infos[19], 64)
		tmpl = tmpl / 100000000

		// 主力流入在流通占比
		tmpll := money / tmpl

		// 成交额
		tmpl2, _ := strconv.ParseFloat(infos[7], 64)
		tmpl2 = tmpl2 / 100000000

		// 筛选
		if tmpll < ltzb || tmpl2 < cje {
			continue
		}

		fmt.Printf("%s, %s, %.2f%%, %.2f亿, %.2f%%, %.2f亿, %.2f亿, %.4f%%",
			code,
			g["f14"].(string),
			g["f3"].(float64),
			money,
			money/tmpl2,
			tmpl2,
			tmpl,
			tmpll,
		)

		fmt.Println("")
	}

}

func getL() map[string][]string {
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

	reader, err := gzip.NewReader(bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	tmpBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	str := string(tmpBytes)
	ret := strings.Split(strings.Split(str, "data:[\"")[1], "\"],")

	strData := ret[0]
	datas := strings.Split(strData, "\",\"")

	var lll = make(map[string][]string)
	for _, data := range datas {
		arr := strings.Split(data, ",")
		lll[arr[1]] = arr
	}

	return lll
}
