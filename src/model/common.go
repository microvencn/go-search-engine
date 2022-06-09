package model

import (
	"github.com/jinzhu/gorm"
	"go-search-engine/src/database"
)

var db *gorm.DB = database.MySqlDb
