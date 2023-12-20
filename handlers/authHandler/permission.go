package authHandler

import (
	"oui/models"
	"oui/models/user"

	"github.com/kataras/iris/v12"
)

func Permission(ctx iris.Context, route models.Route) {

	if ctx.Method() != "GET" || route.Secondary != "" || route.Tertiary != "" || len(route.Tail) > 0 {
		return
	}

	var queries PermissionQueries
	if err := ctx.ReadQuery(&queries); err != nil || queries.Permission == "" {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	var result bool

	var token user.UserToken
	if t := ctx.GetID(); t != nil {
		token = t.(user.UserToken)

		result = user.HasPermission(token.ID, queries.Permission)
	} else {
		result = false
	}

	ctx.StopWithJSON(iris.StatusOK, iris.Map{"result": result})

}
