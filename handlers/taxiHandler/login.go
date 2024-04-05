package taxiHandler

import (
	"oui/models"
	"oui/models/user"

	"github.com/kataras/iris/v12"
)

func Login(c iris.Context, route models.Route) {

	if c.Method() != "POST" || route.Secondary != "" {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}

	var loginForm LoginForm
	if err := c.ReadBody(&loginForm); err != nil || loginForm.Token == "" {
		c.StopWithJSON(iris.StatusBadRequest, iris.Map{"message": "Please fully fill in the login form"})
		return
	}

	currentUser, err := user.GetUserByTaxiToken(loginForm.Token)
	if err != nil {
		c.StopWithJSON(iris.StatusUnauthorized, iris.Map{"message": "Invalid token"})
		return
	}

	c.StopWithJSON(iris.StatusOK, iris.Map{"user": currentUser.ToSelfWebDetail()})
}
