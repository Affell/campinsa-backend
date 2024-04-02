package planning

import (
	"encoding/json"
	"io"
	"os"
)

func InitPlanning() {
	f, err := os.Open("planning.json")
	if err != nil {
		GlobalPlanning = RetrievePlanning()
	} else {
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return
		}
		json.Unmarshal(b, &GlobalPlanning)
	}
}
