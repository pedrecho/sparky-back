package app

import (
	"errors"
	"fmt"
	"github.com/uptrace/bunrouter"
	"net/http"
	"sparky-back/internal/config"
	"sparky-back/internal/controllers"
	"sparky-back/internal/loader"
	"sparky-back/internal/logic"
	"sparky-back/internal/middlewares"
	"sparky-back/pkg/zaplogger"
)

func Run(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("config initialization: %w", err)
	}
	zapsync, err := zaplogger.ReplaceZap(cfg.Logger)
	if err != nil {
		return fmt.Errorf("zaplogger initialization: %w", err)
	}
	defer zapsync()
	l := logic.NewLogic(loader.New(cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName))
	c := controllers.New(l)

	router := bunrouter.New(
		bunrouter.Use(middlewares.Log),
	)
	router.POST("/signup", c.AddUser)
	router.POST("/signin", c.Login)
	router.POST("/update", c.UpdateUser)
	router.GET("/user/:id", c.GetUserByID)
	router.GET("/static/:filename", c.GetFile)
	router.POST("/reaction", c.SetReaction)
	router.POST("/connection", c.ClientConnection)
	router.POST("/message", c.NewMessage)
	router.POST("/recommendations", c.GetRecommendations)
	handler := http.HandlerFunc(router.ServeHTTP)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler,
	}
	if err = httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
