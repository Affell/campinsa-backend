package shotgun

import (
	"database/sql"
	"oui/models/postgresql"

	"github.com/jackc/pgx/v4"
	"github.com/kataras/golog"
)

func LoadShotgunsIntoCache() {
	query := "SELECT id,created_time,unlock_time,form_link,image_link,name,description,ended FROM shotgun"

	conn, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}

	defer conn.Close(postgresql.SQLCtx)

	rows, err := conn.Query(postgresql.SQLCtx, query)
	if err == pgx.ErrNoRows {
		return
	} else if err != nil {
		golog.Errorf("query '%s' return error : %s", query, err)
		return
	}

	for rows.Next() {
		var id, formLink, imageLink, name, description sql.NullString
		var createdTime, unlockTime sql.NullInt64
		var ended bool

		rows.Scan(&id, &createdTime, &unlockTime, &formLink, &imageLink, &name, &description, &ended)

		shotgun := Shotgun{
			Id:          id.String,
			CreatedTime: createdTime.Int64,
			UnlockTime:  unlockTime.Int64,
			FormLink:    formLink.String,
			ImageBytes:  imageLink.String,
			Name:        name.String,
			Description: description.String,
			Ended:       ended,
		}
		if !shotgun.StoreRedis() {
			golog.Errorf("failed to store shotgun %v : '%s'", id, name)
		}
	}
}
