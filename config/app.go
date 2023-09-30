package config

import (
	"os"
	"time"
)

type App struct {
	Name         string
	Version      string
	Port         string
	InternalPort string
	FrontURL     string
	DebugLevel   string
	StartedTime  time.Time
}

func InitApp() (app App) {
	app.StartedTime = time.Now()

	if env := os.Getenv("APP_NAME"); env != "" {
		app.Name = env
	} else {
		app.Name = "API - XConnect"
	}

	if env := os.Getenv("APP_VERSION"); env != "" {
		app.Version = env
	} else {
		app.Version = "v1"
	}

	if env := os.Getenv("APP_LOG_LEVEL"); env != "" {
		app.DebugLevel = env
	} else {
		app.DebugLevel = "debug"
	}

	// For sending email for example
	if env := os.Getenv("APP_FRONT_URL"); env != "" {
		app.FrontURL = env
	} else {
		app.FrontURL = "http://localhost:3000"
	}

	if env := os.Getenv("APP_PORT"); env != "" {
		app.Port = env
	} else {
		app.Port = "4000"
	}

	if env := os.Getenv("APP_INTERNAL_PORT"); env != "" {
		app.InternalPort = env
	} else {
		app.InternalPort = "5000"
	}

	return
}
