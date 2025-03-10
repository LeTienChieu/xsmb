package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

// DayOfXsmb model
type DayOfXsmb struct {
	ID              uint   `gorm:"primary_key"`
	DayPrize        string `gorm:"type:varchar(12)"`
	CreatedAt       time.Time
	SpecialPrize    string            `gorm:"type:varchar(5)"`
	Top1Prize       string            `gorm:"type:varchar(5)"`
	DayOfXsmbDetail []DayOfXsmbDetail `gorm:"foreignkey:DayOfXsmbID"`
}

// DayOfXsmbDetail model
type DayOfXsmbDetail struct {
	ID              uint   `gorm:"primary_key"`
	TypePrizeDetail string `gorm:"type:varchar(20)"`
	Content         string `gorm:"type:varchar(5)"`
	DayOfXsmbID     uint   // Foreign key to User
}

type WrapSQLResult struct {
	date       time.Time
	listNumber string
}

type MapPair struct {
	Pair  [2]string `json:"pair"`
	Count int       `json:"count"`
}

type WrapResultAPI struct {
	Date       string   `json:"date"`
	ListNumber []string `json:"listNumber"`
}

type ResFormatForHuman struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

type WrapFormatForHuman struct {
	StartWith  string              `json:"startWith"`
	Value      int                 `json:"value"`
	NumberList []ResFormatForHuman `json:"numberList"`
}

func main() {
	//var listTime = calculateTime(0)
	//fmt.Printf("listTime: %v\n", listTime)
	//for i := 0; i < len(listTime); i++ {
	//	time.Sleep(1 * time.Second)
	//	responseHtml := fetchData(listTime[i])
	//	//readDataFromResponse(responseHtml)
	//	storeDataResponse(responseHtml)
	//}
	// Đăng ký handler cho route /sample
	http.HandleFunc("/sample/count/start-with-detail", getTopStartNumberBestDetail)
	http.HandleFunc("/sample/count/start-with-2025", getTopStartNumberBest2025)
	http.HandleFunc("/sample/count/start-with-week", getTopStartNumberBestWeek)
	http.HandleFunc("/sample/count/start-with-month", getTopStartNumberBestMonth)
	http.HandleFunc("/sample/count/one-year", getTopNumberBestYear)
	http.HandleFunc("/sample/count/one-month", getTopNumberBestMonth)
	http.HandleFunc("/sample/count/one-week", getTopNumberBestWeek)
	http.HandleFunc("/sample/count/pair", getTopPairNumberBest)

	// Khởi động server và lắng nghe trên cổng 8080
	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func countOccurrences(list []string) map[string]int {
	countMap := make(map[string]int)

	for _, item := range list {
		countMap[item]++
	}
	return countMap
}

func getTopNumberBestYear(w http.ResponseWriter, r *http.Request) {
	// Force input is GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Fatal("Method not allowed")
	}
	var now = time.Now()
	var firstDayOf2025 = time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)
	fromDate := firstDayOf2025.Format("2006-01-02")
	toDate := now.Format("2006-01-02")
	resultBody := loadDataResponse(fromDate, toDate)
	var listNumberString []string
	for i := 0; i < len(resultBody); i++ {
		listNumberString = append(listNumberString, resultBody[i].ListNumber...)
	}
	result := countOccurrences(listNumberString)
	responseForClient := sortMapByValueDesc(result)

	// Handle and response json
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(responseForClient)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Fatal("Failed to encode response")
	}
	w.Write(jsonData)

	//push data to google sheet
	var dataPush = make([][]interface{}, len(responseForClient)+1)
	dataPush = append(dataPush, []interface{}{"Số", "Đếm", fromDate, toDate})
	for i := 0; i < len(responseForClient); i++ {
		dataPush = append(dataPush, []interface{}{responseForClient[i].Key, responseForClient[i].Value})
	}
	//push data to google sheet
	pushToSpreadSheet("StartWith", "A22", dataPush)
}

