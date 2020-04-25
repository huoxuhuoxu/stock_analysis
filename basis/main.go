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

	// 特殊优先级: 0, 输出
	if mode == 0 {
		fmt.Printf("暂无")
		return
	}

	// 价格监控
	for _, group := range groups {
		if group.Level <= mode {
			for _, k := range group.Combination {
				v := varietys[k]
				go getData(k, ORIGIN_URL+v.OriginDataUrl)
			}
		}
	}
	go show()

	<-interrupt
}

// 输出显示 ----- 需要增加独立输出 原油品种 的功能
func show() {
	// 1s刷新输出
	t := time.Tick(time.Second * 1)
	for {
		select {
		case <-t:
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()

			// 输出
			for _, group := range groups {
				curA := varietys[group.Combination[0]]
				curB := varietys[group.Combination[1]]
				if curA.Price == 0 || curB.Price == 0 {
					continue
				}

				// 组合名称
				name := fmt.Sprintf("%s/%s", group.Combination[0], group.Combination[1])
				// 每一组的持仓比
				matching := fmt.Sprintf("%d:%d", curA.Amount, curB.Amount)
				// 每一组的平仓划点需要付出的价格
				expenditure := curA.Dash*curA.DashCoefficient*float64(curA.Amount) + curB.Dash*curB.DashCoefficient*float64(curB.Amount)

				// 有效精度
				ppA := fmt.Sprintf("%%.%df", curA.PricePrecision)
				ppB := fmt.Sprintf("%%.%df", curB.PricePrecision)
				pp := fmt.Sprintf("%s/%s", ppA, ppB)

				// 基差=(A的涨幅*A的手数*A每手的数量 - B同理) / 10, 1点=10元, 规整
				basis := (curA.Value*float64(curA.Amount)*curA.DashCoefficient - curB.Value*float64(curB.Amount)*curB.DashCoefficient) / 10

				// 判断操作行为
				aI := 0
				if basis <= group.Limit {
					aI = 1
				}
				if basis >= group.Profit {
					aI = 2
				}

				// show
				fmt.Printf("%s, %.0f, "+pp+", "+pp+", %s, %.0f, %s, %s\n",
					name,
					basis,
					curA.Price, curB.Price,
					curA.Value, curB.Value,
					matching, expenditure, actions[aI], group.MarginConsumption)

				fmt.Println("----")
			}

		case contract := <-contractChan:
			if v, ok := varietys[contract.code]; ok {
				v.Price = contract.price
				v.Value = contract.price - contract.close
			}
		}
	}
}

// 更新数据
func getData(k, dataUrl string) {
	t := time.Tick(time.Second * 2)
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
