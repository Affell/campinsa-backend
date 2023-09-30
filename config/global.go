package config

import (
	"oui/email"
)

type Config struct {
	App
	Email email.Config
	Redis
}

var Cfg Config
