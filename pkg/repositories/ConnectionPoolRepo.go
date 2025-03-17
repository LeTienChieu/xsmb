package repositories

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"log"
)

func OpenOrmConnection() *gorm.DB {
	db, err := gorm.Open(GetDriverClassName(), GetDataSourceStr())
	if err != nil {
		log.Fatal("Can not create connection to database:", err)
	}
	return db
}

func CloseOrmConnection(db *gorm.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal("Error when close connection")
	}
}

func OpenNativeConnection() *sql.DB {
	db, err := sql.Open(GetDriverClassName(), GetDataSourceStr())
	if err != nil {
		log.Fatal("Can not connect to database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Can not ping database:", err)
	}
	return db
}

func CloseNativeConnection(db *sql.DB) {
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Can not close database connection:", err)
		}
	}(db)
}
