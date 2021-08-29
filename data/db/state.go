package db

import (
	"database/sql"
	"log"

	"kamilachyla.com/go/dent/data"
)

type SQLiteChangeDao struct {
	Db *sql.DB
}

func GetChangeDao() (dao data.ChangeDao, err error) {
	if err := GetDatabase(); err == nil {
		dao = SQLiteChangeDao{db}
	}
	return
}

const (
	QUERY_FOR_VISIT = `
select c.id, visit_id, patient_id, state_id, tooth_num, tooth_side
from change c 
left join visit v on c.visit_id = v.id
where visit_id = ?
`
	QUERY_FOR_ALL_REVERSED = `
select c.id, visit_id, patient_id, state_id, tooth_num, tooth_side
from change c
left join visit v on c.visit_id = v.id
left join patient p on p.id = v.patient_id
where vdatetime <= (select vdatetime from visit where visit.id = ?)
`
)

func readChange(rows *sql.Rows) (change *data.Change) {
	var id int
	var visitid int
	var patientid int
	var stateid int
	var tn int
	var ts int
	err := rows.Scan(&id, &visitid, &patientid, &stateid, &tn, &ts)
	checkErr(err)
	change = data.NewChange(id, visitid, patientid, stateid, tn, ts)
	log.Printf("state.go: ReadChange from db: %s ", change)
	defer func() {
		if r := recover(); r != nil {
			change = nil
		}
	}()
	return
}

func readRows(db *sql.DB, q string, visitId int) (changes []data.Change, err error) {
	rows, err := db.Query(q, visitId)
	checkErr(err)
	defer rows.Close()

	changes = make([]data.Change, 0)
	for rows.Next() {
		if ch := readChange(rows); ch != nil {
			changes = append(changes, *ch)
		}
	}
	return
}

func (d SQLiteChangeDao) ForVisit(visitId int) (states []data.Change, err error) {
	return readRows(d.Db, QUERY_FOR_VISIT, visitId)

}
func (d SQLiteChangeDao) AllReversed(visitId int) (states []data.Change, err error) {

	log.Printf("Called state AllReversed with visitId=%d\n", visitId)
	return readRows(d.Db, QUERY_FOR_ALL_REVERSED, visitId)
}
