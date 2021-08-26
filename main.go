/* Simple http application with SQLite db behind it that stores patients and their visits. */
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
)

func main() {
	r := gin.Default()
	r.Static("assets", "./assets")
	r.Use(favicon.New("assets/favicon.ico"))
	r.LoadHTMLGlob("templates/*")
	DefinePatientRoutes(r)
	DefineVisitRoutes(r)

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
