package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/nanpipat/golang-template-hexagonal/configs"
	"github.com/nanpipat/golang-template-hexagonal/database"
	"github.com/nanpipat/golang-template-hexagonal/migrations"
	"github.com/nanpipat/golang-template-hexagonal/package/logger"
)

type config struct {
	Env string
}

func MigrateRun() {
	var cfg config

	flag.StringVar(&cfg.Env, "env", "", "the environment to use")
	flag.Parse()

	configs.InitViper("./configs", cfg.Env)

	db, err := database.ConnectToDB(
		configs.GetViper().DB.Host,
		configs.GetViper().DB.Port,
		configs.GetViper().DB.Username,
		configs.GetViper().DB.Password,
		configs.GetViper().DB.DBName,
		configs.GetViper().DB.DBType,
	)

	if err != nil {
		logger.Error(err.Error())
		panic(err.Error())
	}
	// Graceful shutdown ...
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("Gracefull shut down ...")
			//TODO: close database or any connection before server has gone ...
			database.DisconnectDatabase(db.DB)
			if err != nil {
				panic("Can't shutdown")
			}
		}
	}()

	err = migrations.RunMigrations(db.DB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Migrate: %v", err)
		os.Exit(1)
	}

}
