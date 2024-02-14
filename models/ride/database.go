package ride

import (
	"oui/models/postgresql"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/kataras/golog"
)

func (ride *Ride) UpsertPgSQL() (success bool) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	if ride.ID == 0 {
		ride.ID = time.Now().UnixNano()
	}

	query := "INSERT INTO ride (id, operator, taxi, clientName, clientNumber, latitudeStart, longitudeStart, latitudeEnd, longitudeEnd) " +
		"VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) " +
		"ON CONFLICT (id) do " +
		"UPDATE set (operator, taxi, clientName, clientNumber, latitudeStart, longitudeStart, latitudeEnd, longitudeEnd) = " +
		"($2,$3,$4,$5,$6,$7,$8,$9)"

	args := []interface{}{
		ride.ID,
		ride.Operator,
		ride.Taxi,
		ride.ClientName,
		ride.ClientNumber,
		ride.Start.Latitude,
		ride.Start.Longitude,
		ride.End.Latitude,
		ride.Start.Longitude,
	}

	_, err = sqlCo.Exec(postgresql.SQLCtx, query, args...)
	if err != nil {
		golog.Error(query, err)
	} else {
		success = true
	}

	return
}
