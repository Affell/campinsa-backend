package membre

type Member struct {
	Id        int64  `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Surname   string `json:"surname"`
	Score     int    `json:"score"`
	Image     string `json:"image"`
	Poste     string `json:"poste"`
	Pays      string `json:"pays"`
}
