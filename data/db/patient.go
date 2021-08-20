package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"kamilachyla.com/go/dent/data"
)

/* Implemets PatientDao */
type SQLitePatientDao struct {
	Db *sql.DB
}

var db *sql.DB

const DEFAULT_DATABASE = "./database.db"

/* Creates database in given location and runs db/create.sql*/
func createEmptyDb(fname string) (res string, err error) {

	if db, err := sql.Open("sqlite3", fname); err == nil {
		defer func() {
			if r := recover(); r != nil {
				res, err = "", r.(error)
			}
			db.Close()
		}()

		c, ioErr := ioutil.ReadFile("db/create.sql")
		checkErr(ioErr)
		sql := string(c)
		_, err := db.Exec(sql)
		checkErr(err)
		return fname, nil
	}
	return "", err
}

/* Gets the name of database passed in commandline arg or
* the name of newly created and initialized db in location
* pointed to by DEFAULT_DATABASE */
func getDbFileName() (string, error) {
	var fname = DEFAULT_DATABASE
	if len(os.Args) > 1 {
		fname = os.Args[1]
	}
	if _, err := os.Stat(fname); err == nil {
		return fname, nil
	} else if os.IsNotExist(err) {
		return createEmptyDb(fname)
	}
	return "", errors.New(fmt.Sprintf("Status pliku %s nieznany", fname))
}

/* Returns package-local var db *sql.DB - initialized lazily */
func GetDatabase() (db *sql.DB, err error) {
	var dbname string
	if db == nil {
		if dbname, err = getDbFileName(); err != nil {
			return nil, err
		}

		if db, err = sql.Open("sqlite3", dbname); err == nil {
			return db, nil
		} else {
			return nil, err
		}
	}
	return db, nil
}

/* Returns an instance of PatientsDao that connects to
* a database given in commandline*/
func GetPatientsDao() (dao data.PatientDao, err error) {
	if db, err := GetDatabase(); err == nil {
		dao = SQLitePatientDao{db}
	}
	return dao, err
}

const GET_BY_ID = "select p.id, p.name, p.surname, p.birthdate, n.id, n.text from patient p left join note n on p.note_id = n.id where p.id=?"

func (dao SQLitePatientDao) GetById(id int) (pat *data.Patient) {

	db := dao.Db
	row, err := db.Query(GET_BY_ID, id)
	checkErr(err)
	for row.Next() {
		pat = readPatient(row)
	}
	log.Printf("GetById [%d]: query %s\nresult: %v\n", id, GET_BY_ID, pat)
	return
}

func (dao SQLitePatientDao) GetAll() (all []data.Patient) {
	db := dao.Db
	res := make([]data.Patient, 0)
	consumer := func(p data.Patient) {
		res = append(res, p)
	}
	ReadPersons(db, consumer)
	return res
}

func (dao SQLitePatientDao) Add(e data.Patient) (id int, er error) {
	log.Printf("Adding: %v\n", e)
	tx, err := dao.Db.Begin()
	checkErr(err)
	var note_id int64
	if strings.Trim(e.Note, " ") != "" {
		res, err := tx.Exec("insert into note(text) values (?)", e.Note)
		if err != nil {
			tx.Rollback()
			checkErr(err)
		}
		note_id, err = res.LastInsertId()
		checkErr(err)
		log.Printf("Added %v", e)
	}

	res, err := tx.Exec("insert into patient(name, surname, birthdate, note_id) values (?, ?, ?, ?)", e.Name, e.Surname, e.Birthdate, note_id)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	pid, err := res.LastInsertId()
	if err == nil {
		checkErr(tx.Commit())
		return int(pid), nil
	}
	return -1, err

}

func (dao SQLitePatientDao) Remove(id int) (ok bool) {
	return Remove(dao.Db, id)
}

func (dao SQLitePatientDao) UpdatePatient(e *data.Patient) (er error) {
	log.Printf("FUNCTION UpdatePatient GOT: %#v\n", e)

	tx, err := dao.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Println("UpdatePatient: recover != nill, rollback")
			log.Printf("Reason: %s\n", r.(error).Error())
			tx.Rollback()
		} else {
			log.Println("UpdatePatient: recover == nill, commit")
			checkErr(tx.Commit())
		}
	}()

	var res sql.Result
	if e.NoteId >= 0 {
		res, err = tx.Exec("update note set text = ? where id = ?", e.Note, e.NoteId)
		checkErr(err)

	} else {
		e.NoteId = -1
		res, err = tx.Exec("insert into note(text) values (?)", e.Note)
		checkErr(err)

		var lastId int64
		lastId, err = res.LastInsertId()
		checkErr(err)
		e.NoteId = int(lastId)
	}

	log.Printf("Issuing update patient sql with note_id %d and patient=%v\n", e.NoteId, e)
	_, err = tx.Exec("update patient set name = ?, surname = ?, birthdate = ?, note_id = ? where id = ?", e.Name, e.Surname, e.Birthdate, e.NoteId, e.Id)
	checkErr(err)
	return nil
}

func Remove(db *sql.DB, id int) (ok bool) {
	stmt, err := db.Prepare("delete from patient where id=?")
	checkErr(err)
	_, err = stmt.Exec(id)
	ok = err == nil
	checkErr(err)
	defer stmt.Close()
	return
}

func readPatient(scanner *sql.Rows) (p *data.Patient) {

	var name string
	var surname string
	var id int
	var birthdate time.Time
	var noteid sql.NullInt64
	var note sql.NullString

	err := scanner.Scan(&id, &name, &surname, &birthdate, &noteid, &note)
	checkErr(err)

	var noteid_int int64 = -1
	if noteid.Valid {
		noteid_int = noteid.Int64
	}

	var note_val string
	if note.Valid {
		note_val = note.String
	}

	p = &data.Patient{Name: name, Id: id, Surname: surname, Birthdate: birthdate, Note: note_val, NoteId: int(noteid_int)}

	log.Printf("readPatient: %v\n", p)
	defer func() {
		if r := recover(); r != nil {
			p = nil
		}
	}()

	return
}

func ReadPersons(db *sql.DB, cons func(p data.Patient)) {
	rows, err := db.Query("select patient.id, name, surname, birthdate, note.id, note.text from patient left join note on patient.note_id=note.id")
	checkErr(err)

	for rows.Next() {
		p := readPatient(rows)
		if p != nil {
			cons(*p)
		}
	}
	rows.Close()
}

/* Panics if error is not nil */
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
