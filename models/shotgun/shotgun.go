package shotgun

import (
	"encoding/json"
	"oui/models/postgresql"

	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/kataras/golog"
)

type Shotgun struct {
	Id         string `json:"id" structs:"id"`
	FormLink   string `json:"form_link" structs:"-"`
	ImageBytes string `json:"image_bytes" structs:"image_bytes"`
	Name       string `json:"name" structs:"name"`
	UnlockTime int64  `json:"unlock_time" structs:"unlock_time"`
}

func (shotgun Shotgun) ToWebDetails() map[string]interface{} {
	return structs.Map(shotgun)
}

func (shotgun *Shotgun) StoreRedis() bool {
	data, err := json.Marshal(shotgun)
	if err != nil {
		return false
	}
	return RedisConn.Set(RedisCtx, shotgun.Id, data, 0).Err() == nil
}

func (shotgun *Shotgun) StorePgSQL() (success bool) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	if shotgun.Id == "" {
		shotgun.Id = uuid.New().String()
	}

	query := "INSERT INTO shotgun (id, unlock_time, form_link, image_link, name) " +
		"VALUES ($1,$2,$3,$4,$5) " +
		"ON CONFLICT (id) do " +
		"UPDATE set (unlock_time, form_link, image_link, name) = " +
		"($2,$3,$4,$5)"

	args := []interface{}{
		shotgun.Id,
		shotgun.UnlockTime,
		shotgun.FormLink,
		shotgun.ImageBytes,
		shotgun.Name,
	}

	_, err = sqlCo.Exec(postgresql.SQLCtx, query, args...)
	if err != nil {
		golog.Error(query, err)
	} else {
		success = true
	}

	return
}

func GetAllShotguns() (shotguns []Shotgun, err error) {

	keys, err := RedisConn.Do(RedisCtx, "KEYS", "*").StringSlice()
	if err != nil {
		return
	}

	for _, key := range keys {
		var (
			data []byte
			s    Shotgun
		)

		data, err = RedisConn.Get(RedisCtx, key).Bytes()
		if err != nil {
			continue
		}

		err = json.Unmarshal(data, &s)
		if err != nil {
			continue
		}

		if s.Id != "" {
			shotguns = append(shotguns, s)
		}
	}

	return
}

func GetShotgun(id string) (shotgun Shotgun, err error) {

	data, err := RedisConn.Get(RedisCtx, id).Bytes()
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &shotgun)
	if err != nil {
		return
	}

	return
}
