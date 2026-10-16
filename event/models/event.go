package models

import (
	"database/sql"
	"time"

	"github.com/jeromedesantos/rest-api/db"
)

// define shape of data / model
type Event struct {
	ID          int64 // --> sesuaikan dengan tipe data di database
	Name        string `binding:"required"` 
	Description string `binding:"required"`
	Location    string `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
	// required = data harus ada
	// userID, id memang ga perlu diisi manual
}

// save event to db
func (e *Event) Save() error{
	var (
		query string
		stmt  *sql.Stmt
		result sql.Result
		id int64
		err error
	)

	query = `
	INSERT INTO events(name, description, location, dateTime, user_id) VALUES (?, ?, ?, ?, ?)` // ? --> placeholder
	
	// prepare query untuk mencegah SQL injection dan mengoptimalkan performa --> simpan ke memory dulu sebelum dieksekusi
	stmt, err = db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close() // --> memastikan koneksi ditutup
	// ga peduli kelar karena error atau sukses (defer)

	// nampung hasil query
	result, err = stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}
	
	// ambil id terakhir yang di generate di database
	id, err = result.LastInsertId()
	if err != nil {
		return err
	}
	
	e.ID = id // --> simpan id ke struct event
	return err
}

// ga dibuat method karena :
// 1. ga perlu receiver () 
// 2. ga perlu data dari struct event (e *Event)
func GetAllEvents() ([]Event, error) {
	var (
		query string
		err error
	)
	
	// query : cuma ambil data doang (GET)
	// prepare + exec : untuk simpan ke memory lalu eksekusi perintah ubah sesuatu (POST, PUT, PATCH, DELETE)
	query = "SELECT * FROM events"
	rows, err := db.DB.Query(query) // multiple rows perlu di loop
	if err != nil {
		return nil, err
	}

	defer rows.Close() // --> memastikan koneksi ditutup
	// ga peduli kelar karena error atau sukses (defer)

	var events []Event // --> buat wadah buat naruh data hasil query

	// iterasi untuk query (looping semua data)
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.Location, &e.DateTime, &e.UserID)
		if err != nil {
			return nil, err
		}
		
		events = append(events, e)
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	var (
		query string
		row *sql.Row
		err error
	)

	query = "SELECT * FROM events WHERE id = ?"
	row = db.DB.QueryRow(query, id) // single row

	var e Event
	err = row.Scan(&e.ID, &e.Name, &e.Description, &e.Location, &e.DateTime, &e.UserID) // ambil semua data dari db by id
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (e *Event) Update() error{
	var (
		query string
		stmt *sql.Stmt
		err error
	)

	query = `
	UPDATE events
	SET name = ?, description = ?, location = ?, dateTime = ?
	WHERE id = ?`

	stmt, err = db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)
	if err != nil {
		return err
	}

	return nil
}

func (e *Event) Delete() error{
	var (
		query string
		stmt *sql.Stmt
		err error
	)

	query = "DELETE FROM events WHERE id = ?"

	stmt, err = db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID)
	if err != nil {
		return err
	}

	return nil
}

func (e *Event) Register(userId int64) error {
	var (
		query string
		stmt *sql.Stmt
		err error
	)

	query = `INSERT INTO registrations (event_id, user_id) VALUES (?,?)`

	stmt, err = db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)
	return err
}

func (e *Event) CancelRegistration(userId int64) error {
	var (
		query string
		stmt *sql.Stmt
		err error
	)

	query = `DELETE FROM registrations WHERE event_id = ? AND user_id = ?`

	stmt, err = db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)
	return err
}