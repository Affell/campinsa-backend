package shotgun

import (
	"database/sql"
	"oui/models/postgresql"

	"github.com/jackc/pgx/v4"
	"github.com/kataras/golog"
)

func LoadShotgunsIntoCache() {
	query := "SELECT id,unlock_time,form_link,image_link,name FROM shotgun"

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
		var id, formLink, imageLink, name sql.NullString
		var unlockTime sql.NullInt64

		rows.Scan(&id, &unlockTime, &formLink, &imageLink, &name)

		shotgun := Shotgun{
			Id:         id.String,
			UnlockTime: unlockTime.Int64,
			FormLink:   formLink.String,
			ImageBytes: imageLink.String,
			Name:       name.String,
		}
		if !shotgun.StoreRedis() {
			golog.Errorf("failed to store shotgun %v : '%s'", id, name)
		}
	}
}
