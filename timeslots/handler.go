package timeslots

import (
	"encoding/json"
	"net/http"
	"timezone-converter/db"
)

type TimeslotData struct {
	OwnerId string       `json:"ownerId"`
	Time    []TimePeriod `json:"time"`
}

// time payload format that works
// "time": [{"From": "2014-11-12T11:45:26.371Z", "To": "2017-11-12T11:45:26.371Z"}]
func CreateTimeslots(w http.ResponseWriter, r *http.Request) {
	var timeslotsData TimeslotData
	json.NewDecoder(r.Body).Decode(&timeslotsData)

	repo := NewRepository(db.DbInstance)

	for _, timePeriod := range timeslotsData.Time {
		if !isTimePeriodValid(timePeriod) {
			return
		}
	}

	sortByTimeFrom(timeslotsData.Time)

	for i := 0; i < len(timeslotsData.Time); i++ {
		timeInUnixTo := timeslotsData.Time[i].To.Unix()
		timeInUnixFrom := timeslotsData.Time[i].From.Unix()

		t := TimeslotInDb{TimeslotBase: TimeslotBase{OwnerId: timeslotsData.OwnerId, Booked: false}, TimeFrom: timeInUnixFrom, TimeTo: timeInUnixTo}

		err := repo.createTimeslots(t)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func BookTimeslot(w http.ResponseWriter, r *http.Request) {
	repo := NewRepository(db.DbInstance)

	var data Timeslot
	json.NewDecoder(r.Body).Decode(&data)

	ts, err := repo.getTimeslot(data)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ts.BookedById = data.BookedById

	if ts.Booked {
		w.WriteHeader(http.StatusConflict)
		return
	} else {
		err := repo.bookTimeslot(ts)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
