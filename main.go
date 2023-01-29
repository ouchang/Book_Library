package main

import (
	"fmt"
	"gobooklibrary/config"
	"gobooklibrary/routes"
	"log"
)

var err error

func main() {

	config.DB, err = config.ConnectMySQLDB()
	if err != nil {
		fmt.Println("Cannot connect to MySQL database")
		log.Fatal("connection error:", err)
	}
	config.ORMDB, err = config.ConnectORMMySQLDB(config.DB)

	if err != nil {
		fmt.Println(err)
	}

	defer config.DB.Close()

	r := routes.SetupRouter()
	r.Run()
}
