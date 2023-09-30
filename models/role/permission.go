package role

import (
	"strings"

	"golang.org/x/exp/slices"
)

func MatchPermission(required, match string) bool {
	if required == "" {
		return true
	}
	if match == "" {
		return false
	}

	if !strings.HasSuffix(required, ".") {
		required = required + "."
	}
	if !strings.HasSuffix(match, ".") {
		match = match + "."
	}

	_required := strings.Split(required, ".")[0]
	_match := strings.Split(match, ".")[0]

	if _match == "*" {
		return true
	}

	if _required == _match {
		return MatchPermission(required[len(_required)+1:], match[len(_match)+1:])
	}

	return false
}

func (role Role) HasPermission(permission string, loaded_roles ...string) bool {
	loaded_roles = append(loaded_roles, role.Name)
	for _, role_perm := range GetRolePermissions(role.Id) {
		if MatchPermission(permission, role_perm) {
			return true
		}
		if strings.HasPrefix(role_perm, ROLE_PREFIX) {
			permission := role_perm[len(ROLE_PREFIX):]
			if !slices.Contains(loaded_roles, permission) {
				_role := GetRoleByName(permission)
				if _role.Name != "" && _role.HasPermission(permission, loaded_roles...) {
					return true
				}
			}
		}
	}
	return false
}
