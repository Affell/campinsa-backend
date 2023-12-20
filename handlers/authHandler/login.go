package authHandler

import (
	"oui/models"
	"oui/models/user"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)

func Login(c iris.Context, route models.Route) {

	if c.Method() != "POST" || route.Tertiary != "" || len(route.Tail) != 0 {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}

	var loginForm LoginForm
	if err := c.ReadBody(&loginForm); err != nil || ((loginForm.Password == "" || loginForm.Username == "") && loginForm.Token == "") {
		c.StopWithJSON(iris.StatusBadRequest, iris.Map{"message": "Please fully fill in the login form"})
		return
	}

	var CurrentUserToken user.UserToken
	var err error
	if loginForm.Token == "" && loginForm.Password != "" && loginForm.Username != "" {
		CurrentUserToken, err = user.GetSQLUserToken(loginForm.Username, loginForm.Password)
		if err != nil {
			c.StopWithJSON(iris.StatusForbidden, iris.Map{"message": "Invalid username or password"})
			return
		}
	} else if loginForm.Token != "" && loginForm.Password == "" && loginForm.Username == "" {
		CurrentUserToken, err = user.GetUserToken(loginForm.Token)
		if err != nil {
			c.StopWithJSON(iris.StatusForbidden, iris.Map{"message": "Invalid token"})
			return
		}
	} else {
		c.StopWithStatus(iris.StatusBadRequest)
		return
	}

	TokenID := CurrentUserToken.Store(loginForm.RememberMe)
	if TokenID == "" {
		golog.Error("during store access token memory")
		c.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	c.StopWithJSON(iris.StatusOK, iris.Map{"token": TokenID, "user": CurrentUserToken.ToUserData()})
}
