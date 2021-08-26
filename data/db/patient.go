package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"kamilachyla.com/go/dent/data"
)

/* Implemets PatientDao */
type SQLitePatientDao struct {
	Db *sql.DB
}

/* Returns an instance of PatientsDao that connects to
* a database given in commandline*/
func GetPatientsDao() (dao data.PatientDao, err error) {
	if err := GetDatabase(); err == nil {
		dao = SQLitePatientDao{db}
	}
	return dao, err
}

const (
	QUERY_GET_PATIENT_BY_ID       = "select p.id, p.name, p.surname, p.birthdate, n.id, n.text from patient p left join note n on p.note_id = n.id where p.id=? and p.deleted=false"
	QUERY_DELETE_PATIENT          = "update patient set deleted=true where id=?"
	QUERY_INSERT_PATIENT          = "insert into patient(name, surname, birthdate, note_id) values (?, ?, ?, ?)"
	QUERY_UPDATE_NOTE             = "update note set text = ? where id = ?"
	QUERY_INSERT_NOTE             = "insert into note(text) values (?)"
	QUERY_UPDATE_PATIENT          = "update patient set name = ?, surname = ?, birthdate = ?, note_id = ? where id = ?"
	QUERY_ALL_PATIENTS_WITH_NOTES = "select patient.id, name, surname, birthdate, note.id, note.text from patient left join note on patient.note_id=note.id where patient.deleted=false order by patient.surname"
)

func (dao SQLitePatientDao) GetById(id int) (pat *data.Patient) {

	db := dao.Db
	row, err := db.Query(QUERY_GET_PATIENT_BY_ID, id)
	checkErr(err)
	for row.Next() {
		pat = readPatient(row)
	}
	log.Printf("GetById [%d]: query %s\nresult: %v\n", id, QUERY_GET_PATIENT_BY_ID, pat)
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

/* Adds patient with note to database. Note is created even if empty*/
func (dao SQLitePatientDao) Add(e data.Patient) (id int, er error) {
	log.Printf("Adding: %v\n", e)
	tx, err := dao.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			id, er = -1, r.(error)
			log.Println("Rollback adding patient: ", er.Error())
			checkErr(tx.Rollback())
		} else {
			log.Println("Commit adding patient: ", e)
			checkErr(tx.Commit())
		}
	}()
	checkErr(err)
	var note_id int64
	res, err := tx.Exec("insert into note(text) values (?)", e.Note)
	checkErr(err)

	note_id, err = res.LastInsertId()
	checkErr(err)
	e.NoteId = int(note_id)
	log.Printf("Addeda note [%d] %s", note_id, e.Note)

	res, err = tx.Exec(QUERY_INSERT_PATIENT, e.Name, e.Surname, e.Birthdate, e.NoteId)
	checkErr(err)

	pid, err := res.LastInsertId()
	e.Id = int(pid)
	log.Printf("Addeda patient [%d] %s", pid, e.String())
	checkErr(err)

	return e.Id, nil
}

func (dao SQLitePatientDao) Remove(id int) (ok bool) {
	return Remove(dao.Db, id)
}

func (dao SQLitePatientDao) UpdatePatient(e *data.Patient) (er error) {
	log.Printf("FUNCTION UpdatePatient GOT: %#v\n", e)

	tx, err := dao.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			er = r.(error)
			log.Println("UpdatePatient: recover != nill, rollback")
			log.Printf("Reason: %s\n", er.Error())
			checkErr(tx.Rollback())
		} else {
			log.Println("UpdatePatient: recover == nill, commit")
			checkErr(tx.Commit())
		}
	}()

	var res sql.Result
	if e.NoteId >= 0 {
		res, err = tx.Exec(QUERY_UPDATE_NOTE, e.Note, e.NoteId)
		checkErr(err)

	} else {
		e.NoteId = -1
		res, err = tx.Exec(QUERY_INSERT_NOTE, e.Note)
		checkErr(err)

		var lastId int64
		lastId, err = res.LastInsertId()
		checkErr(err)
		e.NoteId = int(lastId)
	}

	log.Printf("Issuing update patient sql with note_id %d and patient=%v\n", e.NoteId, e)
	_, err = tx.Exec(QUERY_UPDATE_PATIENT, e.Name, e.Surname, e.Birthdate, e.NoteId, e.Id)
	checkErr(err)
	return nil
}

func Remove(db *sql.DB, id int) (ok bool) {
	stmt, err := db.Prepare(QUERY_DELETE_PATIENT)
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
	rows, err := db.Query(QUERY_ALL_PATIENTS_WITH_NOTES)
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
