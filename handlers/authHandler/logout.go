package authHandler

import (
	"oui/auth"
	"oui/models"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)

func Logout(c iris.Context, route models.Route) {

	if c.Method() != "POST" || route.Tertiary != "" || len(route.Tail) != 0 {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}

	var id string
	if t := c.GetID(); t != nil {
		id = t.(auth.UserToken).TokenID
	} else {
		c.StopWithJSON(iris.StatusUnauthorized, iris.Map{"message": "Invalid session"})
		return
	}

	if ok := auth.RevokeUserToken(id); ok {
		c.StopWithJSON(iris.StatusCreated, iris.Map{"message": "You are no longer connected"})
	} else {
		golog.Error("impossible de déconnecter l'utilisateur")
		c.StopWithStatus(iris.StatusInternalServerError)
	}
}
