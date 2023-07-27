package user

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Timeslots string `json:"timeslots"`
}

type TimeslotStatus struct {
	Booked bool `json:"booked"`
}
