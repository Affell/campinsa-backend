package packHandler

import (
	"oui/auth"
	"oui/models/membre"

	"github.com/kataras/iris/v12"
)

func openPack(token auth.UserToken, packToken string) (code int, data interface{}) {

	m := membre.GetRandomMember()
	if m.Id != -1 {
		code, data = iris.StatusOK, m
	} else {
		code, data = iris.StatusInternalServerError, iris.Map{"message": "Internal server error"}
	}

	return
}
