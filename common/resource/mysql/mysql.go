package mysql

import (
	"fmt"
	"frozen-go-cms/common/mylogrus"
	"frozen-go-cms/common/resource/config"
	_ "github.com/go-sql-driver/mysql" //加载mysql驱动
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"net/url"
	"time"
)

var Db *gorm.DB

func init() {
	var err error
	mysqlConfigData := config.GetConfigMysql()
	options := "?charset=utf8mb4&parseTime=True&loc=Local&time_zone=" + url.QueryEscape("'+8:00'")
	dsn := "" + mysqlConfigData.MYSQL_USERNAME + ":" + mysqlConfigData.MYSQL_PASSWORD + "@(" + mysqlConfigData.MYSQL_HOST + ")/" + mysqlConfigData.MYSQL_DB + options

	sqlLogger := logger.Default.LogMode(logger.Info)
	if file := mylogrus.GetSqlLog(); file != nil {
		//sqlLogger = logger.New(log.New(file, "\r\n", log.Ldate|log.Lmicroseconds), logger.Config{
		sqlLogger = MyNew(log.New(file, "", log.Ldate|log.Lmicroseconds), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      false,
		})
	}

	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: sqlLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("mysql connect error %v", err)
	} else {
		log.Println("mysql connect success")
	}

	if Db.Error != nil {
		fmt.Printf("database error %v", Db.Error)
	}

	if d, err := Db.DB(); err == nil {
		d.SetConnMaxLifetime(time.Minute * 30) // 连接可复用的最大时间。
		d.SetMaxIdleConns(300)                 // 空闲连接数
		d.SetMaxOpenConns(300)                 // 最大连接数
		if err := d.Ping(); err != nil {
			fmt.Printf("database ping error %v", err)
		}
	}
}
