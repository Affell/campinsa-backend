package ride

type LatLng struct {
	Latitude  float64 `json:"latitude" structs:"latitude"`
	Longitude float64 `json:"longitude" structs:"longitude"`
}

type Ride struct {
	ID           int64  `json:"id" structs:"id"`
	Operator     int64  `json:"operator" structs:"operator"`
	Taxi         int64  `json:"taxi" structs:"taxi"`
	ClientName   string `json:"client_name" structs:"client_name"`
	ClientNumber string `json:"client_number" structs:"client_number"`
	Start        LatLng `json:"start" structs:"start"`
	End          LatLng `json:"end" structs:"end"`
}