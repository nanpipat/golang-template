package protocol

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/nanpipat/golang-template-hexagonal/configs"
	"github.com/nanpipat/golang-template-hexagonal/database"
	"github.com/nanpipat/golang-template-hexagonal/internal/core/services"
	"github.com/nanpipat/golang-template-hexagonal/internal/handlers"
	"github.com/nanpipat/golang-template-hexagonal/internal/repo"
	"github.com/nanpipat/golang-template-hexagonal/package/logger"
	"github.com/nanpipat/golang-template-hexagonal/protocol/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type config struct {
	Env string
}

func HTTPStart() error {
	app := fiber.New()
	app.Use(cors.New())

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
			err := app.Shutdown()
			if err != nil {
				panic("Can't shutdown")
			}
		}
	}()

	app.Get("/healthz", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	api := app.Group("/api")
	{
		userRepo := repo.NewUserRepository(db.DB)
		userService := services.NewUserService(userRepo)
		userHandlers := handlers.NewUserHandlers(userService)
		routes.UserRoutes(api, userHandlers)
	}

	err = app.Listen(":" + configs.GetViper().App.Port)
	if err != nil {
		return err
	}
	return nil
}
