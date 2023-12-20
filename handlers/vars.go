package handlers

import (
	"oui/handlers/authHandler"
	"oui/handlers/packHandler"
	"oui/handlers/shotgunHandler"
	"oui/handlers/taxiHandler"
	"oui/models"
)

var (
	Services models.HandlerMap = models.HandlerMap{
		authHandler.Service:    authHandler.HandleRequest,
		shotgunHandler.Service: shotgunHandler.HandleRequest,
		packHandler.Service:    packHandler.HandleRequest,
		taxiHandler.Service:    taxiHandler.HandleRequest,
	}
)
