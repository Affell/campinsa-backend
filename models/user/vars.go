package user

const PERMISSIONS_EDIT_PERMISSION = "edit.user.permission"

var SensibleFields []string = []string{
	"email",
	"password",
}

type User struct {
	ID         int64  `structs:"id"`
	Firstname  string `structs:"firstname"`
	Lastname   string `structs:"lastname"`
	Email      string `structs:"email"`
	Password   string `structs:"password"`
	ResetToken string `structs:"reset_token"`
}
