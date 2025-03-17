package main

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"myproject/pkg/services"
	"time"
)

func main() {
	log.Print("Server is starting...")
	jobRunning()
	//c := cron.New()
	//_, err := c.AddFunc("36 21 * * *", func() {
	//	fmt.Println("Cron job started at: ", time.Now())
	//	jobRunning()
	//})
	//if err != nil {
	//	log.Fatalf("Lỗi khi thêm cron job: %v", err)
	//}
	//c.Start()
	//// Wait forever
	//select {}
}
func jobRunning() {
	//prepare data
	listTime := services.CalculateMaxTimeWithDb()
	fmt.Printf("listTime: %v\n", listTime)
	for i := 0; i < len(listTime); i++ {
		time.Sleep(2 * time.Second)
		responseHtml := services.FetchData(listTime[i])
		services.StoreDataFromResponse(responseHtml)
	}
	services.GetTopStartNumberBest2024V2()
	services.GetTopStartNumberBestCurentMonthV2()
	services.GetTopStartNumberBest2025V2()
	services.GetTopNumberBestCurrentMonthV2()
	services.GetTopNumberBest2025V2()
	services.GetTopNumberBestYear2024V2()
}
