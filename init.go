package main

import (
	"oui/config"
	"time"

	"oui/models/postgresql"

	"github.com/kataras/golog"
	"github.com/provectio/godotenv"
)

func init() {
	// Load Env vars from .env if exist
	err := godotenv.Load()
	if err != nil {
		golog.Warn("No .env file to load")
	}

	config.Cfg.App = config.InitApp()
	golog.SetLevel(config.Cfg.App.DebugLevel)

	postgresql.SQLCtx, postgresql.SQLConn = config.InitPgSQL()
	config.Cfg.Email = config.InitEmailConfig(Folder)
	config.Cfg.Redis = config.InitRedis()
	golog.Debug("init success in " + time.Since(config.Cfg.App.StartedTime).String())
}
