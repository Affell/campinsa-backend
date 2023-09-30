package handlers

import (
	"oui/models"

	"github.com/kataras/iris/v12"
)

func HandleRequest(ctx iris.Context) {
	var route models.Route
	if err := ctx.ReadParams(&route); err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}

	ctx.StatusCode(iris.StatusNotFound)

	if handler, success := Services[route.Service]; success {
		handler(ctx, route)
	}
}
