package database

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

var connectDB = &Database{}

func ConnectToDB(host, port, username, pass, dbname, dbtype string) (*Database, error) {
	fmt.Println("User is: ", username)
	fmt.Println("Password is: ", pass)
	fmt.Println("Host is: ", host)
	fmt.Println("Port is: ", port)
	fmt.Println("Dbname is: ", dbname)
	fmt.Println("DbType is: ", dbtype)

	var dial gorm.Dialector

	if host == "" && port == "" && dbname == "" {
		return nil, errors.New("cannot estabished the connection")
	}

	switch dbtype {
	case "postgres":
		dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v", host, username, pass, dbname, port)
		dial = postgres.Open(dsn)
	case "mysql":
		dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", username, pass, host, port, dbname)
		dial = mysql.Open(dsn)
	case "sqlserver":
		dsn := fmt.Sprintf("sqlserver://%v:%v@%v:%v?database=%v", username, pass, host, port, dbname)
		dial = sqlserver.Open(dsn)
	default:
		fmt.Println("not found database type")
	}

	var err error
	con, err := gorm.Open(dial, &gorm.Config{
		DryRun: false,
	})
	if err != nil {
		panic(err)
	}

	connectDB.DB = con
	return connectDB, nil
}

func DisconnectDatabase(db *gorm.DB) {
	sqlDb, err := db.DB()
	if err != nil {
		panic("close db")
	}
	sqlDb.Close()
	log.Println("Connected with db has closed")
}
