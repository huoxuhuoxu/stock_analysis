package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Stoke struct {
	Name               string `json:"name"`
	StockCode          string `json:"stockCode"`
	AffiliatedExchange string `json:"affiliatedExchange"`
}

func init() {
	log.SetFlags(0)
}

func main() {
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

	stokeDataSet := make(map[string]Stoke)
	for _, data := range datas {
		arr := strings.Split(data, ",")
		stoke := Stoke{arr[2], arr[1], arr[0]}
		stokeDataSet[arr[1]] = stoke
	}

	// 编码字符串
	strDataSet, _ := json.MarshalIndent(stokeDataSet, "", "\t")
	fs, err := os.OpenFile("./data/stokeDataSet.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()

	fs.WriteString(string(strDataSet))
}
