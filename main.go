package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"entetry/gotest/internal/config"
	"entetry/gotest/internal/handlers"
	"entetry/gotest/internal/middleware"
	"entetry/gotest/internal/repository/postgre"
	"entetry/gotest/internal/service"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	jwtCfg, err := config.NewJwtConfig()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)
	db, err := pgxpool.Connect(ctx, cfg.ConnectionString)
	defer cancel()
	defer db.Close()
	if err != nil {
		log.Fatalf("Couldn't connect to database: %s", err)
	}

	refreshSessionRepository := postgre.NewRefresh(db)
	refreshSessionService := service.NewRefreshSession(refreshSessionRepository)

	userRepository := postgre.NewUserRepository(db)
	userService := service.NewUserService(userRepository)

	authService := service.NewAuthService(userService, refreshSessionService, jwtCfg)
	authHandler := handlers.NewAuth(authService)

	companyRepository := postgre.NewCompanyRepository(db)
	companyService := service.NewCompany(companyRepository)
	companyHandler := handlers.NewCompany(companyService)

	e := echo.New()
	e.Validator = middleware.NewCustomValidator(validator.New())
	auth := e.Group("api/auth")
	auth.POST("/refresh-tokens", authHandler.Refresh)
	auth.POST("/sign-in", authHandler.SignIn)
	auth.POST("/sign-up", authHandler.SignUp)
	auth.POST("/logout", authHandler.Logout)

	company := e.Group("company")
	company.Use(middleware.NewJwtMiddleware(jwtCfg.AccessTokenKey))
	company.POST("", companyHandler.Create)
	company.GET("", companyHandler.GetAll)
	company.GET("/:id", companyHandler.GetByID)
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
