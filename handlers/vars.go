package handlers

import (
	"oui/handlers/packHandler"
	"oui/handlers/shotgunHandler"
	"oui/models"
)

var (
	Services models.HandlerMap = models.HandlerMap{
		shotgunHandler.Service: shotgunHandler.HandleRequest,
		packHandler.Service:    packHandler.HandleRequest,
	}
)
