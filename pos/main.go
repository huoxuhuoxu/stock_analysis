package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type P struct {
	name       string        // 股票名称
	initPrice  float64       // 建仓成本
	amount     float64       // 持仓
	realPrice  float64       // 实时价格
	percentage float64       // 涨跌幅百分比
	rwMutex    *sync.RWMutex // 锁, 无写冲突, 其实可以不用锁
	url        string        // 拉取实时
}

func (self *P) positionCalculation(ctx context.Context) {
	go func() {
		chanTick := time.Tick(time.Second * 5)
		for {
			select {
			case <-ctx.Done():
				log.Printf("%s end ...", self.name)
				return
			case <-chanTick:
				resp, err := http.Get(self.url)
				if err != nil {
					log.Println(self.name, "http.Get", err)
				}
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println(self.name, "ioutil.ReadAll", err)
				}
				defer resp.Body.Close()

				str := string(body)
				arr := strings.SplitN(str, ",", 2)
				arr = strings.Split(arr[1], "\"")
				data := strings.Split(arr[0], ",")

				yesterdayPrice, _ := strconv.ParseFloat(data[1], 64)
				realPrice, _ := strconv.ParseFloat(data[2], 64)

				percentage := (realPrice - yesterdayPrice) / yesterdayPrice * 100

				// self.rwMutex.Lock()
				self.percentage = percentage
				self.realPrice = realPrice
				// self.rwMutex.Unlock()
			}
		}
	}()
}

var (
	ps []*P
)

func init() {
	log.SetFlags(0)

	ps = []*P{
		&P{"美的集团", 47.7, 1300, 0, 0, &sync.RWMutex{}, "https://hq.sinajs.cn/?_=0.8803355743806824&list=sz000333"},
		&P{"双汇发展", 25.574, 1700, 0, 0, &sync.RWMutex{}, "https://hq.sinajs.cn/?_=0.5444750447346742&list=sz000895"},
		&P{"伊利股份", 27.713, 400, 0, 0, &sync.RWMutex{}, "https://hq.sinajs.cn/?_=0.8803355743806824&list=sh600887"},
		&P{"新城控股", 42.521, 500, 0, 0, &sync.RWMutex{}, "https://hq.sinajs.cn/?_=0.8803355743806824&list=sh601155"},
		&P{"工商银行", 5.665, 8000, 0, 0, &sync.RWMutex{}, "https://hq.sinajs.cn/?_=0.8803355743806824&list=sh601398"},
	}
}

func main() {
	interrput := make(chan os.Signal, 0)
	signal.Notify(interrput, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	// 这里需要用[]*P而不是[]P, p如果是值就是值拷贝, 是一个有独立存储空间的临时变量, 对它的修改不改变[]内的p
	for _, p := range ps {
		p.positionCalculation(ctx)
	}

	show(ctx)

	for {
		<-interrput
		break
	}
	cancel()
	log.Println("exit ...")
}

func show(ctx context.Context) {
	go func() {
		chanTick := time.Tick(time.Second)
		for {
			select {
			case <-ctx.Done():
				log.Println("show end ...")
				return
			case <-chanTick:
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()

				fmt.Printf("\r\n\r\n\r\n\r\n\r\n  组合收益情况\r\n\r\n\r\n\r\n")
				var t float64
				for _, p := range ps {
					if p.realPrice == 0 {
						fmt.Printf("    %s 暂未拉取数据 \r\n\r\n", p.name)
						continue
					}
					tmpV := p.realPrice*p.amount - p.initPrice*p.amount
					fmt.Printf("    %s %.3f %.2f%% %.3f %.3f %.2f%% \r\n\r\n", p.name, p.realPrice, p.percentage, p.initPrice, tmpV, (p.realPrice-p.initPrice)/p.initPrice*100)
					t += tmpV
				}

				fmt.Printf("\r\n\r\n  持仓盈亏 %.3f \r\n\r\n\r\n", t)
			}
		}
	}()
}
