package dbFunctions

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ryanProd/TelegramBot/internal/structs"
)

const dbHost = "DB_HOST"
const dbName = "DB_NAME"
const dbUser = "DB_USER"
const dbPassword = "DB_PASSWORD"
const dbPort = "DB_PORT"

const stateZero string = "Please review our product."

// connect to DB
func ConnectDB() *sql.DB {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s", os.Getenv(dbUser), os.Getenv(dbPassword),
		os.Getenv(dbName), os.Getenv(dbHost))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	var version string
	if err := db.QueryRow("select version()").Scan(&version); err != nil {
		panic(err)
	}

	return db
}

func QueryDB(db *sql.DB, userId int) (*structs.UserState, error) {
	var userState structs.UserState
	selectQuery := fmt.Sprintf(`SELECT * FROM userstates WHERE user_id = %d`, userId)
	log.Printf(selectQuery)
	if err := db.QueryRow(selectQuery).Scan(&userState.UserId,
		&userState.CurrState, &userState.StateZeroReply, &userState.StateOneReply); err != nil {
		if err == sql.ErrNoRows {
			log.Printf(err.Error())
			return nil, err
		}
		return nil, err
	}
	return &userState, nil
}

func InsertRow(db *sql.DB, userId int) {
	insertQuery := fmt.Sprintf(`INSERT INTO userstates (user_id, currState, stateZeroReply, stateOneReply)
	    VALUES ('%d', '%s', '%s', '%s')`, userId, stateZero, "", "")
	log.Printf(insertQuery)
	_, err := db.Exec(insertQuery)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	log.Printf("User ID: %d Inserted", userId)
}

func UpdateRow(db *sql.DB, userId int, currState string, stateZeroReply string, stateOnereply string) {
	updateQuery := fmt.Sprintf(`UPDATE userstates SET currState = '%s', stateZeroReply = '%s', stateOneReply = '%s' 
	    WHERE user_id = '%d'`, currState, stateZeroReply, stateOnereply, userId)
	log.Printf(updateQuery)
	_, err := db.Exec(updateQuery)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	log.Printf("User ID: %d Updated", userId)
}
