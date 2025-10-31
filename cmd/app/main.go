package main

import (
	"github.com/bmkersey/go-gui-dean/internal/db"
	"github.com/bmkersey/go-gui-dean/internal/gui"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	sqlDB := db.Connect()
	defer sqlDB.Close()

	gui.RunApp(sqlDB)

}
