package db

import (
	"Company-Chat-Register/common/config"
	"fmt"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

var DB *gorm.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress, config.Config.Mysql.DBName)

	logs.Info(dsn)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		logs.Error("err open db = %s", err.Error())
		return
	}
	db.SingularTable(true)
	db.DB().SetMaxOpenConns(config.Config.Mysql.DBMaxOpenConns)
	db.DB().SetMaxIdleConns(config.Config.Mysql.DBMaxIdleConns)
	db.DB().SetConnMaxLifetime(time.Duration(config.Config.Mysql.DBMaxLifeTime) * time.Second)
	DB = db
	logs.Info("DB init success...")
}
