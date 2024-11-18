package main

import "gorm.io/gorm"

type TokenHistory struct {
	Id            int `gorm:"primaryKey"`
	TokenId       int
	Name          string
	Date          string
	Buy_Price     float64
	Count         float64
	Current_Price float64
}

type Infos struct {
	Id            int    `gorm:"primaryKey"`
	Name          string `gorm:"unique"`
	Date          string
	Avg_Price     float64
	Count         float64
	Current_Price float64
	Histories     []TokenHistory `gorm:"foreignKey:TokenId"`
}

type API struct {
	Id  int
	Api string
}

var tokens []Infos
var histories []TokenHistory

func initializeDatabase(DB *gorm.DB) {
	DB.AutoMigrate(&Infos{})
	DB.AutoMigrate(&TokenHistory{})
}
