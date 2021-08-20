/* Simple http server with SQLite db behind it that stores patients and their visits. */
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"kamilachyla.com/go/dent/data"
	"kamilachyla.com/go/dent/data/db"
)

func maintest() {
	// TODO remove Mainb
	db.MainDb()
}
func createFuncMap() template.FuncMap {
	return template.FuncMap{
		"equal": func(a int, b int) bool {
			return a == b
		},
	}
}

func main() {
	r := gin.Default()
	r.SetFuncMap(createFuncMap())
	r.Static("assets", "./assets")
	r.LoadHTMLGlob("templates/*")

	var patients = db.GetPatientsDao()

	// do testowania
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

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

	type FormData struct {
		Action string
		Method string
	}

	r.GET("/patient/:id", func(c *gin.Context) {
		v, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			v = -1
		}
		var all = patients.GetAll()
		c.HTML(http.StatusOK, "patients.tmpl", gin.H{
			"title":    "Pacjenci",
			"nav":      getNav(),
			"pat":      all,
			"curr":     patients.GetById(v),
			"formdata": FormData{Action: fmt.Sprintf("/patient/%d", v), Method: "POST"},
		})
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
		log.Printf("BindPatient: post form id is %s", formId)
		if id, err := strconv.Atoi(formId); err == nil {
			patWithNote.Id = id
		} else {
			return patWithNote, err
		}

		patWithNote.Surname = c.PostForm("surname")
		patWithNote.Birthdate, err = time.Parse("2006-01-02", c.DefaultPostForm("birthdate", "2021-10-12"))
		patWithNote.Note = c.PostForm("note")
		noteId := c.PostForm("note_id")
		log.Printf("BindPatient: post form note_id is %s", noteId)
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
		//log.Printf("Pat with note P.OST /newpatient: %v\n", patWithNote)
		if err == nil {
			id, err = patients.Add(patWithNote)
		}

		patWithNote.Id = id
		c.HTML(http.StatusOK, "patients.tmpl", gin.H{
			"title":    "Dane pacjenta",
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

	r.GET("/visits", func(c *gin.Context) {
		c.HTML(http.StatusOK, "visits.tmpl", gin.H{
			"title": "Wizyty",
			"nav":   getNav(),
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

type TemplMap struct {
	Path  string
	Templ string
}

func getNav() map[string]TemplMap {
	return map[string]TemplMap{
		"Pacjenci": {"/patients", "patients.tmpl"},
		"Wizyty":   {"/visits", "visits.tmpl"},
	}
}
