package role

const (
	ROLE_PREFIX          = "role."
	ROLE_EDIT_PERMISSION = "edit.role.permission"
)

type Role struct {
	Id   int64  `structs:"id"`
	Name string `structs:"name"`
}
