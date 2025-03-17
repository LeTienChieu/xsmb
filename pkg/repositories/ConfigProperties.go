package repositories

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func GetDriverClassName() string {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	// read config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	// get value from config file
	return viper.GetString("db.driver")
}

func GetDataSourceStr() string {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")

	// read config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// get value from config file
	dataSourceStr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.dbname"))
	return dataSourceStr
}

func GeConfigFileGoogleSheet() string {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	return viper.GetString("config.sheet.file")
}

func GetConfigSheetId() string {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	return viper.GetString("config.sheet.id")
}

func GetConfigSourceUrl() string {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	return viper.GetString("config.data.url")
}
