package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
)

var (
	contractChan = make(chan *Contract, 1)
	mode         int
)

func init() {
	flag.IntVar(&mode, "mode", 0, "0: all, 1: 主货, 2: 主货+权")
	flag.Parse()
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	for k, variety := range varietys {
		if variety.IsShow {
			go getData(k, ORIGIN_URL+variety.OriginDataUrl)
		}
	}
	go show()

	<-interrupt
}

func show() {
	t := time.Tick(time.Second * 5)
	var keys sort.StringSlice

	for k, v := range varietys {
		if v.IsShow {
			keys = append(keys, k)
		}
	}
	sort.Sort(keys)
	switch mode {
	case 0:
		keys = sort.StringSlice{"a2009", "m2101", "c2009", "rb2101", "CF009"}
	case 1:
		keys = sort.StringSlice{"a2009", "m2101", "c2009"}
	case 2:
		keys = sort.StringSlice{"a2009", "m2101", "c2009", "rb2101"}
	}

	for {
		select {
		case <-t:
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()

			// 代码, 现价, 现期基差, 短期目标, 趋势判断
			// fmt.Printf("name, price, spot, tmp-aims, trend \n")
			for _, k := range keys {
				if v, ok := varietys[k]; ok {
					basis := 0.0 // 基差
					if v.SpotPrice != 0 {
						basis = v.SpotPrice - v.Price
					}
					fmt.Printf("%s %.0f %.0f %.0f %.0f %s\n", v.Code, v.Price, basis, v.Value, v.Aims, v.Trend)
				}
			}
		case contract := <-contractChan:
			if v, ok := varietys[contract.code]; ok {
				v.Price = contract.price
				v.VolatilityValue = contract.high - contract.low
				v.Value = contract.price - contract.close
			}
		}
	}
}

func getData(k, dataUrl string) {
	t := time.Tick(time.Second * 5)
	client := &http.Client{}
	for {
		select {
		case <-t:
			req, _ := http.NewRequest("GET", dataUrl, nil)
			for hk, hv := range headers {
				req.Header.Set(hk, hv)
			}

			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				break
			}

			respData, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				break
			}
			resp.Body.Close()

			respStr := string(respData)
			respStr = strings.Replace(respStr, "aa(", "", 1)
			respStr = strings.Replace(respStr, "})", "}", 1)

			var data = make(map[string]interface{})
			json.Unmarshal([]byte(respStr), &data)

			if _, ok := data["qt"]; ok {
				qt := data["qt"].(map[string]interface{})
				contract := &Contract{
					code:  k,
					open:  qt["o"].(float64),
					close: qt["qrspj"].(float64),
					high:  qt["h"].(float64),
					low:   qt["l"].(float64),
					price: qt["p"].(float64),
				}

				contractChan <- contract
			}
		}
	}
}
