package database

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
)

var config = viper.New()

func DatabaseInit() *gorm.DB {
	config.SetConfigFile("config.env")
	config.AddConfigPath(".")
	config.AutomaticEnv()

	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}

	host := config.GetString("DB_HOST")
	port := config.GetInt("DB_PORT")
	username := config.GetString("DB_USERNAME")
	password := config.GetString("DB_PASSWORD")
	database := config.GetString("DB_DATABASE")

	dsn := username + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}
