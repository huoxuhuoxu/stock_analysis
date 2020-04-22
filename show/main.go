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
	flag.IntVar(&mode, "mode", 1, "level相关优先级及输出")
	flag.Parse()
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// 说明输出
	if mode == 99 {
		for _, variety := range varietys {
			if variety.IsShow {
				arr := strings.Split(strings.TrimSpace(variety.Describe), "\n")
				for i, v := range arr {
					arr[i] = strings.TrimSpace(v)
				}

				fmt.Println(variety.Name)
				fmt.Println(strings.Join(arr, "\n") + "\n")
			}
		}
		return
	}

	// 价格监控
	for k, variety := range varietys {
		if variety.IsShow && variety.Level <= mode {
			go getData(k, ORIGIN_URL+variety.OriginDataUrl)
		}
	}
	go show()

	<-interrupt
}

func show() {
	var keys []string

	// 按优先级分组
	var arr = make([][]string, mode, mode)
	for key, v := range varietys {
		if v.IsShow && v.Level <= mode {
			i := v.Level - 1
			arr[i] = append(arr[i], key)
		}
	}

	// 分组排序后汇总为一个按照优先级排序的有序key组
	for _, group := range arr {
		if len(group) > 1 {
			var ks sort.StringSlice
			for _, k := range group {
				ks = append(ks, k)
			}
			sort.Sort(ks)
			group = ks
		}
		keys = append(keys, group...)
	}

	// 1s刷新输出
	t := time.Tick(time.Second * 1)
	for {
		select {
		case <-t:
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()

			// 代码, 现价, 现期基差, 短期目标, 趋势判断
			for _, k := range keys {
				if v, ok := varietys[k]; ok {
					pp := fmt.Sprintf("%%.%df", v.PricePrecision)
					fmt.Printf("%s "+pp+" "+pp+" %s\n", k, v.Price, v.Value, v.Aims)
				}
				fmt.Println("------------")
			}
		case contract := <-contractChan:
			if v, ok := varietys[contract.code]; ok {
				v.Price = contract.price
				v.Value = contract.price - contract.close
			}
		}
	}
}

// 3s更新一次数据
func getData(k, dataUrl string) {
	t := time.Tick(time.Second * 3)
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
