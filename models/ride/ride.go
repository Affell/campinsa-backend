package ride

import (
	"github.com/fatih/structs"
)

func Init() {
	LoadRiders()
}

func (ride Ride) ToAppDetails() map[string]interface{} {
	m := structs.Map(ride)
	delete(m, "date")
	return m
}
