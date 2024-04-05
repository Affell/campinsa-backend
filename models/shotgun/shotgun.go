package shotgun

import (
	"encoding/json"
	"oui/models/postgresql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/kataras/golog"
)

type Shotgun struct {
	Id          string `json:"id"`
	FormLink    string `json:"form_link"`
	ImageBytes  string `json:"image_bytes"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UnlockTime  int64  `json:"unlock_time"`
	Ended       bool   `json:"ended"`
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

	query := "INSERT INTO shotgun (id, unlock_time, form_link, image_link, name, description, ended) " +
		"VALUES ($1,$2,$3,$4,$5,$6,$7) " +
		"ON CONFLICT (id) do " +
		"UPDATE set (unlock_time, form_link, image_link, name, description, ended) = " +
		"($2,$3,$4,$5,$6,$7)"

	args := []interface{}{
		shotgun.Id,
		shotgun.UnlockTime,
		shotgun.FormLink,
		shotgun.ImageBytes,
		shotgun.Name,
		shotgun.Description,
		shotgun.Ended,
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
