package main

import (
	"context"
	"entetry/gotest/internal/config"
	"entetry/gotest/internal/handlers"
	"entetry/gotest/internal/repository"
	"entetry/gotest/internal/service"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)
	db, err := pgxpool.Connect(ctx, cfg.ConnectionString)
	if err != nil {
		log.Fatalf("Couldn't connect to database: %s", err)
	}

	e := echo.New()
	companyRepository := repository.NewCompanyRepository(db)
	companyService := service.NewCompany(companyRepository)
	companyHandler := handlers.NewCompany(companyService)
	company := e.Group("company")
	company.POST("", companyHandler.Create)
	company.GET("/:id", companyHandler.GetById)
	company.PUT("", companyHandler.Update)
	company.DELETE("/:id", companyHandler.Delete)
	e.Start(fmt.Sprintf(":%d", cfg.Port))
	log.Info("Server started on ", cfg.Port)
	go func() {
		<-sigChan
		cancel()
		err = e.Shutdown(ctx)
		if err != nil {
			log.Errorf("can't stop server gracefully %v", err)
		}
	}()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
	err = e.Server.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}

}
