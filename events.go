package main

import (
	"database/sql"
	"time"
)

// Event struct
type Event struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Date     time.Time `json:"date"`
	Venue    Venue     `json:"venue"`
	Sponsors []Sponsor `json:"sponsors"`
	Talks    []Talk    `json:"talks"`
}

// Venue struct
type Venue struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Capacity int    `json:"capacity"`
}

// Sponsor struct
type Sponsor struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Contribution int    `json:"contribution"`
}

// Talk struct
type Talk struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Speaker Speaker `json:"speaker"`
}

// Speaker struct
type Speaker struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Twitter string `json:"twitter"`
	Github  string `json:"github"`
	Bio     string `json:"bio"`
	Avatar  string `json:"avatar"`
}

// sqlStatement := `
// INSERT INTO users (age, email, first_name, last_name)
// VALUES ($1, $2, $3, $4)
// RETURNING id`
// id := 0
// err = db.QueryRow(sqlStatement, 30, "jon@calhoun.io", "Jonathan", "Calhoun").Scan(&id)
// if err != nil {
//   panic(err)
// }
// fmt.Println("New record ID is:", id)

func (e *Event) getEvent(db *sql.DB) error {
	return db.QueryRow("SELECT name, date FROM events WHERE id=$1",
		e.ID).Scan(&e.Name, &e.Date)
}

func getEvents(db *sql.DB, start, count int) ([]Event, error) {
	rows, err := db.Query(
		"SELECT id, name FROM events LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []Event{}

	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.Name); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (e *Event) createEvent(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO events(name, date) VALUES($1, $2) RETURNING id",
		e.Name, e.Date).Scan(&e.ID)

	if err != nil {
		return err
	}
	if e.Venue.ID > 0 {
		_, err := db.Exec(
			"INSERT INTO events_rel(events_id, rel_id, type) VALUES($1, $2, $3)",
			e.ID, e.Venue.ID, "v")
		if err != nil {
			return err
		}

	}

	return nil
}

func (e *Event) updateEvent(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE events SET name=$1, date=$2 WHERE id=$3",
			e.Name, e.Date, e.ID)

	return err
}

func (e *Event) deleteEvent(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM events WHERE id=$1", e.ID)

	return err
}
