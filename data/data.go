package data

import (
	"fmt"
	"time"
)

type Patient struct {
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Birthdate time.Time `json:"birthdate"`
	NoteId    int       `json:"note_id"`
	Id        int       `json:"id,omitempty"`
}

func NewPatient(name string, surname string, birthdate time.Time) Patient {
	return Patient{name, surname, birthdate, -1, -1}
}

func (p Patient) GetLink() string {
	return fmt.Sprintf("/patient/%d", p.Id)
}

type PatientDao interface {
	GetAll() (all []Patient)
	Add(e Patient) (id int, ok bool)
	Remove(id int) (ok bool)
}

const DATE_LAYOUT = "2006-01-02"

func Date(s string) time.Time {
	res, err := time.Parse(DATE_LAYOUT, s)
	if err == nil {
		return res
	}
	return time.Now()
}
