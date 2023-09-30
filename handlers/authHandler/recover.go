package authHandler

import (
	"oui/auth"
	"oui/config"
	"oui/email"
	"oui/models"
	"oui/models/user"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)

func Recover(c iris.Context, route models.Route) {

	if route.Tertiary != "" || len(route.Tail) > 0 {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}

	var (
		code int
		data interface{}
	)

	if route.Secondary == "" {
		var askRecoverForm AskRecoverForm
		if err := c.ReadBody(&askRecoverForm); err != nil || askRecoverForm.Email == "" {
			c.StopWithJSON(iris.StatusBadRequest, iris.Map{"message": "Empty email"})
			return
		}
		code, data = postAskRecover(askRecoverForm)
	} else {
		if len(route.Secondary) != 36 {
			c.StopWithJSON(iris.StatusBadRequest, iris.Map{"message": "Invalid reset token"})
			return
		}

		var recoverForm RecoverForm
		if err := c.ReadBody(&recoverForm); err != nil || recoverForm.Password == "" {
			c.StopWithJSON(iris.StatusBadRequest, iris.Map{"message": "Empty password"})
			return
		}

		code, data = postRecover(route.Secondary, recoverForm)
	}

	if code == 0 {
		code = iris.StatusNotFound
	}

	if data != nil {
		c.StopWithJSON(code, data)
	} else {
		c.StopWithStatus(code)
	}

}

func postAskRecover(askRecoverForm AskRecoverForm) (code int, data interface{}) {
	if !auth.ValidEmail(askRecoverForm.Email) {
		code, data = iris.StatusBadRequest, iris.Map{"message": "Invalid Email format"}
		return
	}
	code, data = iris.StatusOK, iris.Map{"message": "If your email is registered as an account, an email will be sent to initiate the password recovery process"}

	ResetToken := user.GenResetToken(askRecoverForm.Email)
	if ResetToken != "" {
		subject := "Password recovery"
		text := "A password reset request has been made for your Tor'Khan website account using your email address. " +
			"To complete this process, click the button below:"
		btnText := "RESET PASSWORD"
		btnURL := config.Cfg.App.FrontURL + "/auth/recover/" + ResetToken

		sucessfullySent := email.New(askRecoverForm.Email, subject, subject, text, btnText, btnURL).Send(config.Cfg.Email)
		if !sucessfullySent {
			golog.Error("failed send mail to : ", askRecoverForm.Email)
			code, data = iris.StatusInternalServerError, iris.Map{"message": "The password reset email could not be sent. Please contact support at support@torkhan.net"}
		} else {
			golog.Info("Reset email sucessfully sent to : ", askRecoverForm.Email)
		}
	}

	return
}

func postRecover(token string, recoverForm RecoverForm) (code int, data interface{}) {

	if err := auth.ValidPassword(recoverForm.Password); err != "" {
		code, data = iris.StatusBadRequest, iris.Map{"message": err}
		return
	}

	if user.DefinePasswordWithResetToken(token, recoverForm.Password) {
		code = iris.StatusOK
	} else {
		code, data = iris.StatusNotFound, iris.Map{"message": "Invalid token"}
	}

	return
}