func getTopNumberBestMonth(w http.ResponseWriter, r *http.Request) {
	// Force input is GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Fatal("Method not allowed")
	}
	var now = time.Now()
	var firstDayOfLastMonth = time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local)
	fromDate := firstDayOfLastMonth.Format("2006-01-02")
	toDate := now.Format("2006-01-02")
	resultBody := loadDataResponse(fromDate, toDate)
	var listNumberString []string
	for i := 0; i < len(resultBody); i++ {
		listNumberString = append(listNumberString, resultBody[i].ListNumber...)
	}
	result := countOccurrences(listNumberString)
	responseForClient := sortMapByValueDesc(result)

	// Handle and response json
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(responseForClient)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Fatal("Failed to encode response")
	}
	w.Write(jsonData)

	//push data to google sheet
	var dataPush = make([][]interface{}, len(responseForClient)+1)
	dataPush = append(dataPush, []interface{}{"Số", "Đếm", fromDate, toDate})
	for i := 0; i < len(responseForClient); i++ {
		dataPush = append(dataPush, []interface{}{responseForClient[i].Key, responseForClient[i].Value})
	}
	//push data to google sheet
	pushToSpreadSheet("StartWith", "H22", dataPush)
}

func getTopNumberBestWeek(w http.ResponseWriter, r *http.Request) {
	// Force input is GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Fatal("Method not allowed")
	}
	var now = time.Now()
	var firstDayOfThisWeek = now.AddDate(0, 0, -int(now.Weekday()))
	fromDate := firstDayOfThisWeek.Format("2006-01-02")
	toDate := now.Format("2006-01-02")
	resultBody := loadDataResponse(fromDate, toDate)
	var listNumberString []string
	for i := 0; i < len(resultBody); i++ {
		listNumberString = append(listNumberString, resultBody[i].ListNumber...)
	}
	result := countOccurrences(listNumberString)
	responseForClient := sortMapByValueDesc(result)

	// Handle and response json
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(responseForClient)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Fatal("Failed to encode response")
	}
	w.Write(jsonData)

	//push data to google sheet
	var dataPush = make([][]interface{}, len(responseForClient)+1)
	dataPush = append(dataPush, []interface{}{"Số", "Đếm", fromDate, toDate})
	for i := 0; i < len(responseForClient); i++ {
		dataPush = append(dataPush, []interface{}{responseForClient[i].Key, responseForClient[i].Value})
	}
	//push data to google sheet
	pushToSpreadSheet("StartWith", "O22", dataPush)
}

func getTopStartNumberBestDetail(w http.ResponseWriter, r *http.Request) {
	// Force input is GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Fatal("Method not allowed")
	}
	queryParams := r.URL.Query()
	resultBody := loadDataResponse(queryParams.Get("fromDate"), queryParams.Get("toDate"))
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
	resultCountOne := countOccurrences(listNumberString)
	var wrapListNumberForRes []WrapFormatForHuman
	for item := range mapStartWith {
		var listNumberForRes []ResFormatForHuman
		for oneNumCout := range resultCountOne {
			startWithStr := oneNumCout[0:1]
			if item == startWithStr {
				listNumberForRes = append(listNumberForRes, ResFormatForHuman{Key: oneNumCout, Value: resultCountOne[oneNumCout]})
			}
		}
		wrapListNumberForRes = append(wrapListNumberForRes, WrapFormatForHuman{StartWith: item, NumberList: listNumberForRes, Value: mapStartWith[item]})
	}

	// Handle and response json
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(wrapListNumberForRes)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Fatal("Failed to encode response")
	}
	w.Write(jsonData)
}

