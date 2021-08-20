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
	Note      string    `json:"note"`
	Id        int       `json:"id,omitempty"`
}

func NewPatient(id int, name string, surname string, birthdate time.Time, note_id int, note string) Patient {
	return Patient{Name: name, Id: id, Surname: surname, Birthdate: birthdate, Note: note, NoteId: note_id}
}

func (p Patient) GetLink() string {
	return fmt.Sprintf("/patient/%d", p.Id)
}
func (p Patient) BirthString() string {
	return p.Birthdate.Format("2006-01-02")
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

type Note struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

func NewNote(text string) Note {
	return Note{-1, text}
}
