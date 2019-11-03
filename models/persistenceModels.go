package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// --------------------------------------- 数据记录
type StockData struct {
	gorm.Model
	RawData    string `gorm:"type:text;not null"`             // 聚合后的原始数据
	FirstData  string `gorm:"type:text;not null"`             // 第一分组数据
	RecordDate string `gorm:"type:varchar(128);unique_index"` // 记录日期, 格式(2019-06-10)
}

type PersistenceDb struct {
	db *gorm.DB
}

func NewPersistenceDb() (*PersistenceDb, error) {
	self := &PersistenceDb{}

	db, err := gorm.Open("sqlite3", "./dbs/persistence.db")
	if err != nil {
		return nil, err
	}
	self.db = db

	self.db.AutoMigrate(&StockData{})
	return self, nil
}

func (self *PersistenceDb) Append(rawData, firstData, recordDate string) error {
	var stockData StockData
	// 判断是否存在
	self.db.Model(&StockData{}).Where(&StockData{
		RecordDate: recordDate,
	}).First(&stockData)

	// add
	if stockData.ID == 0 {
		stockData = StockData{
			RawData:    rawData,
			FirstData:  firstData,
			RecordDate: recordDate,
		}
		ret := self.db.Create(&stockData)
		if ret.RowsAffected != 1 {
			return errors.New("Add failed")
		}

		return nil
	}

	// update
	ret := self.db.Model(&StockData{}).Where(&StockData{
		RecordDate: recordDate,
	}).Updates(map[string]interface{}{
		"RawData":   rawData,
		"FirstData": firstData,
	})
	if ret.RowsAffected != 1 {
		return errors.New("Update failed")
	}

	return nil
}

func (self *PersistenceDb) GetByRecordDate(recordDate string) (*StockData, error) {
	var stockData StockData
	ret := self.db.Where(&StockData{RecordDate: recordDate}).First(&stockData)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &stockData, nil
}

func (self *PersistenceDb) GetAll() ([]*StockData, error) {
	var list []*StockData
	ret := self.db.Where(&StockData{}).Find(&list)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return list, nil
}
