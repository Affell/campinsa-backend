package membre

import (
	"database/sql"
	"oui/models/postgresql"

	"github.com/jackc/pgx/v4"
)

func ScanMember(row pgx.Row) (member Member) {
	var (
		id, score                                        sql.NullInt64
		firstname, lastname, surname, image, poste, pays sql.NullString
	)

	err := row.Scan(&id, &firstname, &lastname, &surname, &score, &image, &poste, &pays)
	if err != nil {
		return Member{Id: -1}
	}

	return Member{
		Id:        id.Int64,
		Firstname: firstname.String,
		Lastname:  lastname.String,
		Surname:   surname.String,
		Score:     int(score.Int64),
		Image:     image.String,
		Poste:     poste.String,
		Pays:      pays.String,
	}
}

func GetAllMembers() (members []Member) {
	conn, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer conn.Close(postgresql.SQLCtx)

	query := "SELECT id, firstname, lastname, surname, score, image, poste, pays FROM membre"
	rows, err := conn.Query(postgresql.SQLCtx, query)
	if err != nil {
		return
	}

	for rows.Next() {
		m := ScanMember(rows)
		if m.Id != -1 {
			members = append(members, m)
		}
	}

	return
}

func GetRandomMember() (member Member) {
	conn, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return Member{Id: -1}
	}
	defer conn.Close(postgresql.SQLCtx)

	query := "SELECT id, firstname, lastname, surname, score, image, poste, pays FROM membre ORDER BY random() LIMIT 1"
	row := conn.QueryRow(postgresql.SQLCtx, query)
	if err != nil {
		return Member{Id: -1}
	}

	return ScanMember(row)
}
