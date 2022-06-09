package database

import (
	"fmt"
	"go-search-engine/src/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var MySqlDb *gorm.DB
var MySqlError error

func init() {
	dbConfig := config.GetMySQLConfig()

	// set database dsn
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s",
		dbConfig.USER,
		dbConfig.PASSWORD,
		dbConfig.DB_HOST,
		dbConfig.DB_PORT,
		dbConfig.DB_NAME,
		dbConfig.CHARSET,
		dbConfig.ParseTime,
	)

	// open connection
	MySqlDb, MySqlError = gorm.Open("mysql", dbDSN)
	// Enable Logger, show detailed log
	MySqlDb.LogMode(true)

	db := MySqlDb.DB()

	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)

	// 禁用默认复数表名
	MySqlDb.SingularTable(true)

	if MySqlError != nil {
		panic("database open error! " + MySqlError.Error())
	}
}
