package ride

var Riders map[int64]Ride
var GolfetteStatus bool

type LatLng struct {
	Latitude  float64 `json:"latitude" structs:"latitude"`
	Longitude float64 `json:"longitude" structs:"longitude"`
	Name      string  `json:"name" structs:"name"`
}

type Ride struct {
	ID           int64  `json:"id" structs:"id"`
	Operator     int64  `json:"operator" structs:"operator"`
	Taxi         int64  `json:"taxi" structs:"taxi"`
	Completed    bool   `json:"completed" structs:"completed"`
	ClientName   string `json:"client_name" structs:"client_name"`
	ClientNumber string `json:"client_number" structs:"client_number"`
	Start        LatLng `json:"start" structs:"start"`
	End          LatLng `json:"end" structs:"end"`
	Task         string `json:"task" structs:"task"`
	OperatorName string `json:"operator_name" structs:"operator_name"`
	TaxiName     string `json:"taxi_name" structs:"taxi_name"`
	Date         int64  `json:"date" structs:"date"`
}
