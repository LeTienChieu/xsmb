package services

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"myproject/pkg/repositories"
)

func PushToSpreadSheet(sheetName string, startPosition string, values [][]interface{}) {
	var credFile = repositories.GeConfigFileGoogleSheet()
	//credFile := "google-auth/sinuous-crow-433203-m1-3e1f8830eebd.json"

	data, err := ioutil.ReadFile(credFile)
	if err != nil {
		log.Fatalf("Can not read file certificate: %v", err)
	}

	config, err := google.JWTConfigFromJSON(data, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Can not create config: %v", err)
	}

	ctx := context.Background()
	client := config.Client(ctx)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Can not init google sheet service: %v", err)
	}

	// ID of Google Sheet
	spreadsheetId := repositories.GetConfigSheetId()

	// Create ValueRange for write data
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Data position (ví dụ: sheet1, bắt đầu từ ô A1)
	rangeData := sheetName + "!" + startPosition

	// Write data to google sheet
	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, rangeData, valueRange).
		ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Can not write data for google sheet: %v", err)
	}

	fmt.Println("Wirte data to google sheet successfully")
}
