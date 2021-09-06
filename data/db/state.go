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
with
	pv(id, patient_id, vdatetime) as (
		select id, patient_id, vdatetime
		from visit
		where visit.id = ?)
select
	c.id, c.visit_id, v.patient_id, c.state_id, c.tooth_num, c.tooth_side
from
	change c
left join
	visit v
on
	c.visit_id = v.id
join
	pv
on
	v.patient_id = pv.patient_id
where
	v.vdatetime <= pv.vdatetime
order by
	v.vdatetime asc,
	c.time asc
`
	INSERT_CHANGE = "insert into change(visit_id, state_id, tooth_num, tooth_side) values (?, ?, ?, ?)"

	QUERY_STATES = "select id, name, whole from state"
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

func (d SQLiteChangeDao) InsertChange(data data.Change) (change *data.Change, err error) {

	if res, err := d.Db.Exec(INSERT_CHANGE, data.VisitId, data.StateId, data.ToothNum, data.ToothSide); err == nil {
		if id64, err := res.LastInsertId(); err == nil {
			data.Id = int(id64)
		} else {
			return &data, err
		}
	}
	return &data, err
}

func (d SQLiteChangeDao) GetStates() (states []data.State, err error) {
	rows, err := d.Db.Query(QUERY_STATES)
	if err != nil {
		return nil, err
	}

	var id int
	var name string
	var val bool

	states = make([]data.State, 0)
	for rows.Next() {
		rows.Scan(&id, &name, &val)
		states = append(states, data.State{Id: id, Name: name, Whole: val})
	}
	log.Printf("SQLiteChangeDao: GetStates : %v\n", states)
	return states, nil

}
