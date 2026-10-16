package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // _ kita butuh import package ini, tapi ga dipake secara langsung, cuma untuk ngasih tau go bahwa driver ini ada
)

var DB *sql.DB // variabel pointer ke database yang akan kita pake di seluruh aplikasi

func InitDB(){ // inisialisasi database
	var err error
	
	DB, err = sql.Open("sqlite3", "api.db") // buka koneksi ke database api.db

	if err != nil{
		panic("Could not connect to database.") // kalo error, stop program
	}
	
	DB.SetMaxOpenConns(10) // berapa koneksi yang di open
	DB.SetMaxIdleConns(5) // berapa koneksi yang di open kalo ga ada yang pake

	createTables() // bikin tabel
}

func createTables(){
	var (
		err error
		createUsersTable string
		createEventsTable string
		createRegistrationsTable string
	)

	// membuat tabel users (jika belum ada)
	createUsersTable = `
	CREATE TABLE IF NOT EXISTS users (
		id 				INTEGER PRIMARY KEY AUTOINCREMENT,
		email 			TEXT NOT NULL UNIQUE,
		password 		TEXT NOT NULL
	);
	`

	// menjalankan query
	_, err = DB.Exec(createUsersTable)
	if err != nil {
		panic("Could not create users table")
	}

	// membuat tabel events (jika belum ada)
	createEventsTable = `
	CREATE TABLE IF NOT EXISTS events (
		id 				INTEGER PRIMARY KEY AUTOINCREMENT,
		name 			TEXT NOT NULL,
		description 	TEXT NOT NULL,
		location 		TEXT NOT NULL,
		dateTime 		DATETIME NOT NULL,
		user_id 		INTEGER,

		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`

	// menjalankan query
	_, err = DB.Exec(createEventsTable)
	if err != nil {
		panic("Could not create events table")
	}
 
	createRegistrationsTable = `
	CREATE TABLE IF NOT EXISTS registrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id INTEGER,
		user_id INTEGER,
		
		FOREIGN KEY(event_id) REFERENCES events(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`

	_, err = DB.Exec(createRegistrationsTable)
	if err != nil {
		panic("Could not create registrations table")
	}
} 