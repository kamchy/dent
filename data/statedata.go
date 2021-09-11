package data

import "fmt"

type State struct {
	Id      int
	Name    string
	Whole   bool
	Color   string
	Svgpath string
}

type Change struct {
	Id        int
	VisitId   int
	PatientId int
	StateId   int
	ToothNum  int
	ToothSide int
}

func NewChange(id int, visitid int, patientid int, stateid int, tn int, ts int) *Change {
	return &Change{Id: id, VisitId: visitid, PatientId: patientid, StateId: stateid, ToothNum: tn, ToothSide: ts}
}

func (c Change) String() string {
	return fmt.Sprintf("Change %d: vis[%d], pat[%d], state[%d], t[%d], side[%d]", c.Id, c.VisitId, c.PatientId, c.StateId, c.ToothNum, c.ToothSide)
}

type ChangeDao interface {
	ForVisit(visitId int) (changes []Change, err error)
	AllReversed(visitId int) (changes []Change, err error)
	InsertChange(data Change) (change *Change, err error)
	GetStates() (states []State, err error)
}
