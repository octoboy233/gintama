package goft

import (
	"github.com/jinzhu/gorm"
	"log"
)

type GormAdapter struct {
	*gorm.DB
}

func (this *GormAdapter) Name() string {
	return "GormAdapter"
}

func NewGormAdapter() *GormAdapter {
	db, err := gorm.Open("mysql",
		"root:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(10)
	return &GormAdapter{DB: db}
}
