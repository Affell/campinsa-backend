package auth

import (
	"oui/models/exceptions"
	"oui/models/user"

	"github.com/kataras/iris/v12"
)

const (
	TokenKeyName = "Oui-Connect-Token"
)

// RouteVars:
// is a variable object for the dynamic routing balancing.
type RouteVars struct {
	Name   string `param:"package"`
	First  string `param:"first"`
	Second string `param:"second"`
}

type Header struct {
	TokenID     string `header:"Oui-Connect-Token"`
	TaxiTokenID string `header:"CariTaxi-Connect-Token"`
}

func AuthRequired() func(iris.Context) {
	return func(c iris.Context) {

		var header Header
		if err := c.ReadHeaders(&header); err != nil || len(header.TokenID) != 36 {
			c.Next()
			return
		}

		Token, err := user.GetUserToken(header.TokenID)
		if err != nil {
			c.Next()
			return
		}

		var route RouteVars
		if err := c.ReadParams(&route); err != nil {
			c.StopWithJSON(exceptions.StatusCode(iris.StatusNotFound))
			return
		}

		c.SetID(Token)
		c.Next()
	}
}

func CheckTaxiAuth(c iris.Context) (connected bool) {

	var header Header
	if err := c.ReadHeaders(&header); err != nil || len(header.TaxiTokenID) != 36 {
		return
	}

	user, err := user.GetUserByTaxiToken(header.TaxiTokenID)
	if err != nil {
		return
	}

	c.SetID(user)
	return true
}
