package services

import (
	"fmt"
	"myproject/pkg/repositories"
	"myproject/pkg/utils"
	"time"
)

type LogicBusinessService struct {
}

func NewLogicBusinessService() *LogicBusinessService {
	return &LogicBusinessService{}
}

func GetTopNumberBestYear2024V2() {
	var firstDayOf2024 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)
	fromDate := firstDayOf2024.Format("2006-01-02")
	var lastDayOf2024 = time.Date(2024, 12, 31, 0, 0, 0, 0, time.Local)
	toDate := lastDayOf2024.Format("2006-01-02")
	resultBody := repositories.LoadDataResponse(fromDate, toDate)
	var listNumberString []string
	for _, item := range resultBody {
		listNumberString = append(listNumberString, item.ListNumber...)
	}
	result := utils.CountOccurrences(listNumberString)
	responseForClient := utils.SortMapByValueDesc(result)
	//push data to google sheet
	var dataPush = make([][]interface{}, len(responseForClient)+1)
	dataPush = append(dataPush, []interface{}{"Số", "Đếm", fromDate, toDate})
	for _, item := range responseForClient {
		dataPush = append(dataPush, []interface{}{item.Key, item.Value})
	}
	//push data to google sheet
	PushToSpreadSheet("Phan_tich_Lo", "A22", dataPush)
}

func GetTopNumberBest2025V2() {
	var now = time.Now()
	var firstDayOfLastMonth = time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)
	fromDate := firstDayOfLastMonth.Format("2006-01-02")
	toDate := now.Format("2006-01-02")
	resultBody := repositories.LoadDataResponse(fromDate, toDate)
	var listNumberString []string
	for _, item := range resultBody {
		listNumberString = append(listNumberString, item.ListNumber...)
	}
	result := utils.CountOccurrences(listNumberString)
	responseForClient := utils.SortMapByValueDesc(result)
	//push data to google sheet
	var dataPush = make([][]interface{}, len(responseForClient)+1)
	dataPush = append(dataPush, []interface{}{"Số", "Đếm", fromDate, toDate})
	for i := 0; i < len(responseForClient); i++ {
		dataPush = append(dataPush, []interface{}{responseForClient[i].Key, responseForClient[i].Value})
	}
	//push data to google sheet
	PushToSpreadSheet("Phan_tich_Lo", "H22", dataPush)
}

func GetTopNumberBestCurrentMonthV2() {
	var now = time.Now()
	var firstDatOfThisMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	fromDate := firstDatOfThisMonth.Format("2006-01-02")
	toDate := now.Format("2006-01-02")
	resultBody := repositories.LoadDataResponse(fromDate, toDate)
	var listNumberString []string
	for i := 0; i < len(resultBody); i++ {
		listNumberString = append(listNumberString, resultBody[i].ListNumber...)
	}
	result := utils.CountOccurrences(listNumberString)
	responseForClient := utils.SortMapByValueDesc(result)
	//push data to google sheet
	var dataPush = make([][]interface{}, len(responseForClient)+1)
	dataPush = append(dataPush, []interface{}{"Số", "Đếm", fromDate, toDate})
	for i := 0; i < len(responseForClient); i++ {
		dataPush = append(dataPush, []interface{}{responseForClient[i].Key, responseForClient[i].Value})
	}
	var dataMissPush = make([][]interface{}, len(responseForClient))
	dataMissPush = append(dataMissPush, []interface{}{"Số chưa về", "Đếm", fromDate, toDate})
	for i := 0; i < 100; i++ {
		iTypeStr := fmt.Sprintf("%d", i)
		if i < 10 {
			iTypeStr = fmt.Sprintf("0%d", i)
		}
		if _, ok := result[iTypeStr]; !ok {
			dataMissPush = append(dataMissPush, []interface{}{iTypeStr, 0})
		}
	}

	//push data to google sheet
	PushToSpreadSheet("Phan_tich_Lo", "O22", dataPush)

	PushToSpreadSheet("Phan_tich_Lo", "V22", dataMissPush)
}

