package mem

import (
	"fmt"

	"kamilachyla.com/go/dent/data"
)

var patients = InMemoryPatients{
	Patients: map[int]data.Patient{
		1: data.NewPatient("Ala", "Naa", data.Date("2001-05-30")),
		2: data.NewPatient("Fuja", "Goo", data.Date("2021-09-21")),
		3: data.NewPatient("JHlajs dlkajs ", "A salkj sd ", data.Date("1978-11-02")),
	},
	NextId: 4,
}

func GetPatientsDao() data.PatientDao {
	return patients
}

type InMemoryPatients struct {
	Patients map[int]data.Patient
	NextId   int
}

/* Gets all patients from database */
func (imp InMemoryPatients) GetAll() []data.Patient {
	var res = make([]data.Patient, 0)
	for id, wp := range imp.Patients {
		wp.Id = id
		fmt.Println(wp)
		res = append(res, wp)
	}
	fmt.Println(res)
	return res
}

/* Adds data.Parient to in-mem db*/
func (imp InMemoryPatients) Add(p data.Patient) (id int, ok bool) {
	ok = false
	id = imp.NextId
	p.Id = id
	imp.Patients[id] = p
	imp.NextId = id + 1
	ok = true

	return
}

/* Removes data.Patient with given id to in-mem db */
func (imp InMemoryPatients) Remove(id int) (ok bool) {
	delete(imp.Patients, id)
	ok = true
	return
}