func getTopStartNumberBest2025(w http.ResponseWriter, r *http.Request) {
	// Force input is GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Fatal("Method not allowed")
	}
	var now = time.Now()
	var firstDayOf2025 = time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)
	fromDate := firstDayOf2025.Format("2006-01-02")
	toDate := now.Format("2006-01-02")
	resultBody := loadDataResponse(fromDate, toDate)
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
	// Handle and response json
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(mapStartWith)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Fatal("Failed to encode response")
	}
	w.Write(jsonData)
	resFormatForHumanList := sortMapByValueDesc(mapStartWith)
	//push data to google sheet
	var dataPush = make([][]interface{}, len(mapStartWith)+1)
	dataPush = append(dataPush, []interface{}{"Đầu số", "Đếm", fromDate, toDate})
	for i := 0; i < len(resFormatForHumanList); i++ {
		dataPush = append(dataPush, []interface{}{resFormatForHumanList[i].Key, resFormatForHumanList[i].Value})
	}
	//push data to google sheet
	pushToSpreadSheet("StartWith", "A1", dataPush)
}

func getTopStartNumberBestWeek(w http.ResponseWriter, r *http.Request) {
	// Force input is GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Fatal("Method not allowed")
	}
	var now = time.Now()
	var firstDayOfThisWeek = now.AddDate(0, 0, -int(now.Weekday()))
	fromDate := firstDayOfThisWeek.Format("2006-01-02")
	toDate := now.Format("2006-01-02")
	resultBody := loadDataResponse(fromDate, toDate)
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
	// Handle and response json
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(mapStartWith)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Fatal("Failed to encode response")
	}
	w.Write(jsonData)

	resFormatForHumanList := sortMapByValueDesc(mapStartWith)
	//push data to google sheet
	var dataPush = make([][]interface{}, len(mapStartWith)+1)
	dataPush = append(dataPush, []interface{}{"Đầu số", "Đếm", fromDate, toDate})
	for i := 0; i < len(resFormatForHumanList); i++ {
		dataPush = append(dataPush, []interface{}{resFormatForHumanList[i].Key, resFormatForHumanList[i].Value})
	}
	//push data to google sheet
	pushToSpreadSheet("StartWith", "O1", dataPush)
}

func getTopStartNumberBestMonth(w http.ResponseWriter, r *http.Request) {
	// Force input is GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Fatal("Method not allowed")
	}
	var now = time.Now()
	var lastMonth = now.AddDate(0, -1, 0)
	var firstDayOfLastMonth = time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, time.Local)
	fromDate := firstDayOfLastMonth.Format("2006-01-01")
	toDate := now.Format("2006-01-02")
	resultBody := loadDataResponse(fromDate, toDate)
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
	// Handle and response json
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(mapStartWith)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Fatal("Failed to encode response")
	}
	w.Write(jsonData)

	resFormatForHumanList := sortMapByValueDesc(mapStartWith)
	//push data to google sheet
	var dataPush = make([][]interface{}, len(mapStartWith)+1)
	dataPush = append(dataPush, []interface{}{"Đầu số", "Đếm", fromDate, toDate})
	for i := 0; i < len(resFormatForHumanList); i++ {
		dataPush = append(dataPush, []interface{}{resFormatForHumanList[i].Key, resFormatForHumanList[i].Value})
	}
	//push data to google sheet
	pushToSpreadSheet("StartWith", "H1", dataPush)
}

// Handler cho API GET /sample
func getTopPairNumberBest(w http.ResponseWriter, r *http.Request) {
	// Force input is GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Fatal("Method not allowed")
	}
	queryParams := r.URL.Query()
	resultBody := loadDataResponse(queryParams.Get("fromDate"), queryParams.Get("toDate"))
	var allPairList [][2]string
	for i := 0; i < len(resultBody); i++ {
		var listPair = findPairs(resultBody[i].ListNumber)
		allPairList = append(allPairList, listPair...)
	}
	var mapPairs = findMostFrequentPairs(allPairList)

	// Handle and response json
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(mapPairs)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Fatal("Failed to encode response")
	}
	w.Write(jsonData)
}

// Hàm để tạo cặp số từ 2 chuỗi trong danh sách
func findPairs(list []string) [][2]string {
	var pairs [][2]string
	n := len(list)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			// Tạo cặp từ hai chuỗi
			pair := [2]string{list[i], list[j]}
			// Sắp xếp cặp để đảm bảo thứ tự không thay đổi giữa (a, b) và (b, a)
			sort.Strings(pair[:])
			pairs = append(pairs, pair)
		}
	}
	return pairs
}

