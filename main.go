package main

import (
	"embed"
	"fmt"
	"oui/auth"
	"oui/config"
	"oui/handlers"
	"oui/handlers/ws"
	"oui/models/planning"
	"oui/models/shotgun"
	"strings"
	"time"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)

//go:embed public/*
var Folder embed.FS

func main() {
	planning.RetrievePlanning()
	return
	// Lancement des routines
	golog.Info("launching routines...")
	shotgun.LoadShotgunsIntoCache()

	// creation du nouveau Engine Iris
	router := iris.New()

	// DEBUG : utilisation du iris.Logger par défaut sinon rien
	router.Logger().SetLevel(config.Cfg.App.DebugLevel)
	if config.Cfg.App.DebugLevel == "debug" {
		router.Use(debugLogger())
	}

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "DELETE"},
		AllowCredentials: true,
	})
	router.UseRouter(crs)
	router.AllowMethods(iris.MethodOptions)
	router.Use(iris.Compression)

	router.Get("/favicon.ico", iris.Cache(30*time.Second), handlers.FaviconHandler(Folder))

	router.Get("/", iris.Cache(15*time.Second), handlers.IndexHandler)

	router.HandleDir("/img", handlers.FileHandler(Folder, "public/img"), iris.DirOptions{IndexName: "/img", Compress: true})

	wsRouter := ws.NewRouter()
	wsRouter.On("mode", ws.OnMode, false)
	wsRouter.On("planning", ws.OnPlanning, true)
	wsRouter.On("currentLocation", ws.OnCurrentLocation, true)
	wsRouter.On("askTaxiLocation", ws.OnAskTaxiLocation, false)
	wsRouter.On("stopLocation", ws.OnStopLocation, true)
	wsRouter.On("newRide", ws.OnNewRide, true)
	wsRouter.On("rides", ws.RetrieveRides, true)
	wsRouter.On("rideAnswer", ws.OnRideAnswer, true)
	wsRouter.On("rideCompleted", ws.OnRideCompleted, true)
	router.Any("/ws/caritaxi", wsRouter.ServeHTTP)

	router.Use(auth.AuthRequired())
	{
		router.Any("/{service}", handlers.HandleRequest)
		router.Any("/{service}/{primary}", handlers.HandleRequest)
		router.Any("/{service}/{primary}/{secondary}", handlers.HandleRequest)
		router.Any("/{service}/{primary}/{secondary}/{tertiary}", handlers.HandleRequest)
		router.Any("/{service}/{primary}/{secondary}/{tertiary}/{tail}", handlers.HandleRequest)
	}

	err := router.Listen(":" + config.Cfg.App.Port)
	if err != nil {
		golog.Fatal(err)
	}
}

func debugLogger() iris.Handler {
	return func(c iris.Context) {
		t := time.Now()
		c.Next()
		params := []string{
			fmt.Sprint(c.GetStatusCode()),
			c.Request().Method,
			c.RequestPath(false),
			time.Since(t).String(),
		}
		golog.Debug(strings.Join(params, " | "))
	}
}
