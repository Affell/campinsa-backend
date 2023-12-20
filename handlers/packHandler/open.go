package packHandler

import (
	"oui/models/membre"
	"oui/models/user"

	"github.com/kataras/iris/v12"
)

func openPack(token user.UserToken, packToken string) (code int, data interface{}) {

	m := membre.GetRandomMember()
	if m.Id != -1 {
		code, data = iris.StatusOK, m
	} else {
		code, data = iris.StatusInternalServerError, iris.Map{"message": "Internal server error"}
	}

	return
}