// Hàm tìm cặp string xuất hiện nhiều nhất và sắp xếp theo số lần xuất hiện giảm dần
func findMostFrequentPairs(allPairList [][2]string) []MapPair {
	// Bước 1: Đếm số lần xuất hiện của từng cặp string
	pairCount := make(map[[2]string]int)
	for _, pair := range allPairList {
		pairCount[pair]++
	}
	var pairCountList []MapPair
	for pair, count := range pairCount {
		pairCountList = append(pairCountList, MapPair{pair, count})
	}

	// Bước 3: Sắp xếp slice theo số lần xuất hiện giảm dần
	sort.Slice(pairCountList, func(i, j int) bool {
		return pairCountList[i].Count > pairCountList[j].Count
	})

	return pairCountList
}

//// Handler cho API GET /sample
//func getTopThreeNumberBest(w http.ResponseWriter, r *http.Request) {
//	// Force input is GET method
//	if r.Method != http.MethodGet {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		log.Fatal("Method not allowed")
//	}
//	queryParams := r.URL.Query()
//	resultBody := loadDataResponse(queryParams.Get("fromDate"), queryParams.Get("toDate"))
//	var allList [][]string
//	for i := 0; i < len(resultBody); i++ {
//		allList = append(allList, resultBody[i].ListNumber)
//	}
//	result := countTripletOccurrences(allList)
//
//	resultFroClient := sortMapByValueThreeDesc(result)
//	// Handle and response json
//	w.Header().Set("Content-Type", "application/json")
//	jsonData, err := json.Marshal(resultFroClient)
//	if err != nil {
//		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
//		log.Fatal("Failed to encode response")
//	}
//	w.Write(jsonData)
//}
//
//// Hàm tìm các bộ 3 chuỗi trong một danh sách
//func findTriplets(list []string) [][3]string {
//	var triplets [][3]string
//	n := len(list)
//	// Duyệt qua danh sách để tạo các bộ 3 chuỗi
//	for i := 0; i < n-2; i++ {
//		for j := i + 1; j < n-1; j++ {
//			for k := j + 1; k < n; k++ {
//				triplet := [3]string{list[i], list[j], list[k]}
//				triplets = append(triplets, triplet)
//			}
//		}
//	}
//	return triplets
//}
//
//// Hàm đếm bộ 3 chuỗi xuất hiện nhiều lần nhất và in kết quả
//func countTripletOccurrences(lists [][]string) map[[3]string]int {
//	tripletCount := make(map[[3]string]int) // Map lưu trữ số lần xuất hiện của bộ 3 chuỗi
//
//	// Duyệt qua tất cả các danh sách con
//	for _, list := range lists {
//		// Tìm các bộ 3 chuỗi trong danh sách con
//		triplets := findTriplets(list)
//		// Đếm tần suất các bộ 3 chuỗi
//		for _, triplet := range triplets {
//			tripletCount[triplet]++
//		}
//	}
//
//	return tripletCount
//}

func openConnection() *sql.DB {
	// Format: "username:password@tcp(host:port)/dbname"
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.1.72:3306)/go-lang-xsmb?parseTime=true")
	if err != nil {
		log.Fatal("Không thể kết nối database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Không thể ping database:", err)
	}
	return db
}

func closeConnection(db *sql.DB) {
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Không thể đóng kết nối database:", err)
		}
	}(db) // Đóng kết nối khi xong
}

