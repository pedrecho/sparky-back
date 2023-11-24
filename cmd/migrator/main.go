package main

import (
	"context"
	"fmt"
	"os"
	"sparky-back/internal/config"
	"sparky-back/internal/loader"
	"sparky-back/internal/models"
)

const configEnv = "CONFIG"

func main() {
	configPath, ok := os.LookupEnv(configEnv)
	if !ok {
		panic("no config env")
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("config initialization: %v", err))
	}
	db := loader.New(cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	_, err = db.NewCreateTable().
		Model(&models.Reaction{}).
		IfNotExists().
		Exec(context.Background())
	_, err = db.NewCreateTable().
		Model(&models.User{}).
		IfNotExists().
		Exec(context.Background())
	fmt.Println("successful migration")
}