func GetTopStartNumberBest2025V2() {
	var now = time.Now()
	var firstDayOfLastMonth = time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)
	fromDate := firstDayOfLastMonth.Format("2006-01-01")
	toDate := now.Format("2006-01-02")
	resultBody := repositories.LoadDataResponse(fromDate, toDate)
	var listNumberString []string
	for i := 0; i < len(resultBody); i++ {
		listNumberString = append(listNumberString, resultBody[i].ListNumber...)
	}
	mapStartWith := make(map[string]int)
	for i := 0; i < len(listNumberString); i++ {
		str := listNumberString[i][0:1]
		if _, ok := mapStartWith[str]; !ok {
			mapStartWith[str] = 1
		} else {
			mapStartWith[str] = mapStartWith[str] + 1
		}
	}
	resFormatForHumanList := utils.SortMapByValueDesc(mapStartWith)
	//push data to google sheet
	var dataPush = make([][]interface{}, len(mapStartWith)+1)
	dataPush = append(dataPush, []interface{}{"Đầu số", "Đếm", fromDate, toDate})
	for i := 0; i < len(resFormatForHumanList); i++ {
		dataPush = append(dataPush, []interface{}{resFormatForHumanList[i].Key, resFormatForHumanList[i].Value})
	}
	//push data to google sheet
	PushToSpreadSheet("Phan_tich_Lo", "H1", dataPush)
}

func GetTopStartNumberBestCurentMonthV2() {
	var now = time.Now()
	var firstDatOfThisMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	fromDate := firstDatOfThisMonth.Format("2006-01-02")
	toDate := now.Format("2006-01-02")
	resultBody := repositories.LoadDataResponse(fromDate, toDate)
	var listNumberString []string
	for i := 0; i < len(resultBody); i++ {
		listNumberString = append(listNumberString, resultBody[i].ListNumber...)
	}
	mapStartWith := make(map[string]int)
	for i := 0; i < len(listNumberString); i++ {
		str := listNumberString[i][0:1]
		if _, ok := mapStartWith[str]; !ok {
			mapStartWith[str] = 1
		} else {
			mapStartWith[str] = mapStartWith[str] + 1
		}
	}
	resFormatForHumanList := utils.SortMapByValueDesc(mapStartWith)
	//push data to google sheet
	var dataPush = make([][]interface{}, len(mapStartWith)+1)
	dataPush = append(dataPush, []interface{}{"Đầu số", "Đếm", fromDate, toDate})
	for i := 0; i < len(resFormatForHumanList); i++ {
		dataPush = append(dataPush, []interface{}{resFormatForHumanList[i].Key, resFormatForHumanList[i].Value})
	}
	//push data to google sheet
	PushToSpreadSheet("Phan_tich_Lo", "O1", dataPush)
}

func GetTopStartNumberBest2024V2() {
	var firstDayOf2024 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)
	fromDate := firstDayOf2024.Format("2006-01-02")
	var lastDayOf2024 = time.Date(2024, 12, 31, 0, 0, 0, 0, time.Local)
	toDate := lastDayOf2024.Format("2006-01-02")
	resultBody := repositories.LoadDataResponse(fromDate, toDate)
	var listNumberString []string
	for i := 0; i < len(resultBody); i++ {
		listNumberString = append(listNumberString, resultBody[i].ListNumber...)
	}
	mapStartWith := make(map[string]int)
	for i := 0; i < len(listNumberString); i++ {
		str := listNumberString[i][0:1]
		if _, ok := mapStartWith[str]; !ok {
			mapStartWith[str] = 1
		} else {
			mapStartWith[str] = mapStartWith[str] + 1
		}
	}
	resFormatForHumanList := utils.SortMapByValueDesc(mapStartWith)
	//push data to google sheet
	var dataPush = make([][]interface{}, len(mapStartWith)+1)
	dataPush = append(dataPush, []interface{}{"Đầu số", "Đếm", fromDate, toDate})
	for i := 0; i < len(resFormatForHumanList); i++ {
		dataPush = append(dataPush, []interface{}{resFormatForHumanList[i].Key, resFormatForHumanList[i].Value})
	}
	//push data to google sheet
	PushToSpreadSheet("Phan_tich_Lo", "A1", dataPush)
}