func loadDataResponse(fromDate string, toDate string) []WrapResultAPI {
	//open connection
	db := openConnection()
	if toDate == "" {
		toDate = time.Now().Format("2006-01-02")

	}

	// Query dữ liệu từ table
	rows, err := db.Query("select dox.day_of_prize as `date`, GROUP_CONCAT(doxd.content) as list_number "+
		"from day_of_xsmbs dox join day_of_xsmb_details doxd  "+
		"on dox.id = doxd.day_of_xsmb_id "+
		"where dox.day_of_prize <= cast(? as date) and dox.day_of_prize >= cast(? as date) "+
		"group by (dox.day_of_prize) "+
		"order by dox.day_of_prize desc", toDate, fromDate)
	if err != nil {
		log.Fatal("Lỗi khi query:", err)
	}
	closeConnection(db)

	// Tạo slice để chứa danh sách users
	var dayOfXsmbs []WrapSQLResult

	// Duyệt qua các row trả về
	for rows.Next() {
		var dayOfXsmb WrapSQLResult
		// Scan dữ liệu từ row vào struct
		err := rows.Scan(&dayOfXsmb.date, &dayOfXsmb.listNumber)
		if err != nil {
			log.Fatal("Lỗi khi scan dữ liệu:", err)
		}
		dayOfXsmbs = append(dayOfXsmbs, dayOfXsmb)
	}

	// Kiểm tra lỗi sau khi duyệt rows
	if err = rows.Err(); err != nil {
		log.Fatal("Lỗi khi duyệt rows:", err)
	}
	//dayOfXsmb.

	if err != nil {
		log.Fatal("Lỗi khi parse ngày:", err)
	}
	var dayOfXsmbHandled []WrapResultAPI
	for _, u := range dayOfXsmbs {
		parts := strings.Split(u.listNumber, ",")
		// Result list
		result := make([]string, len(parts))

		// For each
		for i, part := range parts {
			// Trim
			part = strings.TrimSpace(part)
			// get two last of string
			if len(part) >= 2 {
				result[i] = part[len(part)-2:]
			}
		}
		dayOfXsmbHandled = append(dayOfXsmbHandled, WrapResultAPI{Date: u.date.Format("2006-01-02"), ListNumber: result})
	}
	return dayOfXsmbHandled
}

func storeDataResponse(responseStr string) {
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
		// Lấy tên giải
		if "Giải nhất" == s.Find("td").Eq(0).Text() {
			top1Prize = s.Find("td").Eq(1).Text()
		}
		if "Đặc biệt" == s.Find("td").Eq(0).Text() {
			specialPrize = s.Find("td").Eq(1).Text()
		}
	})
	var arrayPrize []DayOfXsmbDetail
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		var checkIgnore = top1Prize != "?????"
		// Lấy tên giải
		award := s.Find("td").Eq(0).Text()

		if award == "Ký hiệu" {
			return
		}
		// Lấy kết quả xổ số
		result := s.Find("td").Eq(1).Text()
		resultArray := strings.Split(result, " ")
		if len(resultArray) > 0 && checkIgnore {
			for i := 0; i < len(resultArray); i++ {
				if strings.TrimSpace(resultArray[i]) == "" {
					continue
				}
				arrayPrize = append(arrayPrize, DayOfXsmbDetail{TypePrizeDetail: award, Content: strings.TrimSpace(resultArray[i])})
			}
		}
	})
	// Cấu hình kết nối MySQL
	db, err := gorm.Open("mysql", "root:123456@tcp(192.168.1.72:3306)/go-lang-xsmb?charset=utf8&parseTime=True")
	if err != nil {
		log.Fatal("Không thể kết nối tới cơ sở dữ liệu:", err)
	}
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Error when close connection")
		}
	}(db)

	// Tự động tạo các bảng
	db.AutoMigrate(&DayOfXsmb{}, &DayOfXsmbDetail{})

	// Tạo một user mới với trường CreatedAt
	dayOfXsmb := DayOfXsmb{DayPrize: dateStr, CreatedAt: time.Now(), SpecialPrize: specialPrize, Top1Prize: top1Prize, DayOfXsmbDetail: arrayPrize}
	db.Create(&dayOfXsmb) // Tạo user trong bảng
}

func readDataFromResponse(response string) {
	// Phân tích chuỗi HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		// Lấy tên giải
		award := s.Find("td").Eq(0).Text()

		if award == "Ký hiệu" {
			return
		}

		// Lấy kết quả xổ số
		result := s.Find("td").Eq(1).Text()

		// Kiểm tra nếu có thông tin cần thiết
		if award != "" && result != "" {
			// In ra kết quả
			fmt.Printf("Giải: %s\n", award)
			fmt.Printf("Kết quả: %s\n", result)
		}
	})
}

