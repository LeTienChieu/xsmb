package repositories

import (
	"time"
)

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

func autoInitTable() {
	db := OpenOrmConnection()
	db.AutoMigrate(&DayOfXsmb{}, &DayOfXsmbDetail{})
	CloseOrmConnection(db)
}

func storeData(dayOfXsmb DayOfXsmb) {
	db := OpenOrmConnection()
	db.Create(&dayOfXsmb)
	CloseOrmConnection(db)
}
