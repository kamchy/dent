package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"kamilachyla.com/go/dent/data"
	"kamilachyla.com/go/dent/data/db"
)

type FormData struct {
	Action string
	Method string
}

func DefinePatientRoutes(r *gin.Engine) {
	patients, err := db.GetPatientsDao()
	exitIfErr(err)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Gabinet",
			"nav":   getNav(),
		})
	})

	r.GET("/patients", func(c *gin.Context) {
		c.HTML(http.StatusOK, "patients.tmpl", gin.H{
			"title": "Pacjenci",
			"nav":   getNav(),
			"pat":   patients.GetAll(),
		})
	})

	r.GET("/patient/:id", func(c *gin.Context) {
		v := GetId(c)
		var all = patients.GetAll()
		var data = gin.H{
			"title":    "Pacjenci",
			"nav":      getNav(),
			"pat":      all,
			"formdata": FormData{Action: fmt.Sprintf("/patient/%d", v), Method: "POST"},
		}
		if curr := patients.GetById(v); curr != nil {
			data["curr"] = curr
		}
		c.HTML(http.StatusOK, "patients.tmpl", data)
	})

	r.GET("/newpatient", func(c *gin.Context) {
		var patWithNote data.Patient
		c.HTML(http.StatusOK, "patients.tmpl", gin.H{
			"title":    "Nowy pacjent",
			"nav":      getNav(),
			"pat":      patients.GetAll(),
			"curr":     patWithNote,
			"formdata": FormData{Action: "/newpatient", Method: "POST"},
		})
	})

	var BindPatient = func(c *gin.Context) (data.Patient, error) {
		var patWithNote data.Patient

		var err error
		patWithNote.Name = c.PostForm("name")
		formId := c.PostForm("id")
		if id, err := strconv.Atoi(formId); err == nil {
			patWithNote.Id = id
		} else {
			return patWithNote, err
		}

		patWithNote.Surname = c.PostForm("surname")
		patWithNote.Birthdate, err = time.Parse("2006-01-02", c.DefaultPostForm("birthdate", "2021-10-12"))
		patWithNote.Note = c.PostForm("note")
		noteId := c.PostForm("note_id")
		if nid, err := strconv.Atoi(noteId); err == nil {
			patWithNote.NoteId = nid
		} else {
			patWithNote.NoteId = -1
		}
		return patWithNote, err
	}

	r.POST("/newpatient", func(c *gin.Context) {
		var patWithNote, err = BindPatient(c)
		var id int
		if err == nil {
			id, err = patients.Add(patWithNote)
		}

		patWithNote.Id = id
		c.HTML(http.StatusOK, "patients.tmpl", gin.H{
			"title":    "Nowy pacjent",
			"nav":      getNav(),
			"pat":      patients.GetAll(),
			"curr":     patients.GetById(id),
			"err":      err,
			"formdata": FormData{Action: fmt.Sprintf("/patient/%d", patWithNote.Id), Method: "POST"},
		})

	})

	r.POST("/patient/:id", func(c *gin.Context) {
		patWithNote, err := BindPatient(c)
		log.Printf("POSTED /patient/%d: %#v\n", patWithNote.Id, patWithNote)
		if err == nil {
			err = patients.UpdatePatient(&patWithNote)
		}

		pat := patients.GetById(patWithNote.Id)
		log.Printf("AFTER POSTED: %#v\n", pat)
		c.HTML(http.StatusOK, "patients.tmpl", gin.H{
			"title":    "Dane pacjenta",
			"nav":      getNav(),
			"pat":      patients.GetAll(),
			"curr":     patients.GetById(patWithNote.Id),
			"err":      err,
			"formdata": FormData{Action: fmt.Sprintf("/patient/%d", patWithNote.Id), Method: "POST"},
		})
	})

}
