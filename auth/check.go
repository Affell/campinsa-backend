package auth

import (
	"net/mail"
	"regexp"
	"strings"
)

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// function to check the password complexity
// if err == nil {valid password} else {return err}
func ValidPassword(password string) (err string) {

	containUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	containLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	containSpecialChar := strings.ContainsAny(password, " \"!#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")
	containNumeric := strings.ContainsAny(password, "0123456789")

	if len(password) < 8 {
		err = "le mot de passe doit contenir plus de 8 charactères"
	} else if !containLower {
		err = "le mot de passe doit contenir une minuscule"
	} else if !containSpecialChar {
		err = "le mot de passe doit contenir un caractère spécial"
	} else if !containNumeric {
		err = "le mot de passe doit contenir un chiffre"
	} else if !containUpper {
		err = "le mot de passe doit contenir une majuscule"
	}

	return
}

func ValidUsername(username string) (err string) {
	var re = regexp.MustCompile(`^(?mi)([\w.\-]+)$`)
	if len(username) < 3 || len(username) > 20 || !re.MatchString(username) {
		return "username must only contains alphanumerics and '-' or '_' and must be between 3 and 20 characters long"
	}

	return
}
func ValidBase64(base64 string) bool {
	var re = regexp.MustCompile(`[^-A-Za-z0-9+/=]|=[^=]|={3,}$`)
	return re.MatchString(base64)
}
