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
		&P{"丰乐种业", 0, 0, &sync.RWMutex{}, "https://hq.sinajs.cn/?_=0.8803355743806824&list=sz000713"},
		&P{"五洲交通", 0, 0, &sync.RWMutex{}, "https://hq.sinajs.cn/?_=0.8803355743806824&list=sh600368"},
		&P{"双汇发展", 0, 0, &sync.RWMutex{}, "https://hq.sinajs.cn/?_=0.8803355743806824&list=sz000895"},
		&P{"格力电器", 0, 0, &sync.RWMutex{}, "https://hq.sinajs.cn/?_=0.8803355743806824&list=sh600651"},
		&P{"石大胜华", 0, 0, &sync.RWMutex{}, "https://hq.sinajs.cn/?_=0.8803355743806824&list=sh603028"},
		&P{"安洁科技", 0, 0, &sync.RWMutex{}, "https://hq.sinajs.cn/?_=0.8803355743806824&list=sz002635"},
		&P{"来伊份", 0, 0, &sync.RWMutex{}, "https://hq.sinajs.cn/?_=0.8803355743806824&list=sh603777"},
		&P{"航天科技", 0, 0, &sync.RWMutex{}, "https://hq.sinajs.cn/?_=0.8803355743806824&list=sh600677"},
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

				fmt.Println("Real Time:\r\n")
				fmt.Printf("    %6s %10s %10s%\r\n\r\n", "名称", "价格", "涨幅")
				for _, p := range ps {
					if p.realPrice == 0 {
						fmt.Printf("    %6s 暂未拉取数据 \r\n\r\n", p.name)
						continue
					}

					sRealPrice := fmt.Sprintf("%.2f", p.realPrice)
					sPercentage := fmt.Sprintf("%.2f", p.percentage)
					fmt.Printf("    %6s %10s %10s%%\r\n\r\n", p.name, sRealPrice, sPercentage)
				}
			}
		}
	}()
}
