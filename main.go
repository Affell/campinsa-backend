package main

import (
	"embed"
	"fmt"
	"oui/auth"
	"oui/config"
	"oui/handlers"
	"oui/handlers/ws"
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

	cariTaxiWsRouter := ws.NewRouter(true)
	cariTaxiWsRouter.On("currentLocation", ws.OnCurrentLocation)
	cariTaxiWsRouter.On("askTaxiLocation", ws.OnAskTaxiLocation)
	cariTaxiWsRouter.On("stopLocation", ws.OnStopLocation)
	cariTaxiWsRouter.On("newRide", ws.OnNewRide)
	cariTaxiWsRouter.On("rides", ws.RetrieveRides)
	cariTaxiWsRouter.On("rideAnswer", ws.OnRideAnswer)
	cariTaxiWsRouter.On("rideCompleted", ws.OnRideCompleted)
	router.Any("/ws/caritaxi/{token}", cariTaxiWsRouter.ServeHTTP)

	webWsRouter := ws.NewRouter(false)
	webWsRouter.On("updateTaxiLocation", ws.OnWebUpdateTaxiLocation)
	router.Any("/ws/web", webWsRouter.ServeHTTP)

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
