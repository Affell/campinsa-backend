package shotgunHandler

import "oui/models"

const Service = "shotgun"

var (
	Handlers models.HandlerMap = models.HandlerMap{
		"list": List,
		"link": GetLink,
	}
)
