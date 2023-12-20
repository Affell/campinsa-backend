package authHandler

import (
	"oui/models"
	"oui/models/user"

	"github.com/kataras/iris/v12"
)

func Me(c iris.Context, route models.Route) {

	if c.Method() != "GET" || route.Secondary != "" || route.Tertiary != "" || len(route.Tail) > 0 {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}

	var token user.UserToken
	if t := c.GetID(); t != nil {
		token = t.(user.UserToken)
	} else {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}

	c.Redirect("/auth/user/"+token.Email, iris.StatusPermanentRedirect)
}

func getUser(token user.UserToken, email string) (code int, data interface{}) {
	u, err := user.GetUserByEmail(email)
	if err == "" {
		code = iris.StatusOK
		var m map[string]interface{}
		if token.IsNil() || token.Email == email {
			code = iris.StatusUnauthorized
			return
		}
		m = u.ToSelfWebDetail()
		if token.ID != 0 && user.HasPermission(token.ID, user.PERMISSIONS_EDIT_PERMISSION) {
			m["permission"] = user.GetUserPermissions(u.ID)
		}
		data = m
	} else {
		code, data = iris.StatusNotFound, iris.Map{"message": err}
	}

	return
}
