package main

import (
	db2 "labsystem/db"
	"labsystem/model"
)

func main() {
	db := db2.NewMySQL()
	for _, v := range model.Models {
		if err := db.DB.AutoMigrate(v); err != nil {
			// TODO print log
		}
	}
}
