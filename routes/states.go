package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"kamilachyla.com/go/dent/data"
	"kamilachyla.com/go/dent/data/db"

	"net/http"
)

func DefineStateRoutes(r *gin.Engine) {
	changes, err := db.GetChangeDao()
	ExitIfErr(err)

	r.GET("/api/visits/:id/visitchanges", func(c *gin.Context) {
		id := GetId(c)
		if changes, err := changes.ForVisit(id); err == nil {
			c.JSON(http.StatusOK, changes)
		} else {
			c.JSON(http.StatusNotFound, fmt.Sprintf("There was error %s", err.Error()))
		}
	})

	r.GET("/api/visits/:id/renderchanges", func(c *gin.Context) {
		id := GetId(c)
		if changes, err := changes.AllReversed(id); err == nil {
			c.JSON(http.StatusOK, changes)
		} else {
			c.JSON(http.StatusNotFound, fmt.Sprintf("There was error %s", err.Error()))
		}
	})

	r.POST("/api/visits/:id/add", func(c *gin.Context) {
		ch := data.Change{}
		c.BindJSON(&ch)
		c.String(http.StatusOK, "Got change %v", ch)

	})
}
