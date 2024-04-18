package ride

import (
	"database/sql"
	"oui/models/postgresql"
	"oui/models/user"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/kataras/golog"
)

func (ride *Ride) TranslateRideIds() {
	var ids []interface{}
	ids = append(ids, ride.Operator)
	if ride.Taxi != 0 {
		ids = append(ids, ride.Taxi)
	}

	translations, err := user.GetFullNameById(ids...)
	if err != nil {
		return
	}
	if operator, ok := translations[ride.Operator]; ok {
		ride.OperatorName = operator
	}
	if taxi, ok := translations[ride.Taxi]; ok {
		ride.TaxiName = taxi
	}
}

func (ride *Ride) UpsertPgSQL() (success bool) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	if ride.ID == 0 {
		ride.ID = time.Now().UnixMilli()
	}
	query := "INSERT INTO ride (id, operator, taxi, completed, clientName, clientNumber, startName, latitudeStart, longitudeStart, endName, latitudeEnd, longitudeEnd, task) " +
		"VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) " +
		"ON CONFLICT (id) do " +
		"UPDATE set (operator, taxi, completed, clientName, clientNumber, startName, latitudeStart, longitudeStart, endName, latitudeEnd, longitudeEnd, task) = " +
		"($2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)"

	args := []interface{}{
		ride.ID,
		ride.Operator,
		ride.Taxi,
		ride.Completed,
		ride.ClientName,
		ride.ClientNumber,
		ride.Start.Name,
		ride.Start.Latitude,
		ride.Start.Longitude,
		ride.End.Name,
		ride.End.Latitude,
		ride.End.Longitude,
		ride.Task,
	}

	_, err = sqlCo.Exec(postgresql.SQLCtx, query, args...)
	if err != nil {
		golog.Error(query, err)
		return
	}

	return true
}

func ScanRide(row pgx.Row) (ride Ride) {
	var (
		id, operator, taxi, date                                 sql.NullInt64
		completed                                                sql.NullBool
		clientName, clientNumber, startName, endName, task       sql.NullString
		latitudeStart, longitudeStart, latitudeEnd, longitudeEnd sql.NullFloat64
	)

	err := row.Scan(
		&id,
		&operator,
		&taxi,
		&completed,
		&clientName,
		&clientNumber,
		&startName,
		&latitudeStart,
		&longitudeStart,
		&endName,
		&latitudeEnd,
		&longitudeEnd,
		&task,
		&date,
	)
	if err != nil {
		return Ride{}
	}

	return Ride{
		ID:           id.Int64,
		Operator:     operator.Int64,
		Taxi:         taxi.Int64,
		Completed:    completed.Bool,
		ClientName:   clientName.String,
		ClientNumber: clientNumber.String,
		Start: LatLng{
			Latitude:  latitudeStart.Float64,
			Longitude: longitudeStart.Float64,
			Name:      startName.String,
		},
		End: LatLng{
			Latitude:  latitudeEnd.Float64,
			Longitude: longitudeEnd.Float64,
			Name:      endName.String,
		},
		Task: task.String,
		Date: date.Int64,
	}
}

func GetAllRides(bypassDate bool) (rides []Ride) {
	conn, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer conn.Close(postgresql.SQLCtx)

	now := time.Now().UnixMilli()
	query := "SELECT id, operator, taxi, completed, clientName, clientNumber, startName, latitudeStart, longitudeStart, endName, latitudeEnd, longitudeEnd, task, date FROM ride"
	rows, err := conn.Query(postgresql.SQLCtx, query)
	if err != nil {
		return
	}
	for rows.Next() {
		r := ScanRide(rows)
		if r.ID != 0 {
			if r.Date == 0 || bypassDate || r.Date <= now {
				r.TranslateRideIds()
				rides = append(rides, r)
			}
		}
	}

	return
}

func GetRideByID(id int64) (ride Ride, err error) {

	conn, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer conn.Close(postgresql.SQLCtx)

	query := "SELECT id, operator, taxi, completed, clientName, clientNumber, startName, latitudeStart, longitudeStart, endName, latitudeEnd, longitudeEnd, task, date FROM ride WHERE id=$1"
	row := conn.QueryRow(postgresql.SQLCtx, query, id)
	ride = ScanRide(row)
	ride.TranslateRideIds()

	return
}

func LoadRiders() {

	Riders = make(map[int64]Ride)

	conn, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer conn.Close(postgresql.SQLCtx)

	query := "SELECT id, operator, taxi, completed, clientName, clientNumber, startName, latitudeStart, longitudeStart, endName, latitudeEnd, longitudeEnd, task, date FROM ride WHERE NOT completed AND taxi IS NOT NULL"
	rows, err := conn.Query(postgresql.SQLCtx, query)
	if err != nil {
		return
	}

	for rows.Next() {
		r := ScanRide(rows)
		if r.ID != 0 {
			r.TranslateRideIds()
			Riders[r.Taxi] = r
		}
	}

}
