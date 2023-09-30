package user

import (
	"strings"

	"oui/models/role"

	"github.com/fatih/structs"
	"golang.org/x/exp/slices"
)

func (user User) ToSelfWebDetail() map[string]interface{} {
	m := structs.Map(user)
	delete(m, "password")
	delete(m, "reset_token")
	return m
}

func HasPermission(id int64, permission string, loaded_roles ...string) bool {
	for _, user_perm := range GetUserPermissions(id) {
		if role.MatchPermission(permission, user_perm) {
			return true
		}
		if strings.HasPrefix(user_perm, role.ROLE_PREFIX) {
			role_name := user_perm[len(role.ROLE_PREFIX):]
			if !slices.Contains(loaded_roles, role_name) {
				_role := role.GetRoleByName(role_name)
				if _role.Name != "" && _role.HasPermission(permission, loaded_roles...) {
					return true
				}
			}
		}
	}
	return false
}
