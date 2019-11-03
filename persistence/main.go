package main

import (
	"encoding/json"
	"log"
	"sort"
	"stock_analysis/base"
	"stock_analysis/models"
	"time"
)

const (
	filterLen = 50
)

func main() {
	// 第一分组
	gds := base.GetMoneyData()
	var list base.GeneralDatas
	for _, gd := range gds {
		gd.Score = gd.Changehands + gd.MainP
		list = append(list, gd)
	}
	sort.Sort(list)
	firstList := list[0:filterLen]

	// db
	db, err := models.NewPersistenceDb()
	if err != nil {
		log.Fatal(err)
	}

	// write
	byteRawData, _ := json.Marshal(list)
	byteFirstData, _ := json.Marshal(firstList)
	recordDate := time.Now().Format("2006-01-02")

	err = db.Append(string(byteRawData), string(byteFirstData), recordDate)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Succ !")

	// get-test
	// db, err := models.NewPersistenceDb()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// v, err := db.GetAll()
	// for _, c := range v {
	// 	log.Printf("%+v, %+v", c.RecordDate, err)
	// }
}
