package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Data struct {
	Open         float64
	PrevousClose float64
}

const (
	FILE_NAME = "2009.txt"
)

func main() {
	fs, _ := os.Open("./ETF_50/" + FILE_NAME)
	bytes, _ := ioutil.ReadAll(fs)
	log.Printf("%s", bytes)
	fs.Close()

	lines := strings.Split(string(bytes), "\n")

	var datas []*Data = make([]*Data, 0)
	var oldClose float64
	for _, line := range lines {
		if line == "" {
			continue
		}

		strs := strings.Split(line, ",")
		open, _ := strconv.ParseFloat(strs[1], 64)
		close, _ := strconv.ParseFloat(strs[2], 64)
		data := &Data{
			Open:         open,
			PrevousClose: oldClose,
		}
		datas = append(datas, data)
		oldClose = close
	}

	// 计算盈亏
	var total float64 = 0
	for _, data := range datas {
		open := data.Open
		provClose := data.PrevousClose

		if open == 0 || provClose == 0 {
			continue
		}

		v := (open - provClose) / provClose
		// log.Println(v)
		total += v
	}

	log.Println(total * 100)
}
