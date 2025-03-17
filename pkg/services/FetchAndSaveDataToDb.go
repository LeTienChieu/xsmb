package services

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"myproject/pkg/repositories"
	"net/http"
	"strings"
	"time"
)

func StoreDataFromResponse(responseStr string) {
	// Split responseStr HTML by tag
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(responseStr))
	if err != nil {
		log.Fatal(err)
	}
	var dateStr = time.Now().Format("20060102")
	var specialPrize string
	var top1Prize string
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		//get date data from HTML
		dateStr = doc.Find("#txtLotteryDate").Text()
		if "Giải nhất" == s.Find("td").Eq(0).Text() {
			top1Prize = s.Find("td").Eq(1).Text()
		}
		if "Đặc biệt" == s.Find("td").Eq(0).Text() {
			specialPrize = s.Find("td").Eq(1).Text()
		}
	})
	var checkIgnore = strings.Contains(top1Prize, "?") ||
		strings.Contains(specialPrize, "?") ||
		strings.Contains(top1Prize, " ") ||
		strings.Contains(specialPrize, " ")
	if checkIgnore {
		return
	}
	var arrayPrize []repositories.DayOfXsmbDetail
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		award := s.Find("td").Eq(0).Text()
		if award == "Ký hiệu" {
			return
		}
		result := s.Find("td").Eq(1).Text()
		resultArray := strings.Split(result, " ")
		if len(resultArray) > 0 {
			for i := 0; i < len(resultArray); i++ {
				if strings.TrimSpace(resultArray[i]) == "" {
					continue
				}
				arrayPrize = append(arrayPrize, repositories.DayOfXsmbDetail{TypePrizeDetail: award, Content: strings.TrimSpace(resultArray[i])})
			}
		}
	})
	var db = repositories.OpenOrmConnection()
	// Auto create DDL
	db.AutoMigrate(&repositories.DayOfXsmb{}, &repositories.DayOfXsmbDetail{})

	// create new dayOfXsmb
	dayOfXsmb := repositories.DayOfXsmb{DayPrize: dateStr, CreatedAt: time.Now(), SpecialPrize: specialPrize, Top1Prize: top1Prize, DayOfXsmbDetail: arrayPrize}
	db.Create(&dayOfXsmb)
	repositories.CloseOrmConnection(db)
}

func CalculateMaxTimeWithDb() []string {
	var (
		currentTime = time.Now()
		_index      = 0
		arr         []string
	)
	var maxTimeInDb = repositories.GetMaxCreateDate()
	if currentTime.Hour() < 19 {
		_index = 1
	}
	// Count number of day
	var dayNumber = int(currentTime.Sub(maxTimeInDb).Hours() / 24)
	// case ignore data
	if dayNumber == 0 && maxTimeInDb.Hour() > 12 {
		return arr
	}
	for i := _index; i <= dayNumber; i++ {
		if i == 0 {
			arr = append(arr, fmt.Sprintf("date=%s", currentTime.Format("20060102")))
			continue
		}
		var timeCalculate = currentTime.AddDate(0, 0, -i)
		arr = append(arr, fmt.Sprintf("date=%s", timeCalculate.Format("20060102")))
	}

	return arr
}

func Calculate3YearMinTimeWithDb() []string {
	var arr []string
	var minTimeInDb = repositories.GetMinCreateDate()
	maxTimeForCalculate := minTimeInDb.AddDate(0, 0, -1)
	minTimeForCalculate := maxTimeForCalculate.AddDate(-3, 0, 0)
	// Count number of day
	var dayNumber = int(maxTimeForCalculate.Sub(minTimeForCalculate).Hours() / 24)
	for i := 0; i <= dayNumber; i++ {
		if i == 0 {
			arr = append(arr, fmt.Sprintf("date=%s", maxTimeForCalculate.Format("20060102")))
			continue
		}
		var timeCalculate = maxTimeForCalculate.AddDate(0, 0, -i)
		arr = append(arr, fmt.Sprintf("date=%s", timeCalculate.Format("20060102")))
	}
	return arr
}

func FetchDataFromSource() {
	var arr = Calculate3YearMinTimeWithDb()
	for i := 0; i < len(arr); i++ {
		time.Sleep(1 * time.Second)
		var bodyStr = arr[i]
		var responseStr = FetchData(bodyStr)
		StoreDataFromResponse(responseStr)
	}
}

func FetchData(bodyStr string) string {
	url := repositories.GetConfigSourceUrl()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyStr)))
	log.Printf("Fetch data: %v", bodyStr)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "vi")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "http://xosothudo.com.vn")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	return string(body)
}
