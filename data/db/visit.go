package db

import (
	"database/sql"
	"log"
	"time"

	"kamilachyla.com/go/dent/data"
)

type SQLiteVisitDao struct {
	Db *sql.DB
}

func GetVisitsDao() (dao data.VisitDao, err error) {
	if err := GetDatabase(); err == nil {
		log.Printf("GetVisitsDao: db=%v\n", db)
		dao = SQLiteVisitDao{db}
	}
	return
}

const (
	GET_VISITS_NO_COND = ` 
		select 
			v.id, v.vdatetime, v.patient_id, v.note_id, n.text, p.name, p.surname 
		from 
			visit v 
		left join note n on v.note_id=n.id 
		left join patient p on v.patient_id=p.id
		`
	GET_VISITS               = GET_VISITS_NO_COND + " where (v.deleted = false) and (p.deleted = false)"
	GET_VISIT_BY_ID          = GET_VISITS_NO_COND + " where (v.id = ?) and (v.deleted = false)"
	GET_VISITS_BY_PATIENT_ID = GET_VISITS_NO_COND + " where (v.patient_id = ?) and (v.deleted = false) and (p.deleted = false)"
	ADD_VISIT                = `
	insert into visit(vdatetime, patient_id, note_id) 
	values(?, ?, ?)`
	ADD_NOTE     = "insert into note(text) values (?)"
	UPDATE_NOTE  = "update note set text = ? where id = ?"
	UPDATE_VISIT = `
	update visit 
	set vdatetime = ?, patient_id = ?, note_id = ? 
	where id = ?`
	DELETE_VISIT = "update visit set deleted = true where id = ?"
)

func (dao SQLiteVisitDao) GetById(id int) (vis *data.Visit) {
	db := dao.Db
	row, err := db.Query(GET_VISIT_BY_ID, id)
	checkErr(err)
	for row.Next() {
		vis = readVisit(row)
	}
	return
}

func (dao SQLiteVisitDao) GetByPatientId(patientId int) []data.Visit {
	res := make([]data.Visit, 0)
	consumer := func(v data.Visit) {
		res = append(res, v)
	}
	readPatientVisits(db, consumer, patientId)

	return res
}

func (dao SQLiteVisitDao) GetAll() []data.Visit {
	res := make([]data.Visit, 0)
	consumer := func(v data.Visit) {
		res = append(res, v)
	}
	readAllVisits(db, consumer)

	return res
}

/*TODO what id? rolllback with 0*/
func (dao SQLiteVisitDao) Add(v data.Visit) (id int, err error) {
	log.Printf("Adding: %v\n", v)
	tx, err := dao.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			id = -1
			log.Println("Add visit: recover != nill, rollback")
			log.Printf("Reason: %s\n", err.Error())
			tx.Rollback()
		} else {
			log.Println("Add visit: recover == nill, commit")
			checkErr(tx.Commit())
		}
	}()

	checkErr(err)
	var note_id int64
	res, err := tx.Exec(ADD_NOTE, v.Note)
	checkErr(err)

	note_id, err = res.LastInsertId()
	checkErr(err)
	v.NoteId = int(note_id)

	res, err = tx.Exec(ADD_VISIT, v.VisitDate, v.PatientId, v.NoteId)
	checkErr(err)

	pid, err := res.LastInsertId()
	checkErr(err)
	id = int(pid)
	v.Id = id
	return
}

func (dao SQLiteVisitDao) UpdateVisit(v data.Visit) error {
	log.Printf("FUNCTION UpdateVisitGOT: %#v\n", v)

	tx, err := dao.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Println("UpdateVisit: recover != nill, rollback")
			log.Printf("Reason: %s\n", r.(error).Error())
			tx.Rollback()
		} else {
			log.Println("UpdateVisit: recover == nill, commit")
			checkErr(tx.Commit())
		}
	}()

	var res sql.Result
	if v.NoteId >= 0 {
		res, err = tx.Exec(UPDATE_NOTE, v.Note, v.NoteId)
		checkErr(err)

	} else {
		v.NoteId = -1
		res, err = tx.Exec(ADD_NOTE, v.Note)
		checkErr(err)

		var lastId int64
		lastId, err = res.LastInsertId()
		checkErr(err)
		v.NoteId = int(lastId)
	}

	log.Printf("Issuing update visit sql with note_id %d and visit=%v\n", v.NoteId, v)
	_, err = tx.Exec(UPDATE_VISIT, v.VisitDate, v.PatientId, v.NoteId, v.Id)
	checkErr(err)
	return nil
}

func (dao SQLiteVisitDao) Delete(id int) error {
	db := dao.Db
	_, err := db.Exec(DELETE_VISIT, id)
	log.Printf("Deleting visit with id %d, err=%v", id, err)
	checkErr(err)
	return err

}
func readAllVisits(db *sql.DB, cons func(v data.Visit)) {
	rows, err := db.Query(GET_VISITS)
	readVisits(rows, err, cons)
}

func readPatientVisits(db *sql.DB, cons func(v data.Visit), pid int) {
	rows, err := db.Query(GET_VISITS_BY_PATIENT_ID, pid)
	readVisits(rows, err, cons)
}

/* TODO refactor wiht ReadPersons*/
func readVisits(rows *sql.Rows, err error, cons func(v data.Visit)) {
	defer rows.Close()
	checkErr(err)

	for rows.Next() {
		v := readVisit(rows)
		if v != nil {
			cons(*v)
		}
	}
}

func readVisit(scanner *sql.Rows) (v *data.Visit) {
	var id int
	var visitdate time.Time
	var noteid sql.NullInt64
	var note sql.NullString
	var patientid int64
	var name string
	var surname string

	err := scanner.Scan(&id, &visitdate, &patientid, &noteid, &note, &name, &surname)
	checkErr(err)

	var noteid_int int64 = -1
	if noteid.Valid {
		noteid_int = noteid.Int64
	}

	var note_val string
	if note.Valid {
		note_val = note.String
	}

	v = data.NewVisit(id, visitdate, int(patientid), int(noteid_int), note_val, name, surname)

	log.Printf("readVisit: %v\n", v)
	defer func() {
		if r := recover(); r != nil {
			v = nil
		}
	}()

	return
}
