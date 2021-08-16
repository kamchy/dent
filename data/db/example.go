package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"kamilachyla.com/go/dent/data"
)

/* Implemets PatientDao */
type SQLitePatientDao struct {
	Db *sql.DB
}

func GetPatientsDao() data.PatientDao {
	db, err := sql.Open("sqlite3", os.Args[1])

	checkErr(err)
	return SQLitePatientDao{db}
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

func (dao SQLitePatientDao) Add(e data.Patient) (id int, ok bool) {
	id64, isok := AddPerson(dao.Db, e)
	if isok {
		return int(id64), isok
	}
	return -1, false
}

func (dao SQLitePatientDao) Remove(id int) (ok bool) {
	return Remove(dao.Db, id)
}

func Remove(db *sql.DB, id int) (ok bool) {
	fmt.Printf("DELETE PERSON: %d----\v", id)
	stmt, err := db.Prepare("delete from patient where id=?")
	checkErr(err)
	_, err = stmt.Exec(id)
	ok = err == nil
	checkErr(err)
	defer stmt.Close()
	return
}

func AddPerson(db *sql.DB, p data.Patient) (id int64, ok bool) {
	fmt.Printf("ADD PERSON: %v----\v", p)
	stmt, err := db.Prepare("insert into patient(name, surname, birthdate) values  (?, ?, ?)")
	checkErr(err)
	var res sql.Result
	res, err = stmt.Exec(p.Name, p.Surname, p.Birthdate)
	checkErr(err)
	id, err = res.LastInsertId()
	ok = err == nil
	checkErr(err)
	defer stmt.Close()
	return
}

func ReadPersons(db *sql.DB, cons func(p data.Patient)) {
	fmt.Println("READ PERSONS: ----")
	rows, err := db.Query("select * from patient")
	checkErr(err)
	var name string
	var surname string
	var id int
	var birthdate time.Time
	var noteid sql.NullInt64
	for rows.Next() {
		err = rows.Scan(&id, &name, &surname, &birthdate, &noteid)
		checkErr(err)
		var noteid_int int64 = -1
		if noteid.Valid {
			noteid_int = noteid.Int64
		}
		p := data.NewPatient(name, surname, birthdate)
		p.NoteId = int(noteid_int)
		p.Id = id
		cons(p)
	}
	rows.Close()
}

func MainDb() {
	db, err := sql.Open("sqlite3", os.Args[1])

	checkErr(err)
	defer db.Close()

	ReadPersons(db, PrintPerson)
	kasiaid, ok := AddPerson(db, data.NewPatient("Kasia", "Mściborska", data.Date("2012-03-03")))
	ReadPersons(db, PrintPerson)
	if ok {
		Remove(db, int(kasiaid))
		ReadPersons(db, PrintPerson)
	}
}

func PrintPerson(p data.Patient) {
	fmt.Printf("%v (link: %s)\n", p, p.GetLink())
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
