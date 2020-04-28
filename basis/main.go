package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
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
				if k == "" {
					continue
				}

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
				// 套差与趋势跟随
				curA := varietys[group.Combination[0]]
				curB := varietys[group.Combination[1]]
				if curA.Price == 0 || curB.Price == 0 {
					continue
				}

				// 对冲配比
				matchingA := group.Matching[0]
				matchingB := group.Matching[1]

				// 组合名称
				name := fmt.Sprintf("%s/%s", group.Combination[0], group.Combination[1])
				// 每一组的持仓比
				matching := fmt.Sprintf("%d:%d", matchingA, matchingB)
				// 每一组的平仓划点需要付出的价格, 不考虑划点了, 如果盈利值得平仓, 可以不考虑划点的折损
				// expenditure := curA.Dash*curA.DashCoefficient*float64(matchingA) + curB.Dash*curB.DashCoefficient*float64(matchingB)

				// 有效精度
				ppA := fmt.Sprintf("%%.%df", curA.PricePrecision)
				ppB := fmt.Sprintf("%%.%df", curB.PricePrecision)
				pp := fmt.Sprintf("%s/%s", ppA, ppB)

				// 基差=(A的涨幅*A的手数*A每手的数量 - B同理) / 10, 1点=10元, 规整
				basis := (curA.Value*float64(matchingA)*curA.DashCoefficient - curB.Value*float64(matchingB)*curB.DashCoefficient) / 10

				// 判断操作行为
				aI := 0
				// 按0基差, 上下一定基差数, 作为适合开仓的定义
				if math.Abs(basis) <= group.Limit {
					aI = 1
				}

				// show
				if group.IsAll {
					// 相对价格的基点距离
					curAReaPrice := curA.Price - group.ReasonablePrice[0]
					curBReaPrice := curB.Price - group.ReasonablePrice[1]

					// 相对比例情况
					reaRatio := group.ReasonablePrice[0] / group.ReasonablePrice[1]
					/*
						相对回归的基点距离与原始比例的差值, 正: 0多了, 负: 1多了
						假设:
							差值比在 20% 以内, 认为是正常
							> +20%, 认为 可以进行 反向反套
							< -20%, 认为 可以进行 反套
					*/
					priceRatio := -(reaRatio - curAReaPrice/curBReaPrice) * 100

					/*
						当前的基点相对距离出现大的单边
							+20%, 反向建仓
							-20%, 正向建仓, 不需要管limit
					*/
					if priceRatio > 20 || priceRatio < -20 {
						if priceRatio > 20 {
							aI = 2
						}
						if priceRatio < -20 {
							aI = 1
						}
					}

					fmt.Printf("%s, %.0f, "+pp+", "+pp+", %s, %s, %s, "+pp+", %.2f\n",
						name,
						basis,
						curA.Price, curB.Price,
						curA.Value, curB.Value,
						matching,
						actions[aI],
						group.MarginConsumption,
						curAReaPrice, curBReaPrice,
						priceRatio,
					)
				} else {
					fmt.Printf("%s, %.0f, "+pp+", "+pp+", %s, %s\n",
						name,
						basis,
						curA.Price, curB.Price,
						curA.Value, curB.Value,
						matching,
						group.MarginConsumption,
					)
				}

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
