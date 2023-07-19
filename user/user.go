package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"timezone-converter/db"

	"github.com/gorilla/mux"
)

type User struct {
	Id        string                    `json:"id"`
	Username  string                    `json:"username"`
	Password  string                    `json:"password"`
	Timeslots map[string]TimeslotStatus `json:"timeslots"`
}

type TimeslotStatus struct {
	Booked bool `json:"booked"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB{db: db.DbInstance}
	// get posted data for user and store into db
	var u User
	defaultTimeslots := map[string]TimeslotStatus{
		"9:00":  {Booked: false},
		"10:00": {Booked: false},
		"11:00": {Booked: false},
		"12:00": {Booked: false},
		"13:00": {Booked: false},
		"14:00": {Booked: false},
		"15:00": {Booked: false},
		"16:00": {Booked: false},
	}
	u.Timeslots = defaultTimeslots
	jsonWithTimeslots, err := json.Marshal(defaultTimeslots)

	if err != nil {
		log.Panic(err)
	}

	json.NewDecoder(r.Body).Decode(&u)

	db.exec("INSERT INTO users(username, password, timeslots) values(?,?,?)",
		u.Username, u.Password, jsonWithTimeslots,
	)

	json.NewEncoder(w).Encode(u)
}

func BookTime(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB{db: db.DbInstance}

	params := mux.Vars(r)
	userId := params["userId"]
	user := db.queryRow("SELECT * FROM users WHERE id=?", userId)

	// for now
	fmt.Println(user)

	// TODO: Timeslots seems to be empty there, figure out why
}
