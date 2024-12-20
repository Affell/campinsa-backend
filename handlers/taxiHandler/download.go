package taxiHandler

import (
	"os"
	"oui/models"

	"github.com/kataras/iris/v12"
)

func Download(c iris.Context, route models.Route) {

	if c.Method() != "GET" || route.Tertiary != "" || route.Secondary == "" {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}

	switch route.Secondary {
	case "android":
		c.SendFile("public/caritaxi/caritaxi.apk", "caritaxi.apk")
	case "ios":
		c.Redirect(os.Getenv("APPLE_STORE_LINK"), iris.StatusFound)
	default:
		c.StopWithStatus(iris.StatusNotFound)
	}

}
