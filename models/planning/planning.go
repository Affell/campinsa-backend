package planning

import (
	"encoding/json"
	"io"
	"os"

	"github.com/kataras/golog"
)

func InitPlanning() {
	updateToken := os.Getenv("UPDATE_TOKEN")
	if updateToken == "" {
		golog.Fatal("No UPDATE_TOKEN found in .env")
	}
	UpdateToken = updateToken

	f, err := os.Open("planning.json")
	if err != nil {
		UpdatePlanning()
	} else {
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return
		}
		json.Unmarshal(b, &GlobalPlanning)
	}
}

func UpdatePlanning() {
	GlobalPlanning = RetrievePlanning()
}
