package activity

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"secondlap.go/entry"
)

type ProfRep struct {
	Mu sync.Mutex
	db *sql.DB
}

const file string = "data/profiles.db"

const create string = `
  CREATE TABLE IF NOT EXISTS profiles (
  id INTEGER NOT NULL PRIMARY KEY,
  time DATETIME NOT NULL,
  email TEXT NOT NULL,
  firstname TEXT NOT NULL,
  lastname TEXT NOT NULL,
  entry INTEGER NOT NULL,
  last_entered DATETIME
  );`

func InitDB() (*ProfRep, error) {

	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("error opening db: %s", err)
	}

	if _, err := db.Exec(create); err != nil {
		return nil, err
	}

	return &ProfRep{
		db: db,
	}, nil

}

func (p *ProfRep) insert(email, firstname, lastname string) (int, error) {
	res, err := p.db.Exec("INSERT INTO profiles VALUES(NULL,?,?,?,?,?,?);", time.Now(), email, firstname, lastname, 0, nil)
	if err != nil {
		return -1, err
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return -1, err
	}

	fmt.Println("inserted: ", email)
	return int(id), nil

}

func (p *ProfRep) doesExist(email string) bool {

	query := "SELECT id, firstname, email FROM profiles WHERE email = ?"

	var id int
	var name, foundEmail string
	err := p.db.QueryRow(query, email).Scan(&id, &name, &foundEmail)

	return err != sql.ErrNoRows

	// return true

}

func (p *ProfRep) GetTableSize() (int, error) {

	var count int
	err := p.db.QueryRow("SELECT COUNT (*) FROM profiles").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil

}

func (p *ProfRep) SelectProfiles(count int) []entry.Profile {

	var profs []entry.Profile
	q := "SELECT id, firstname, lastname, email, entry FROM profiles ORDER BY entry ASC LIMIT ?"

	rows, err := p.db.Query(q, count)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id, entrycount int
		var firstname, lastname, email string
		err := rows.Scan(&id, &firstname, &lastname, &email, &entrycount)
		if err != nil {
			log.Fatal(err)
		}

		e := entry.Profile{
			FirstName:  firstname,
			LastName:   lastname,
			Email:      email,
			EntryCount: entrycount,
		}

		profs = append(profs, e)
		fmt.Printf("ID: %d, FirstName: %s, LastName: %s, Email: %s, Entry Count: %d\n", id, firstname, lastname, email, entrycount)

	}

	return profs

}

func (p *ProfRep) UpdateProfile(email string) error {

	if err := p.updateEntryCount(email); err != nil {
		return fmt.Errorf("error updateing entry count: %s", err)
	}

	if err := p.updateEntryTime(email); err != nil {
		return fmt.Errorf("error updating entry time: %s", err)
	}

	return nil

}

// updates the entry count and time last entered given an email
func (p *ProfRep) updateEntryCount(email string) error {

	result, err := p.db.Exec("UPDATE profiles SET last_entered = ? WHERE email = ?", time.Now(), email)
	if err != nil {
		log.Fatal(err)
	}

	aff, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if aff != 0 {
		return nil
	} else {
		return err
	}

}

func (p *ProfRep) updateEntryTime(email string) error {
	result, err := p.db.Exec("UPDATE profiles SET entry = entry + 1 WHERE email = ?", email)
	if err != nil {
		log.Fatal(err)
	}
	aff, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if aff != 0 {
		return nil
	} else {
		return err
	}

}

func (p *ProfRep) ResetEntryCount(email string) bool {

	result, err := p.db.Exec("UPDATE profiles SET entry = 0 WHERE email = ?", email)
	if err != nil {
		log.Fatal(err)
	}

	aff, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if aff != 0 {
		return true
	} else {
		return false
	}

}

func (p *ProfRep) PrintAllProfiles() {

	var entryCount int

	rows, err := p.db.Query("SELECT email, entry, last_entered FROM profiles")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var entry int
		var email string
		var last_entered *time.Time

		err := rows.Scan(&email, &entry, &last_entered)
		if err != nil {
			log.Fatal(err)
		}
		entryCount = entryCount + entry

		if last_entered != nil {
			fmt.Printf("[%s]: %d entries | %s \n", email, entry, formatTime(last_entered))
		} else {
			fmt.Printf("[%s]: %d entries | %s \n", email, entry, "N/A")

		}
	}

	fmt.Printf("Total entries: %d\n", entryCount)

}

func ordinal(n int) string {
	if n >= 11 && n <= 13 {
		return "th"
	}

	switch n % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

func amPm(n int) string {
	if n < 12 {
		return "AM"
	} else {
		return "PM"
	}
}

func formatTime(t *time.Time) string {

	return fmt.Sprintf("%s %d%s, %d:%02d %s", t.Month(), t.Day(), ordinal(t.Day()), t.Hour()%12, t.Minute(), amPm(t.Hour()))

}
