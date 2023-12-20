package taxiHandler

import (
	"oui/models"

	"github.com/kataras/iris/v12"
)

func HandleRequest(ctx iris.Context, route models.Route) {
	if route.Primary == "" {
		return
	}

	if handler, success := Handlers[route.Primary]; success {
		handler(ctx, route)
	}
}
