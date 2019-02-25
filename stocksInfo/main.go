package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func init() {
	log.SetFlags(0)
}

// http://pdfm.eastmoney.com/EM_UBG_PDTI_Fast/api/js?rtntype=5&token=4f1862fc3b5e77c150a2b985b12db0fd&cb=jQuery18306500869817051207_1551101184728&id=0025472&type=k&authorityType=&_=1551101309291
func main() {
	fs, err := os.Open("./data/stoke_dataset.json")
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()

	stocksStr, err := ioutil.ReadAll(fs)
	if err != nil {
		log.Fatal(err)
	}

	var stocksDataSet map[string]interface{}
	err = json.Unmarshal(stocksStr, &stocksDataSet)
	if err != nil {
		log.Fatal("uncoded json", err)
	}

	log.Println("total", len(stocksDataSet))
	curIndex := 0
	singleMutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	go func() {
		chanTick := time.Tick(time.Second)
		for {
			<-chanTick
			cmd := exec.Command("clear") //Linux example, its tested
			cmd.Stdout = os.Stdout
			cmd.Run()

			singleMutex.Lock()
			progress := float64(curIndex) / float64(len(stocksDataSet)) * 100.0
			log.Printf("当前完成 %d 个, 进度 %.2f%%", curIndex, progress)
			singleMutex.Unlock()
		}
	}()

	wg.Add(len(stocksDataSet))
	for _, v := range stocksDataSet {
		tmp := v.(map[string]interface{})
		stockCode := tmp["stockCode"].(string)
		affiliatedExchange := tmp["affiliatedExchange"].(string)
		name := tmp["name"].(string)
		// log.Println(name, stockCode, affiliatedExchange)

		go func() {
			defer func() {
				singleMutex.Lock()
				curIndex++
				// progress := float64(curIndex) / float64(len(stocksDataSet)) * 100.0
				// log.Printf("当前完成 %d 个, 进度 %.2f%%", curIndex, progress)
				singleMutex.Unlock()
				wg.Done()
			}()
			var b bool
			for !b {
				b = Get(stockCode, affiliatedExchange, name)
				time.Sleep(time.Millisecond * 1000)
			}
		}()
		time.Sleep(time.Millisecond * 200)
	}

	wg.Wait()
	log.Println("succ")
}

func Get(stockCode, ae, name string) bool {
	url := "http://pdfm.eastmoney.com/EM_UBG_PDTI_Fast/api/js?rtntype=5&token=4f1862fc3b5e77c150a2b985b12db0fd&cb=jQuery18306500869817051207_1551101184728&id=" + stockCode + ae + "&type=k&authorityType=&_=1551101309291"

	resp, err := http.Get(url)
	if err != nil {
		log.Println(name, err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(name, err)
		return false
	}

	str := string(body)
	ret := strings.Split(strings.Split(str, "data\":")[1], "})")

	var infos []interface{}
	err = json.Unmarshal([]byte(ret[0]), &infos)
	if err != nil {
		log.Println("uncoded json", name, err)
		return false
	}

	fs, err := os.OpenFile("./data/details/"+stockCode, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Println(name, err)
		return false
	}
	defer fs.Close()

	for _, tmp := range infos {
		str := tmp.(string)
		fs.WriteString(str + "\r\n")
	}

	return true
}
