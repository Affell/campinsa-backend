package config

import (
	"embed"
	"html/template"
	"os"
	"oui/email"

	"github.com/kataras/golog"
)

// InitEmailConfig:
// This function Init the configuration for Emails
func InitEmailConfig(folder embed.FS) (config email.Config) {
	if env := os.Getenv("SMTP_HOST"); env != "" {
		config.Host = env
	} else {
		golog.Fatal("'SMTP_HOST' not correct in env")
	}

	if env := os.Getenv("SMTP_PORT"); env != "" {
		config.Port = env
	} else {
		golog.Fatal("'SMTP_PORT' not correct in env")
	}

	if env := os.Getenv("SMTP_USER"); env != "" {
		config.User = env
	} else {
		golog.Fatal("'SMTP_USER' not correct in env")
	}

	if env := os.Getenv("SMTP_PASSWORD"); env != "" {
		config.Password = env
	} else {
		golog.Fatal("'SMTP_PASSWORD' not correct in env")
	}

	if env := os.Getenv("SMTP_DISPLAYNAME"); env != "" {
		config.From = env
	} else {
		golog.Fatal("'SMTP_DISPLAYNAME' not correct in env")
	}

	config.Template = template.Must(template.ParseFS(folder, "public/email.html"))

	return
}
