package repositories

import (
	"log"
	"strings"
	"time"
)

type WrapResultAPI struct {
	Date       string   `json:"date"`
	ListNumber []string `json:"listNumber"`
}

type WrapSQLResult struct {
	date       time.Time
	listNumber string
}

func GetMaxCreateDate() time.Time {
	// Query max date from day_of_xsmbs table
	var db = OpenNativeConnection()
	rows, err := db.Query("SELECT MAX(dox.created_at) as `date` from day_of_xsmbs dox")
	if err != nil {
		log.Fatal(err)
	}
	//close connection
	CloseNativeConnection(db)
	var maxTimeInDb time.Time
	for rows.Next() {
		// Scan dữ liệu
		err := rows.Scan(&maxTimeInDb)
		if err != nil {
			log.Print("Lỗi khi scan dữ liệu:", err)
		}
	}
	// Kiểm tra lỗi sau khi duyệt rows
	if rows.Err() != nil {
		log.Fatal(rows.Err())
	}
	return maxTimeInDb
}

func GetMinCreateDate() time.Time {
	// Query max date from day_of_xsmbs table
	var db = OpenNativeConnection()
	rows, err := db.Query("SELECT MIN(dox.day_of_prize) as `date` from day_of_xsmbs dox")
	if err != nil {
		log.Fatal(err)
	}
	//close connection
	CloseNativeConnection(db)
	var minTimeInDb time.Time
	for rows.Next() {
		// Scan dữ liệu
		err := rows.Scan(&minTimeInDb)
		if err != nil {
			log.Print("Lỗi khi scan dữ liệu:", err)
		}
	}
	// Kiểm tra lỗi sau khi duyệt rows
	if rows.Err() != nil {
		log.Fatal(rows.Err())
	}
	return minTimeInDb
}

func LoadDataResponse(fromDate string, toDate string) []WrapResultAPI {
	//open connection
	var db = OpenNativeConnection()
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
	CloseNativeConnection(db)

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
	if err = rows.Err(); err != nil {
		log.Fatal("Lỗi khi duyệt rows:", err)
	}
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
