package authHandler

import (
	"oui/models"
)

const (
	Service string = "auth"
)

var (
	Handlers models.HandlerMap = models.HandlerMap{
		"login":      Login,
		"logout":     Logout,
		"me":         Me,
		"permission": Permission,
	}
)

type (
	LoginForm struct {
		Username   string `form:"username" json:"username"`
		Password   string `form:"password" json:"password"`
		Token      string `form:"token" json:"token"`
		RememberMe bool   `form:"remember_me" json:"remember_me"`
	}

	AskRecoverForm struct {
		Email string `form:"email" json:"email" binding:"required"`
	}

	RecoverForm struct {
		Password string `form:"password" json:"password" binding:"required"`
	}

	PermissionQueries struct {
		Permission string `form:"permission" json:"permission" binding:"required"`
	}
)
