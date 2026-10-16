package models

import (
	"database/sql"
	"errors"

	"github.com/jeromedesantos/rest-api/db"
	"github.com/jeromedesantos/rest-api/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	var (
		query string
		stmt *sql.Stmt
		result sql.Result
		id int64
		err error
		hashedPassword string
	)

	query = `
	INSERT INTO users(email, password)
	VALUES (?, ?)`

	stmt, err = db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err = utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err = stmt.Exec(u.Email, hashedPassword) // save hash pwd
	if err != nil {
		return err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = id
	return nil
}

func (u *User) ValidateCredentials() error {
	var (
		query string
		row	*sql.Row
		err error
		retrievedPassword string
		passwordIsValid bool
	)

	// cari user: dari password berdaasarkan email (yg unique)
	query = `SELECT id, password FROM users WHERE email = ?`
	row = db.DB.QueryRow(query, u.Email) 

	err = row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		return errors.New("Credentials invalid")
	}

	// bandingkan apakah password yang di input (plain) sesuai dengan database (hash)
	passwordIsValid = utils.CheckPasswordHash(u.Password, retrievedPassword)
	if !passwordIsValid {
		return errors.New("Credentials invalid")
	}

	return nil
}