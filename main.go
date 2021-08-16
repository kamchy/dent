package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"kamilachyla.com/go/dent/data/db"
	dat "kamilachyla.com/go/dent/data/db"
)

func maintest() {
	db.MainDb()
}
func main() {
	r := gin.Default()
	r.Static("assets", "./assets")
	r.LoadHTMLGlob("templates/*")

	var patients = dat.GetPatientsDao()

	// do testowania
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	var getIndex = func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Gabinet",
			"nav":   getNav(),
		})
	}

	r.GET("/index", getIndex)
	r.GET("/", getIndex)

	r.GET("/patients", func(c *gin.Context) {
		c.HTML(http.StatusOK, "patients.tmpl", gin.H{
			"title": "Pacjenci",
			"nav":   getNav(),
			"pat":   patients.GetAll(),
		})
	})

	r.GET("/patients/:id", func(c *gin.Context) {
		v, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			v = -1
		}
		c.HTML(http.StatusOK, "patients.tmpl", gin.H{
			"title":           "Pacjenci",
			"nav":             getNav(),
			"pat":             patients.GetAll(),
			"SelectedPatient": v,
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
		"Gabinet":  {"/index", "index.tmpl"},
		"Pacjenci": {"/patients", "patients.tmpl"},
		"Wizyty":   {"/visits", "visits.tmpl"},
	}
}
