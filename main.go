package main

import (
	"log"
	"os"

	"github.com/nanpipat/golang-template-hexagonal/cmd"
	"github.com/nanpipat/golang-template-hexagonal/consts"
)

func main() {
	switch os.Getenv("APP_SERVICE") {
	case consts.ServiceAPI:
		cmd.ApiRun()
	case consts.ServiceMigration:
		cmd.MigrateRun()
	default:
		log.Fatal("Service not support")
	}
}
