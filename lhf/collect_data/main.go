package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	/// 单一历史
	_data_url_1 = "http://23.push2his.eastmoney.com/api/qt/stock/kline/get?cb=jQuery3510524062431730671_1686623015466&secid=%s.%s&ut=fa5fd1943c7b386f172d6893dbfba10b&fields1=f1,f2,f3,f4,f5,f6&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61&klt=101&fqt=1&end=20500101&lmt=%s&_=1686623015519"

	/// 指数历史
	_data_url_7 = "http://62.push2his.eastmoney.com/api/qt/stock/kline/get?cb=jQuery351023429915107462018_1688610668005&secid=%s.%s&ut=fa5fd1943c7b386f172d6893dbfba10b&fields1=f1,f2,f3,f4,f5,f6&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61&klt=101&fqt=1&end=20500101&lmt=%s&_=1688610668069"

	/// 全深
	_data_url_11 = "http://45.push2.eastmoney.com/api/qt/clist/get?cb=jQuery112408281517578650619_1686625536243&pn=1&pz=5000&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&wbp2u=|0|0|0|web&fid=f3&fs=m:0+t:6,m:0+t:80&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152&_=1686625536244"

	/// 全沪
	_data_url_12 = "http://99.push2.eastmoney.com/api/qt/clist/get?cb=jQuery112404511119106118804_1686625611258&pn=1&pz=5000&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&wbp2u=|0|0|0|web&fid=f3&fs=m:1+t:2,m:1+t:23&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152&_=1686625611259"
)

/// 结构
type StrockInfo struct {
	Url001 string        // 历史数据
	Name   string        // 名称
	Ticker string        // 代码
	ExFlag string        // 交易所识别号
	Klines []interface{} // K线

	Shiying        float64 // 市盈率
	GunDongShiying float64 // 滚动市盈率
	Shijing        float64 // 市净率
	Price          float64 // 最新价
	Huanshou       float64 // 换手
	Liangbi        float64 // 量比
}

/// 获取单只股票历史数据
func _get_history(_strock_info *StrockInfo) bool {
	resp, err := http.Get(_strock_info.Url001)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()

	// 解析
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("异常:", err)
		return false
	}
	arr := strings.Split(string(body), "({")
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

	// 组合
	d := data["data"].(map[string]interface{})
	_strock_info.Name = d["name"].(string)
	_strock_info.Klines = d["klines"].([]interface{})
	return true
}

/// 获取全系列
func _get_all(_url string) map[string]map[string]interface{} {
	resp, err := http.Get(_url)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()

	// 解析
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("异常:", err)
		return nil
	}
	arr := strings.Split(string(body), "({")
	if len(arr) < 2 {
		log.Println("异常:", arr)
		return nil
	}
	json_str := "{" + arr[1][:len(arr[1])-2]
	var data = make(map[string]interface{})
	err = json.Unmarshal([]byte(json_str), &data)
	if err != nil {
		log.Println("异常:", err)
		return nil
	}

	d := data["data"].(map[string]interface{})
	ds := d["diff"].([]interface{})
	mm := make(map[string]map[string]interface{})
	for _, _tmp := range ds {
		_x := _tmp.(map[string]interface{})
		mm[_x["f12"].(string)] = _x
	}
	return mm
}

/// 写入csv
func _write_klines_to_csv(_strock_info *StrockInfo) bool {
	fs, err := os.OpenFile("./data_csv/"+_strock_info.Ticker+".csv", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Println("异常:", err)
		return false
	}
	defer fs.Close()

	fs.WriteString("日期,开盘价,收盘价,最高,最低,成交量,成交金额,振幅,涨跌幅,涨跌额,换手率\r\n")
	for _, _tmp := range _strock_info.Klines {
		kline := _tmp.(string)
		fs.WriteString(kline + "\r\n")
	}
	return true
}

/// 写入csv02
func _write_info_to_csv(_stocks map[string]*StrockInfo) bool {
	fs, err := os.OpenFile("./data_csv/result_code.csv", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Println("异常:", err)
		return false
	}
	defer fs.Close()

	fs.WriteString("股票代码,股票名称,动态市盈率,滚动市盈率,市净率,最新价,换手率,量比\r\n")
	for _, _stock := range _stocks {
		_s := fmt.Sprintf("%s,%s,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f", _stock.Ticker, _stock.Name, _stock.Shiying, _stock.GunDongShiying, _stock.Shijing, _stock.Price, _stock.Huanshou, _stock.Liangbi)
		fs.WriteString(_s + "\r\n")
	}
	return true
}

func main() {
	// 读取目标文件
	fs, err := os.Open("./target_code.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()
	lines_bytes, err := ioutil.ReadAll(fs)
	if err != nil {
		log.Fatal(err)
	}

	// 生成结构对象进行操作
	stocks := make(map[string]*StrockInfo)
	lines := strings.Split(string(lines_bytes), "\n")
	for _, line := range lines {
		params := strings.Split(line, ",")
		_ticker := strings.TrimSpace(params[1])
		_exflag := strings.TrimSpace(params[0])
		_lmt := strings.TrimSpace(params[2])
		stocks[_ticker] = &StrockInfo{
			Url001: fmt.Sprintf(_data_url_1, _exflag, _ticker, _lmt),
			Name:   "",
			Ticker: _ticker,
			ExFlag: _exflag,
			Klines: nil,
		}

		// 拉取数据
		log.Println("\r\n\r\n", _ticker, "获取数据中, ...")
		time.Sleep(time.Second * 2)
		b := _get_history(stocks[_ticker])
		if !b {
			log.Println(_ticker, "获取数据失败")
			continue
		}
		log.Println(_ticker, "获取数据成功")

		// 写入文件
		log.Println(_ticker, "写入文件中, ...")
		b = _write_klines_to_csv(stocks[_ticker])
		if !b {
			log.Println(_ticker, "写入文件失败")
			continue
		}
		log.Println(_ticker, "写入文件成功")
	}

	// 获取其他信息, 组成复合数据
	sz := _get_all(_data_url_11)
	if sz == nil {
		log.Println("拉取全深数据失败")
	}
	sh := _get_all(_data_url_12)
	if sh == nil {
		log.Println("拉取全沪数据失败")
	}

	// 补全
	for k, stock := range stocks {
		var _mm_x map[string]interface{}
		if stock.ExFlag == "0" {
			if sz != nil {
				if v, ok := sz[k]; ok {
					_mm_x = v
				}
			}
		}
		if stock.ExFlag == "1" {
			if sh != nil {
				if v, ok := sh[k]; ok {
					_mm_x = v
				}
			}
		}

		if _mm_x != nil {
			_1 := _mm_x["f9"].(float64)
			_2 := _mm_x["f23"].(float64)
			_3 := _mm_x["f2"].(float64)
			_4 := _mm_x["f8"].(float64)
			_5 := _mm_x["f10"].(float64)
			_6 := _mm_x["f115"].(float64)
			stock.Shiying = _1
			stock.GunDongShiying = _6
			stock.Shijing = _2
			stock.Price = _3
			stock.Huanshou = _4
			stock.Liangbi = _5
		}
	}

	log.Println("\r\n\r\n\r\n写入结束文件中, ...")
	b := _write_info_to_csv(stocks)
	if !b {
		log.Println("写入结束文件失败")
		return
	}
	log.Println("写入结束文件成功")
}
