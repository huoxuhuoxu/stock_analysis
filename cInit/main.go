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

	// PE(5, datas)

	// PB(1, datas)

	// PbMultiplicationPe(10, datas)

	HighChange(15, datas)
}

// PE
func PE(peV float64, datas []string) {
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
			if tmpV < peV && tmpV > 0 {
				fmt.Printf("|"+strings.Repeat("%-14s", 7)+"\r\n",
					arr[2], arr[1], arr[3], arr[5], arr[7], arr[16], arr[17])
				count++
			}
		}
	}

	fmt.Println(strings.Repeat("-", 102))
	fmt.Printf("|%-92s|\r\n", fmt.Sprintf("符合条件可选股票数目: %d", count))
	fmt.Println(strings.Repeat("-", 102))
}

// PB
func PB(pbV float64, datas []string) {
	fmt.Println(strings.Repeat("-", 102))
	fmt.Printf(strings.Repeat("|%-11s", 7)+"|\r\n", "名称", "代码", "价格", "涨幅", "成交额", "市盈率", "市净率")
	fmt.Println(strings.Repeat("-", 102))

	count := 0
	for _, data := range datas {
		arr := strings.Split(data, ",")
		if arr[4] != "-" && !strings.HasPrefix(arr[2], "*ST") && !strings.HasPrefix(arr[2], "ST") {
			tmpStrFloat := arr[17]
			tmpV, err := strconv.ParseFloat(tmpStrFloat, 64)
			if err != nil {
				continue
			}
			if tmpV < pbV && tmpV > 0 {
				fmt.Printf("|"+strings.Repeat("%-14s", 7)+"\r\n",
					arr[2], arr[1], arr[3], arr[5], arr[7], arr[16], arr[17])
				count++
			}
		}
	}

	fmt.Println(strings.Repeat("-", 102))
	fmt.Printf("|%-92s|\r\n", fmt.Sprintf("符合条件可选股票数目: %d", count))
	fmt.Println(strings.Repeat("-", 102))
}

// pb * pe
func PbMultiplicationPe(mV float64, datas []string) {
	fmt.Println(strings.Repeat("-", 102))
	fmt.Printf(strings.Repeat("|%-11s", 7)+"|\r\n", "名称", "代码", "价格", "涨幅", "成交额", "市盈率", "市净率")
	fmt.Println(strings.Repeat("-", 102))

	count := 0
	for _, data := range datas {
		arr := strings.Split(data, ",")
		if arr[4] != "-" && !strings.HasPrefix(arr[2], "*ST") && !strings.HasPrefix(arr[2], "ST") {
			tmpPeFloat := arr[16]
			tmpPbFloat := arr[17]
			tmpChangeFloat := arr[15]
			tmpPe, err := strconv.ParseFloat(tmpPeFloat, 64)
			if err != nil {
				continue
			}
			tmpPb, err := strconv.ParseFloat(tmpPbFloat, 64)
			if err != nil {
				continue
			}
			tmpChange, err := strconv.ParseFloat(tmpChangeFloat, 64)
			if err != nil {
				continue
			}
			tmpV := tmpPb * tmpPe
			if tmpV < mV && tmpV > 0 && tmpChange > 2 {
				fmt.Printf("|"+strings.Repeat("%-14s", 7)+"\r\n",
					arr[2], arr[1], arr[3], arr[5], arr[7], arr[16], arr[17])
				count++
			}
		}
	}

	fmt.Println(strings.Repeat("-", 102))
	fmt.Printf("|%-92s|\r\n", fmt.Sprintf("符合条件可选股票数目: %d", count))
	fmt.Println(strings.Repeat("-", 102))
}

// 高换手
func HighChange(hcV float64, datas []string) {
	fmt.Println(strings.Repeat("-", 102))
	fmt.Printf(strings.Repeat("|%-11s", 7)+"|\r\n", "名称", "代码", "价格", "涨幅", "成交额", "市盈率", "市净率")
	fmt.Println(strings.Repeat("-", 102))

	count := 0
	for _, data := range datas {
		arr := strings.Split(data, ",")
		if arr[4] != "-" && !strings.HasPrefix(arr[2], "*ST") && !strings.HasPrefix(arr[2], "ST") {
			tmpH1Float := arr[15]
			tmpH2Float := arr[8]
			tmpV1, err := strconv.ParseFloat(tmpH1Float, 64)
			if err != nil {
				continue
			}
			tmpV2, err := strconv.ParseFloat(tmpH2Float, 64)
			if err != nil {
				continue
			}
			if tmpV1 > hcV && tmpV2 > 10 {
				fmt.Printf("|"+strings.Repeat("%-14s", 7)+"\r\n",
					arr[2], arr[1], arr[3], arr[5], arr[7], arr[16], arr[17])
				count++
			}
		}
	}

	fmt.Println(strings.Repeat("-", 102))
	fmt.Printf("|%-92s|\r\n", fmt.Sprintf("符合条件可选股票数目: %d", count))
	fmt.Println(strings.Repeat("-", 102))
}