func fetchData(bodyStr string) string {
	// URL của API cần gọi
	url := "http://xosothudo.com.vn/LotteryResult/LastResultHome"

	// Tạo yêu cầu HTTP với headers tùy chỉnh
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

	// Gửi yêu cầu và nhận phản hồi
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Đọc nội dung trả về từ API
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// In ra kết quả nhận được
	return string(body)
}

func calculateTime(dayNumber int) []string {
	var (
		// Lấy thời gian hiện tại
		currentTime = time.Now()
		_index      = 0
		arr         []string
	)
	if currentTime.Hour() < 19 {
		_index = 1
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

// Hàm để sắp xếp map theo giá trị giảm dần
func sortMapByValueDesc(m map[string]int) []ResFormatForHuman {
	// Chuyển map thành slice các cặp key-value
	var sortedList []ResFormatForHuman
	for k, v := range m {
		var resultForHuman = ResFormatForHuman{Key: k, Value: v}
		sortedList = append(sortedList, resultForHuman)
	}

	// Sắp xếp slice theo giá trị giảm dần
	sort.Slice(sortedList, func(i, j int) bool {
		return strings.Compare(sortedList[i].Key, sortedList[j].Key) > 0
	})

	return sortedList
}

// Hàm để sắp xếp map theo giá trị giảm dần
func sortMapByValuePairDesc(m map[[2]string]int) []ResFormatForHuman {
	var ss []ResFormatForHuman
	for k, v := range m {
		if v > 3 {
			ss = append(ss, ResFormatForHuman{k[0] + "," + k[1], v})
		}
	}

	// Sắp xếp slice theo giá trị int giảm dần
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	// In kết quả đã sắp xếp
	for _, kv := range ss {
		fmt.Printf("%v: %d\n", kv.Key, kv.Value)
	}
	return ss
}

// Hàm để sắp xếp map theo giá trị giảm dần
func sortMapByValueThreeDesc(m map[[3]string]int) []ResFormatForHuman {
	var ss []ResFormatForHuman
	for k, v := range m {
		if v > 3 {
			ss = append(ss, ResFormatForHuman{k[0] + "," + k[1] + "," + k[2], v})
		}
	}

	// Sắp xếp slice theo giá trị int giảm dần
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	// In kết quả đã sắp xếp
	for _, kv := range ss {
		fmt.Printf("%v: %d\n", kv.Key, kv.Value)
	}
	return ss
}

func pushToSpreadSheet(sheetName string, startPosition string, values [][]interface{}) {
	// Đường dẫn tới file credentials.json
	credFile := "google-auth/sinuous-crow-433203-m1-3e1f8830eebd.json"

	// Đọc file credentials và tạo cấu hình
	data, err := ioutil.ReadFile(credFile)
	if err != nil {
		log.Fatalf("Không thể đọc file credentials: %v", err)
	}

	config, err := google.JWTConfigFromJSON(data, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Không thể tạo config: %v", err)
	}

	// Tạo client
	ctx := context.Background()
	client := config.Client(ctx)

	// Khởi tạo service cho Google Sheets
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Không thể tạo Sheets service: %v", err)
	}

	// ID của Google Sheet (lấy từ URL của sheet)
	spreadsheetId := "1-qgSkMXo6lSq19RcxZDcAb6OB8V_NJWeseKrmTs5d0s"

	// Tạo ValueRange để ghi dữ liệu
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Vị trí ghi dữ liệu (ví dụ: sheet1, bắt đầu từ ô A1)
	rangeData := sheetName + "!" + startPosition

	// Ghi dữ liệu vào Google Sheet
	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, rangeData, valueRange).
		ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Không thể ghi dữ liệu: %v", err)
	}

	fmt.Println("Đã ghi dữ liệu thành công!")
}
