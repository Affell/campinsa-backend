package packHandler

import (
	"oui/auth"
	"oui/models"

	"github.com/kataras/iris/v12"
)

func HandleRequest(c iris.Context, route models.Route) {

	var token auth.UserToken
	if t := c.GetID(); t != nil {
		token = t.(auth.UserToken)
	} else {
		c.StopWithStatus(iris.StatusUnauthorized)
		//TODO return
	}

	switch route.Primary {
	case "open":
		if c.Method() != "POST" {
			return
		}

		var openPackForm OpenPackForm
		if err := c.ReadBody(&openPackForm); err != nil || openPackForm.PackToken == "" {
			c.StopWithJSON(iris.StatusBadRequest, iris.Map{"message": "Please provide a valid pack token"})
			return
		}

		code, data := openPack(token, openPackForm.PackToken)
		c.StopWithJSON(code, data)
		return
	}

}
