package data

import (
	"fmt"
	"time"
)

type Visit struct {
	Id        int
	VisitDate time.Time
	PatientId int
	NoteId    int
	Note      string
	Name      string
	Surname   string
}

func NewVisit(id int, date time.Time, patientid int, noteid int, note string, name string, surname string) *Visit {
	return &Visit{Id: id, VisitDate: date, PatientId: patientid, NoteId: noteid, Note: note, Name: name, Surname: surname}
}

func NewForPatient(date time.Time, p *Patient) *Visit {
	return &Visit{-1, date, p.Id, -1, "", p.Name, p.Surname}
}

type VisitDao interface {
	GetAll() []Visit
	GetById(id int) *Visit
	GetByPatientId(patientId int) []Visit
	Add(v Visit) (id int, err error)
	UpdateVisit(v Visit) error
}

const VISIT_DATE_LAYOUT = "2006-01-02 15:04"
const VISIT_DATE_ONLY_LAYOUT = "2006-01-02"
const VISIT_TIME_ONLY_LAYOUT = "15:04"

func (v Visit) VisitDateString() string {
	return v.VisitDate.Local().Format(VISIT_DATE_ONLY_LAYOUT)
}

func (v Visit) VisitTimeString() string {
	return v.VisitDate.Local().Format(VISIT_TIME_ONLY_LAYOUT)
}
func (v Visit) GetLink() string {
	return fmt.Sprintf("/patients/%d/visits/%d", v.PatientId, v.Id)
}

func (v Visit) From(d time.Time, t time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
}

func (v Visit) String() string {
	return fmt.Sprintf("VISIT [%d] %v - patient[%d]<%s %s> - note[%d]<%s>\n", v.Id, v.VisitDateString(), v.PatientId, v.Name, v.Surname, v.NoteId, v.Note)
}
