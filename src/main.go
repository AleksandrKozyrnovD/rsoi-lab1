package main

import (
	"context"
	"fmt"
	"lab1/api"
	"lab1/config"
	"lab1/database"
	"lab1/logger"
	"lab1/middleware"
	"lab1/repository"
	"lab1/service"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const configPath string = "./config.yml"

func main() {
	logger.Init() //pass settings and/or name/path

	config, err := config.InitConfig(configPath)
	if err != nil {
		slog.Error("init config", err)
		return
	}

	urlDatabase := database.GetConnectionString(&config.Database)
	database := database.Init(urlDatabase)
	personsRepository := repository.NewPersonRepository(database)
	personService := service.NewPersonService(personsRepository)

	controller := api.NewV1(personService)
	port := config.HTTP.Port

	r := gin.New()

	r.Use(middleware.SlogMiddleware())
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/persons", controller.HandlePersonsGet)
		v1.POST("/persons", controller.HandlePersonsPost)
		v1.GET("/persons/:id", controller.HandlePersonsGetId)
		v1.DELETE("/persons/:id", controller.HandlePersonsDeleteId)
		v1.PATCH("/persons/:id", controller.HandlePersonsPatchId)
	}

	slog.Info("Running at ", config.HTTP.Host, config.HTTP.Port)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port),
		Handler: r,
	}

	go func() {
		slog.Info("Server running", slog.Int64("Port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exited")
}
