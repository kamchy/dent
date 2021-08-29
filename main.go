/* Simple http application with SQLite db behind it that stores patients and their visits. */
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
	"kamilachyla.com/go/dent/routes"
)

func main() {
	r := gin.Default()
	r.Static("assets", "./assets")
	r.Use(favicon.New("assets/favicon.ico"))
	r.LoadHTMLGlob("templates/*")
	routes.DefinePatientRoutes(r)
	routes.DefineVisitRoutes(r)
	routes.DefineStateRoutes(r)

	r.Run() // listen and serve on 0.0.0.0:8080
}
