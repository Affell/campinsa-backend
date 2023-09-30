package role

import (
	"database/sql"
	"fmt"
	"strings"

	"oui/models/postgresql"

	"github.com/jackc/pgx/v4"
	"github.com/kataras/golog"
)

func GetRoleByName(name string) (role Role) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)
	var (
		id sql.NullInt64
	)
	query := "SELECT id FROM role WHERE name=$1"
	err = sqlCo.QueryRow(postgresql.SQLCtx, query, name).Scan(
		&id,
	)
	if err == nil {
		role = Role{
			Id:   id.Int64,
			Name: name,
		}
	}
	return
}

func GetRolePermissions(role int64) (permissions []string) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}

	defer sqlCo.Close(postgresql.SQLCtx)

	query := "SELECT permission FROM role_permission WHERE role=$1"
	rows, err := sqlCo.Query(postgresql.SQLCtx, query, role)
	if err != nil {
		golog.Errorf("execution query '%s':\n%s", query, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err == pgx.ErrNoRows {
			return
		} else if err != nil {
			golog.Errorf("psql scan '%v' failed with error : %v", query, err)
			return
		}
		permissions = append(permissions, permission)
	}
	return
}

func UpsertRole(role *Role) (msg string) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return "Internal server error"
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := "INSERT INTO role (id, name) " +
		"VALUES ($1,$2) ON CONFLICT(id) " +
		"DO UPDATE SET name = $2 RETURNING id"

	err = sqlCo.QueryRow(postgresql.SQLCtx, query,
		role.Id,
		role.Name).Scan(&role.Id)
	if err != nil {
		return "Role name not available"
	}
	return
}

func UpdateRolePermissions(role int64, permissions []string) (msg string) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return "Internal server error"
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	sqlCo.Exec(postgresql.SQLCtx, "DELETE FROM role_permission WHERE role=$1", role)

	if len(permissions) > 0 {
		var values []string
		var args []interface{}

		for i, v := range permissions {
			values = append(values, fmt.Sprintf("($%d,$%d)", 2*i+1, 2*(i+1)))
			args = append(args, role, v)
		}

		query := "INSERT INTO role_permission(role,permission) VALUES " + strings.Join(values, ",")

		_, err := sqlCo.Exec(postgresql.SQLCtx, query, args...)
		if err != nil {
			return "Unable to insert role permissions"
		}
	}
	return
}
