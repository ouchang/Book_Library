package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *sql.DB
var ORMDB *gorm.DB

// DBConfig represents db configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func ConnectMySQLDB() (db *sql.DB, err error) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "123456"
	dbName := "book_library"
	dbPort := "4306"
	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:"+dbPort+")/"+dbName)
	return
}

// https://gorm.io/docs/connecting_to_the_database.html

func ConnectORMMySQLDB(db *sql.DB) (ormdb *gorm.DB, err error) {
	ormdb, err = gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	return
}
