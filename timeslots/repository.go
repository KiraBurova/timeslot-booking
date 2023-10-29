package timeslots

import (
	"database/sql"
	"log"
	"timezone-converter/db"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func createTimeslotsTable() {
	const create = `CREATE TABLE IF NOT EXISTS timeslots(id TEXT, ownerId TEXT, bookedById TEXT, timeFrom INTEGER, timeTo INTEGER, booked INTEGER)`

	if _, err := db.DbInstance.Exec(create); err != nil {
		log.Fatal(err)
	}
}

func (r Repository) createTimeslot(timeslot Timeslot) error {
	createTimeslotsTable()

	// for i := 0; i < len(timeslots); i++ {
	// add id, TimeFrom, TimeTo here
	//}

	// for _, timeslot := range timeslots {
	// somehow construct SQL query here
	// }

	timeInUnixFrom := timeslot.TimeFrom.Unix()
	timeInUnixTo := timeslot.TimeTo.Unix()

	t := TimeslotInDB{TimeslotBase: TimeslotBase{Id: timeslot.Id, OwnerId: timeslot.OwnerId, Booked: false}, TimeFrom: timeInUnixFrom, TimeTo: timeInUnixTo}

	query := "INSERT INTO timeslots(id, ownerId, bookedById, timeFrom, timeTo, booked) values(?,?,?,?,?,?)"

	_, err := db.DbInstance.Exec(query, t.Id, t.OwnerId, t.BookedById, t.TimeFrom, t.TimeTo, t.Booked)

	if err != nil {
		return err
	}

	return nil
}

func (r Repository) getTimeslotById(id string) (TimeslotInDB, error) {
	timeslot := TimeslotInDB{}
	query := "SELECT * FROM timeslots WHERE id=?"

	row := db.DbInstance.QueryRow(query, id)

	err := row.Scan(&timeslot.Id, &timeslot.OwnerId, &timeslot.BookedById, &timeslot.TimeFrom, &timeslot.TimeTo, &timeslot.Booked)

	if err != nil {
		return timeslot, err
	}

	return timeslot, nil
}

func (r Repository) bookTimeslot(timeslot TimeslotInDB) error {
	query := `UPDATE timeslots SET booked = $1, bookedById = $2 WHERE id=$3 AND ownerId=$4`
	_, err := db.DbInstance.Exec(query, 1, timeslot.BookedById, timeslot.Id, timeslot.OwnerId)

	if err != nil {
		return err
	}

	return nil
}
