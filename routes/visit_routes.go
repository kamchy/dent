package routes

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

func DefineVisitRoutes(r *gin.Engine) {
	visits, err := db.GetVisitsDao()
	ExitIfErr(err)
	patients, err := db.GetPatientsDao()
	ExitIfErr(err)
	changes, err := db.GetChangeDao()
	ExitIfErr(err)
	/* Displays all visits */
	r.GET("/visits", func(c *gin.Context) {
		c.HTML(http.StatusOK, "visits.tmpl", gin.H{
			"title":  "Gabinet",
			"nav":    getNav(),
			"visits": visits.GetAll(),
		})
	})

	/* Displays single visit */
	r.GET("/visits/:id", func(c *gin.Context) {
		v := GetVId(c)
		vis := visits.GetById(v)
		c.HTML(http.StatusOK, "visits.tmpl", gin.H{
			"title":    "Gabinet",
			"nav":      getNav(),
			"visits":   visits.GetAll(),
			"currvis":  vis,
			"formdata": FormData{Action: vis.GetLink(), Method: "POST"},
		})
	})

	/* Displays list of visits for given patient */
	r.GET("patients/:id/visits", func(c *gin.Context) {
		p := GetId(c)
		c.HTML(http.StatusOK, "visits.tmpl", gin.H{
			"title":  "Gabinet",
			"nav":    getNav(),
			"visits": visits.GetByPatientId(p),
			"curr":   patients.GetById(p),
		})
	})

	r.GET("patients/:id/visits/new", func(c *gin.Context) {
		p := GetId(c)
		pat := patients.GetById(p)
		vis := data.NewForPatient(time.Now(), pat)

		c.HTML(http.StatusOK, "visits.tmpl", gin.H{
			"title":    "Nowa wizyta",
			"nav":      getNav(),
			"visits":   visits.GetByPatientId(p),
			"currvis":  vis,
			"curr":     pat,
			"formdata": FormData{Action: fmt.Sprintf("/patients/%d/visits/new", p), Method: "POST"},
		})
	})

	var BindVisit = func(c *gin.Context) (data.Visit, error) {
		var vis data.Visit
		if visId, err := strconv.Atoi(c.PostForm("vid")); err != nil {
			return vis, err
		} else {
			vis.Id = visId
		}

		if noteId, err := strconv.Atoi(c.PostForm("nid")); err != nil {
			return vis, err
		} else {
			vis.NoteId = noteId
		}
		vis.Note = c.PostForm("note")

		if patientId, err := strconv.Atoi(c.PostForm("pid")); err != nil {
			return vis, err
		} else {
			vis.PatientId = patientId
		}
		vis.Note = c.PostForm("note")

		var visitdateonly time.Time
		var visittimeonly time.Time

		if visitdateonly, err = time.ParseInLocation(data.VISIT_DATE_ONLY_LAYOUT, c.DefaultPostForm("visitdateonly", time.Now().Local().Format(data.VISIT_DATE_ONLY_LAYOUT)), time.Local); err != nil {
			return vis, err
		}
		if visittimeonly, err = time.ParseInLocation(data.VISIT_TIME_ONLY_LAYOUT, c.DefaultPostForm("visittimeonly", time.Now().Local().Format(data.VISIT_TIME_ONLY_LAYOUT)), time.Local); err != nil {
			return vis, err
		}

		vis.VisitDate = vis.From(visitdateonly, visittimeonly)

		return vis, err

	}
	r.POST("patients/:id/visits/new", func(c *gin.Context) {
		vis, err := BindVisit(c)
		var v int
		if err == nil {
			v, err = visits.Add(vis)
			vis.Id = v
		}
		p := GetId(c)

		c.HTML(http.StatusOK, "visits.tmpl", gin.H{
			"title":    "Nowa wizyta",
			"nav":      getNav(),
			"visits":   visits.GetByPatientId(p),
			"currvis":  vis,
			"curr":     patients.GetById(p),
			"formdata": FormData{Action: vis.GetLink(), Method: "POST"},
		})
	})

	r.POST("visits/:id/delete", func(c *gin.Context) {
		v := GetId(c)
		status := http.StatusFound
		log.Printf("POST visits/%d/delete", v)
		if err := visits.Delete(v); err != nil {
			status = http.StatusGone
		}
		c.Redirect(status, "/visits")
	})

	var respondWithVisitData = func(c *gin.Context, vis *data.Visit) {
		v := GetVId(c)
		p := GetId(c)

		var sts []data.State
		if err == nil {
			sts, err = changes.GetStates()
		}
		var chs []data.Change
		chs, err = changes.AllReversed(v)

		c.HTML(http.StatusOK, "visits.tmpl", gin.H{
			"title":    "Wizyta",
			"nav":      getNav(),
			"visits":   visits.GetByPatientId(p),
			"currvis":  visits.GetById(v),
			"curr":     patients.GetById(p),
			"formdata": FormData{Action: vis.GetLink(), Method: "POST"},
			"states":   sts,
			"changes":  chs,
			"err":      err,
			"pid":      p,
			"vid":      v,
		})
	}
	r.GET("patients/:id/visits/:vid", func(c *gin.Context) {
		v := GetVId(c)
		vis := visits.GetById(v)
		respondWithVisitData(c, vis)
	})

	r.POST("patients/:id/visits/:vid", func(c *gin.Context) {
		vis, err := BindVisit(c)
		if err == nil {
			err = visits.UpdateVisit(vis)
		}
		respondWithVisitData(c, &vis)
	})

}
