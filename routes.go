package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getParamOr(c *gin.Context, s string, def int) int {
	v, err := strconv.Atoi(c.Param(s))
	if err != nil {
		v = def
	}
	return v
}
func GetVId(c *gin.Context) int {
	return getParamOr(c, "vid", -1)
}

func GetId(c *gin.Context) int {
	return getParamOr(c, "id", -1)
}

func exitIfErr(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
