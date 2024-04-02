package main

import (
	"oui/config"
	"oui/models/planning"
	"oui/models/postgresql"
	"oui/models/ride"
	"time"

	"github.com/kataras/golog"
	"github.com/provectio/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		golog.Warn("No .env file to load")
	}

	config.Cfg.App = config.InitApp()
	golog.SetLevel(config.Cfg.App.DebugLevel)

	postgresql.SQLCtx, postgresql.SQLConn = config.InitPgSQL()
	config.Cfg.Email = config.InitEmailConfig(Folder)
	config.Cfg.Redis = config.InitRedis()
	ride.Init()
	planning.InitPlanning()
	golog.Debug("init success in " + time.Since(config.Cfg.App.StartedTime).String())
}
