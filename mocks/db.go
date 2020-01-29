package mocks

import (
	"fmt"

	"github.com/airdb/mina-api/model/po"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // nolint:golint
)

func dropRecords(db *gorm.DB) {
	db.DropTableIfExists(&po.User{})
}

func setUpRecords(db *gorm.DB) {
	db.Create(User1)
}

func SetUpMockDatabases() (*gorm.DB, error) {
	// Set up mock database.
	dbName := "dns_db"
	// db, err := gorm.Open("sqlite3", "dev_shopee_space_dns_db")

	db, err := gorm.Open("sqlite3", dbName)
	db.SingularTable(true)
	db.Callback().Delete().Remove("gorm:delete")
	db.Callback().Update().Remove("gorm:update_time_stamp")
	db.Callback().Create().Remove("gorm:update_time_stamp")

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("=====Using `sqlite3` for testing=====")

	// Migrate the schema.
	db.AutoMigrate(&po.User{})

	// Create records.
	setUpRecords(db)
	/*
	   // Hook test DB into dbUtils.
	   err = dbutils.InitTestDB(dbName, db)
	   if err != nil {
	           return nil, err
	   }
	*/

	return db, nil
}

func DestroyMockDatabases(db *gorm.DB) {
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	// 	defer dbutils.ReleaseTestDB()
	defer dropRecords(db)
}
