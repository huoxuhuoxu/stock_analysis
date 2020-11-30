package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	URL = "http://98.push2his.eastmoney.com/api/qt/stock/kline/get?cb=jQuery1124038063739276411046_1606743014143&secid=1.510050&ut=fa5fd1943c7b386f172d6893dbfba10b&fields1=f1%2Cf2%2Cf3%2Cf4%2Cf5%2Cf6&fields2=f51%2Cf52%2Cf53%2Cf54%2Cf55%2Cf56%2Cf57%2Cf58%2Cf59%2Cf60%2Cf61&klt=101&fqt=0&end=20500101&lmt=3000&_=1606743014166"
)

func main() {
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 读取数据
	rawDataStr := string(body)
	originDataStr := strings.Split(strings.Split(rawDataStr, "(")[1], ")")[0]

	// 分析数据
	var data = make(map[string]interface{})
	err = json.Unmarshal([]byte(originDataStr), &data)
	if err != nil {
		log.Fatal(err)
	}

	tmpData := data["data"].(map[string]interface{})
	list := tmpData["klines"].([]interface{})

	//
	var yearsMap = make(map[string][]string)
	for _, v := range list {
		s := v.(string)
		year := s[0:4]
		if arr, ok := yearsMap[year]; ok {
			arr = append(arr, s)
			yearsMap[year] = arr
		} else {
			yearsMap[year] = []string{s}
		}
	}

	for k, arr := range yearsMap {
		log.Println(k)
		fs, _ := os.OpenFile(fmt.Sprintf("./ETF_50/%s.txt", k), os.O_CREATE|os.O_WRONLY, 0666)
		for _, v := range arr {
			fs.WriteString(v + "\n")
		}
		fs.Close()
	}
}
