package models

type Route struct {
	Service   string   `param:"service"`
	Primary   string   `param:"primary"`
	Secondary string   `param:"secondary"`
	Tertiary  string   `param:"tertiary"`
	Tail      []string `param:"tail"`
}
