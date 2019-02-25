package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var dataSet = []map[string]string{
	map[string]string{
		"name":               "上证指数",
		"stockCode":          "000001",
		"affiliatedExchange": "1",
	},
	map[string]string{
		"name":               "深证成指",
		"stockCode":          "399001",
		"affiliatedExchange": "2",
	},
	map[string]string{
		"name":               "上证50",
		"stockCode":          "000016",
		"affiliatedExchange": "1",
	},
	map[string]string{
		"name":               "沪深300",
		"stockCode":          "000300",
		"affiliatedExchange": "1",
	},
	map[string]string{
		"name":               "A股指数",
		"stockCode":          "000002",
		"affiliatedExchange": "1",
	},
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(5)
	for _, tmp := range dataSet {
		stockCode := tmp["stockCode"]
		affiliatedExchange := tmp["affiliatedExchange"]
		name := tmp["name"]
		log.Println(name, stockCode, affiliatedExchange)

		go func() {
			defer func() {
				wg.Done()
			}()
			var b bool
			for !b {
				b = Get(stockCode, affiliatedExchange, name)
				time.Sleep(time.Millisecond * 1000)
			}
		}()
		time.Sleep(time.Millisecond * 200)
	}

	wg.Wait()
	log.Println("succ")

}

func Get(stockCode, ae, name string) bool {
	url := "http://pdfm.eastmoney.com/EM_UBG_PDTI_Fast/api/js?rtntype=5&token=4f1862fc3b5e77c150a2b985b12db0fd&cb=jQuery18306500869817051207_1551101184728&id=" + stockCode + ae + "&type=k&authorityType=&_=1551101309291"

	resp, err := http.Get(url)
	if err != nil {
		log.Println(name, err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(name, err)
		return false
	}

	str := string(body)
	ret := strings.Split(strings.Split(str, "data\":")[1], "})")

	var infos []interface{}
	err = json.Unmarshal([]byte(ret[0]), &infos)
	if err != nil {
		log.Println("uncoded json", name, err)
		return false
	}

	fs, err := os.OpenFile("./data/index_details/"+stockCode, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Println(name, err)
		return false
	}
	defer fs.Close()

	for _, tmp := range infos {
		str := tmp.(string)
		fs.WriteString(str + "\r\n")
	}

	return true
}
