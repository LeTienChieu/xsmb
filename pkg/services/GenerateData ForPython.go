package services

import (
	"myproject/pkg/repositories"
	"strings"
	"time"
)

type StoreDataWithDateAndNumber struct {
	Date          time.Time `gorm:"primary_key"`
	ListNumberStr string    `gorm:"type:varchar(255)"`
}

func GenerateDataForOneNumber() {
	var result []StoreDataWithDateAndNumber
	resultSql := repositories.LoadAllData()
	for _, item := range resultSql {
		parts := strings.Split(item.ListNumberStr, ",")
		var output []string
		for _, s := range parts {
			if len(s) >= 2 {
				output = append(output, s[len(s)-2:])
			}
		}
		resultStr := strings.Join(output[:], ",")
		elementResult := StoreDataWithDateAndNumber{item.Date, resultStr}
		result = append(result, elementResult)
	}
	db := repositories.OpenOrmConnection()
	for _, item := range result {
		db.Create(&item)
	}
	repositories.CloseOrmConnection(db)
}
