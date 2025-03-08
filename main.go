package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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

type WrapResultAPI struct {
	Date       string   `json:"date"`
	ListNumber []string `json:"listNumber"`
}

type formatForHuman struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

func main() {
	//var listTime = calculateTime(2)
	//for i := 0; i < len(listTime); i++ {
	//	time.Sleep(1 * time.Second)
	//	responseHtml := fetchData(listTime[i])
	//	//readDataFromResponse(responseHtml)
	//	storeDataResponse(responseHtml)
	//}
	// Đăng ký handler cho route /sample
	http.HandleFunc("/sample/one", getTopOneNumberBest)
	http.HandleFunc("/sample/pair", getTopPairNumberBest)

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

// Handler cho API GET /sample
func getTopOneNumberBest(w http.ResponseWriter, r *http.Request) {
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
	result := countOccurrences(listNumberString)

	//// In ra kết quả
	//fmt.Println("Occurrences of elements:")
	//for key, value := range result {
	//	if value > 10 {
	//		fmt.Printf("%s: %d\n", key, value)
	//	}
	//}
	//responseForClient := sortMapByValueDesc(result)

	// Handle and response json
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Fatal("Failed to encode response")
	}
	w.Write(jsonData)
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
	var allList [][]string
	for i := 0; i < len(resultBody); i++ {
		allList = append(allList, resultBody[i].ListNumber)
	}
	result := countPairOccurrences(allList)

	fmt.Println("Các cặp chuỗi và số lần xuất hiện:")
	for pair, count := range result {
		if count > 0 {
			fmt.Printf("Cặp (%s, %s): %d lần\n", pair[0], pair[1], count)
		}
	}

	// Handle and response json
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(result)
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

// Hàm đếm cặp chuỗi xuất hiện nhiều lần nhất và in kết quả
func countPairOccurrences(lists [][]string) map[[2]string]int {
	pairCount := make(map[[2]string]int) // Map lưu trữ số lần xuất hiện của cặp chuỗi

	// Duyệt qua tất cả các danh sách con
	for _, list := range lists {
		// Tìm các cặp chuỗi trong danh sách con
		pairs := findPairs(list)
		// Đếm tần suất các cặp chuỗi
		for _, pair := range pairs {
			pairCount[pair]++
		}
	}

	return pairCount
}

func openConnection() *sql.DB {
	// Format: "username:password@tcp(host:port)/dbname"
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.1.72:3306)/go-lang-xsmb?parseTime=true")
	if err != nil {
		log.Fatal("Không thể kết nối database:", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Không thể đóng kết nối database:", err)
		}
	}(db) // Đóng kết nối khi xong

	err = db.Ping()
	if err != nil {
		log.Fatal("Không thể ping database:", err)
	}
	return db
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
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows) // Đóng rows khi xong

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
func sortMapByValueDesc(m map[string]int) []formatForHuman {
	// Chuyển map thành slice các cặp key-value
	var sortedList []formatForHuman
	for k, v := range m {
		var resultForHuman = formatForHuman{Key: k, Value: v}
		sortedList = append(sortedList, resultForHuman)
	}

	// Sắp xếp slice theo giá trị giảm dần
	sort.Slice(sortedList, func(i, j int) bool {
		return sortedList[i].Value > sortedList[j].Value
	})

	return sortedList
}
