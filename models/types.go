package models

import (
	"github.com/kataras/iris/v12"
)

type (
	HandlerMap map[string]func(iris.Context, Route)
)
