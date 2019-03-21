package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
* @doc
*
*
*
* @字段分析
*	2			1=上证/2=深成		0
*	002547		股票代码			1
*	春兴精工	 股票名称		     2
*	7.77		当前价格			3
*	0.71		上涨价格			4
*	10.06		涨幅				5
*	2062612		206万手成交量		6
*	1549976928	成交额				7
*	13.31		振幅				8
*	7.77		盘中最高			9
*	6.83		盘中最低			10
*	6.93		开盘价				11
*	7.06		昨收盘价			12
*	0.00		5分钟涨速?			13
*	1.96		量比				14
*	25.86		换手率				15
*	148.87		动态市盈率			16
*	3.19		市净率				17
*	8765004174	总市值
*	6197961277	流通市值
*	81.12%		60日涨幅
*	36.8%		年初至今涨幅
*	0.00		涨速?
*	上市日期
*	最新有效日期
*
 */

func init() {
	log.SetFlags(0)
}

func main() {

	// pageIndex := "1"
	// pageCount := "10000"
	// initUrl := "http://nufm.dfcfw.com/EM_Finance2014NumericApplication/JS.aspx?cb=jQuery112407378636087109598_1528454080519&type=CT&token=4f1862fc3b5e77c150a2b985b12db0fd&js=(%7Bdata%3A%5B(x)%5D%2CrecordsTotal%3A(tot)%2CrecordsFiltered%3A(tot)%7D)&cmd=C._A&sty=FCOIATC&st=(ChangePercent)&sr=-1&p=" + pageIndex + "&ps=" + pageCount + "&_=1528454080520"
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

	fmt.Println(strings.Repeat("-", 102))
	fmt.Printf(strings.Repeat("|%-11s", 7)+"|\r\n", "名称", "代码", "价格", "涨幅", "成交额", "市盈率", "市净率")
	fmt.Println(strings.Repeat("-", 102))

	count := 0
	for _, data := range datas {
		arr := strings.Split(data, ",")
		if arr[4] != "-" && !strings.HasPrefix(arr[2], "*ST") && !strings.HasPrefix(arr[2], "ST") {
			tmpStrFloat := arr[16]
			tmpV, err := strconv.ParseFloat(tmpStrFloat, 64)
			if err != nil {
				continue
			}
			if tmpV < 5 {
				fmt.Printf("|"+strings.Repeat("%-14s", 7)+"\r\n",
					arr[2], arr[1], arr[3], arr[5], arr[7], arr[16], arr[17])

				count++
			}
		}
	}

	fmt.Println(strings.Repeat("-", 102))
	fmt.Printf("|%-92s|\r\n", fmt.Sprintf("今日可选股票数目: %d", count))
	fmt.Println(strings.Repeat("-", 102))

}
