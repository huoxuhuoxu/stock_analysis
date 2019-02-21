package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	Nums()

}

func Nums() {

	pageIndex := "1"
	pageCount := "10000"
	initUrl := "http://nufm.dfcfw.com/EM_Finance2014NumericApplication/JS.aspx?cb=jQuery112407378636087109598_1528454080519&type=CT&token=4f1862fc3b5e77c150a2b985b12db0fd&js=(%7Bdata%3A%5B(x)%5D%2CrecordsTotal%3A(tot)%2CrecordsFiltered%3A(tot)%7D)&cmd=C._A&sty=FCOIATC&st=(ChangePercent)&sr=-1&p=" + pageIndex + "&ps=" + pageCount + "&_=1528454080520"

	resp, err := http.Get(initUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	str := string(body)
	ret := strings.Split(strings.Split(str, "data:[\"")[1], "\"],")

	strData := ret[0]
	datas := strings.Split(strData, "\",\"")

	log.Println(strings.Repeat("-", 102))
	log.Printf(strings.Repeat("|%-11s", 7)+"|", "名称", "代码", "价格", "涨幅", "成交额", "市盈率", "市净率")
	log.Println(strings.Repeat("-", 102))

	count := 0
	for _, data := range datas {
		arr := strings.Split(data, ",")

		if arr[4] != "-" && !strings.HasPrefix(arr[2], "*ST") && !strings.HasPrefix(arr[2], "ST") {
			log.Printf("|"+strings.Repeat("%-14s", 7),
				arr[2], arr[1], arr[3], arr[5], arr[7], arr[16], arr[17])

			count++
		}
	}

	log.Println(strings.Repeat("-", 102))
	log.Printf("|%-92s|", fmt.Sprintf("今日可选股票数目: %d", count))
	log.Println(strings.Repeat("-", 102))

}

func Fill() {

}
