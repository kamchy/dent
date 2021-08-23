package data

import (
	"fmt"
	"time"
)

type Patient struct {
	Name      string
	Surname   string
	Birthdate time.Time
	NoteId    int
	Note      string
	Id        int
}

func NewPatient(id int, name string, surname string, birthdate time.Time, note_id int, note string) Patient {
	return Patient{Name: name, Id: id, Surname: surname, Birthdate: birthdate, Note: note, NoteId: note_id}
}

func (p Patient) GetLink() string {
	return fmt.Sprintf("/patient/%d", p.Id)
}
func (p Patient) GetVisits() string {
	return fmt.Sprintf("/patients/%d/visits", p.Id)

}
func (p Patient) BirthString() string {
	return p.Birthdate.Format(DATE_LAYOUT)
}

func (p Patient) String() string {
	return fmt.Sprintf("PATIENT[%d] %s %s %s - note [%d]<%s>",
		p.Id, p.Name, p.Surname, p.BirthString(), p.NoteId, p.Note)
}

type PatientDao interface {
	GetAll() (all []Patient)
	Add(e Patient) (id int, err error)
	Remove(id int) (ok bool)
	UpdatePatient(e *Patient) error
	GetById(id int) *Patient
}

const DATE_LAYOUT = "2006-01-02"

func Date(s string) time.Time {
	res, err := time.Parse(DATE_LAYOUT, s)
	if err == nil {
		return res
	}
	return time.Now()
}
