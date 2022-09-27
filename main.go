package main

import (
	"context"
	repository2 "entetry/gotest/internal/auth/repository"
	"entetry/gotest/internal/config"
	"entetry/gotest/internal/handlers"
	"entetry/gotest/internal/repository"
	"entetry/gotest/internal/repository/mongodb"
	"entetry/gotest/internal/repository/postgre"
	"entetry/gotest/internal/service"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	var companyRepository repository.CompanyRepository
	var authRepository repository.Auth
	if cfg.IsMongo {
		db, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.ConnectionString))
		defer func(db *mongo.Client, ctx context.Context) {
			err := db.Disconnect(ctx)
			if err != nil {
				log.Fatal(err)
			}
		}(db, context.Background())
		companyRepository = mongodb.NewCompanyRepository(db.Database("companydb"))
		if err != nil {
			log.Fatalf("Couldn't connect to mongo database: %s", err)
		}
	} else {
		db, err := pgxpool.Connect(ctx, cfg.ConnectionString)
		defer db.Close()
		authRepository := repository2.NewAuthRepository(db)
		companyRepository = postgre.NewCompanyRepository(db)
		if err != nil {
			log.Fatalf("Couldn't connect to database: %s", err)
		}
	}

	e := echo.New()

	authService := postgre.NewAuth(authRepository)
	authHandler := handlers.NewAuth(authService)
	auth := e.Group("api/auth")
	auth.POST("/refresh-tokens", authHandler.Refresh)
	auth.POST("/signin", authHandler.SignIn)
	auth.POST("/signup", authHandler.SignUp)
	auth.POST("/logout", authHandler.Logout)

	companyService := service.NewCompany(companyRepository)
	companyHandler := handlers.NewCompany(companyService)
	company := e.Group("company")
	company.POST("", companyHandler.Create)
	company.GET("", companyHandler.GetAll)
	company.GET("/:id", companyHandler.GetById)
	company.PUT("", companyHandler.Update)
	company.DELETE("/:id", companyHandler.Delete)
	err = e.Start(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
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
