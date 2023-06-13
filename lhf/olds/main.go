package main

import (
	"net/http"
	"io/ioutil"
	"log"
	"strings"
	"encoding/json"
	"os"
	"time"
)


var (
	base_stock = []string{
		"http://31.push2his.eastmoney.com/api/qt/stock/kline/get?cb=jQuery35109319016647768255_1686540747303&secid=",
		"&ut=fa5fd1943c7b386f172d6893dbfba10b&fields1=f1%2Cf2%2Cf3%2Cf4%2Cf5%2Cf6&fields2=f51%2Cf52%2Cf53%2Cf54%2Cf55%2Cf56%2Cf57%2Cf58%2Cf59%2Cf60%2Cf61&klt=103&fqt=1&beg=0&end=20500101&smplmt=460&lmt=1000000&_=1686540747348",
	}
)

func Get (url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("异常:", err)
		return false
	}

	str := string(body)
	arr := strings.Split(str, "({")
	if len(arr) < 2 {
		log.Println("异常:", arr)
		return false
	}

	json_str := "{" + arr[1][:len(arr[1])-2]
	var data = make(map[string]interface{})
	err = json.Unmarshal([]byte(json_str), &data)
	if err != nil {
		log.Println("异常:", err)
		return false
	}

	d := data["data"].(map[string]interface{})

	name := d["code"].(string)
	fs, err := os.OpenFile("./data_csv/" + name + ".csv", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Println("异常:", err)
		return false
	}
	defer fs.Close()

	fs.WriteString("日期,开盘价,收盘价,最高,最低,成交量,成交金额,振幅,涨跌幅,涨跌额,换手率\r\n")
	klines := d["klines"].([]interface{})
	for _, _tmp := range klines {
		kline := _tmp.(string)
		fs.WriteString(kline+"\r\n")
	}

	return true
}

func main (){
	fs, err := os.Open("./target_code.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()
	lines_bytes, err := ioutil.ReadAll(fs)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(lines_bytes), "\n")
	for _, line := range lines {
		_line := strings.TrimSpace(line)
		log.Println("处理中, ", _line)

		s_url := base_stock[0] + _line + base_stock[1]
		b := Get(s_url)
		if !b {
			log.Println("处理失败")
		}
		time.Sleep(time.Second * 2)
	}
}
