package dao

import (
	"database/sql"
	"fmt"
	"github.com/Axope/JOJ/configs"
	"github.com/Axope/JOJ/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func autoMigrate() error {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

	return nil
}

func InitMysql() error {
	cfg := configs.GetMysqlConfig()
	username := cfg.Username
	password := cfg.Password
	host := cfg.Host
	port := cfg.Port
	name := cfg.Name
	timeout := cfg.Timeout

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		username, password, host, port, name, timeout)

	var err error
	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		return err
	}

	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)

	autoMigrate()

	return nil
}

func GetMysql() *gorm.DB {
	return db
}
