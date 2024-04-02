package planning

var Days = [...]string{"Lundi", "Mardi", "Mercredi", "Jeudi", "Vendredi"}
var GlobalPlanning Planning

type (
	Planning         map[int64]UserWeekPlanning
	UserWeekPlanning map[int]UserDayPlanning
	UserDayPlanning  []TimeSlot
	TimeSlot         struct {
		Time       string `json:"time"`
		Title      string `json:"title"`
		Additional string `json:"additional"`
	}
	SheetSlot struct {
		Start   int
		End     int
		Content string
	}
)
