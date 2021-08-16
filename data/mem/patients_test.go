package mem

import (
	"testing"

	"kamilachyla.com/go/dent/data"
)

var Patients = GetPatientsDao()

func assertLen(expLen int, ps data.PatientDao, t *testing.T) {
	var got = len(ps.GetAll())
	if got != expLen {
		t.Errorf("Expected %d, got %d", expLen, got)
	}
}

func TestInital(t *testing.T) {
	assertLen(3, Patients, t)
}

func TestFields(t *testing.T) {
	var p = data.NewPatient("kam", "234")
	p.Id = 7

	if p.Name != "kam" {
		t.Errorf("Name should be %s", "kam")
	}
	if p.Surname != "234" {
		t.Errorf("Surname shoud be %d", 234)
	}
	if p.Id != 7 {
		t.Errorf("Number shoud be %d", 7)
	}
}

func TestLink(t *testing.T) {
	var p = data.NewPatient("kam", "234")
	p.Id = 7
	var exp = "/patient/7"
	var got = p.GetLink()
	if got != exp {
		t.Errorf("Link shoud be %s, was %s", got, exp)
	}
}

func TestAdd(t *testing.T) {
	var np = data.NewPatient("kam", "234")
	if npid, ok := Patients.Add(np); ok {
		if npid != 4 {
			t.Errorf("Expected np be 4, is %d", npid)
		}

		assertLen(4, Patients, t)
	} else {
		t.Error("Cannot add")
	}
}
