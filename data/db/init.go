package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

/* global, lazily initialized database handle*/
var db *sql.DB

/* Default database name in case it is not given in program arguments */
const DEFAULT_DATABASE = "./database.db"

/* Creates database in given location and runs db/create.sql*/
func createAndInitDb(fname string) (err error) {
	log.Printf("Creating empty database %s\n", fname)

	if db, err := sql.Open("sqlite3", fname); err == nil {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
				log.Printf("In defer, recover != nil, %v", err)
			}
			db.Close()
		}()

		log.Println("Reading file db/create.sql")
		c, ioErr := ioutil.ReadFile("db/create.sql")
		checkErr(ioErr)
		sql := string(c)
		_, err := db.Exec(sql)
		checkErr(err)
	}
	return err
}

/* Gets the name of database passed in commandline arg or
*  DEFAULT_DATABASE */
func getDbFileName() string {
	var fname = DEFAULT_DATABASE
	if len(os.Args) > 1 {
		fname = os.Args[1]
	}
	return fname
}

func createIfNotExist(fname string) (err error) {
	log.Printf("CreateIfNotExist: fname: %s\n", fname)
	if fi, err := os.Stat(fname); err != nil {

		log.Printf("CreateIfNotExist: stat: %v\n", fi)
		if os.IsNotExist(err) {
			return createAndInitDb(fname)
		}
		if os.IsPermission(err) {
			return errors.New(fmt.Sprintf("Nie masz uprawnień: %v\n", err))
		}
		return errors.New(fmt.Sprintf("Status pliku %s nieznany", fname))
	}

	return
}

func openOrCreate(dbname string) (err error) {
	log.Println("In openOrCreate")
	if err = createIfNotExist(dbname); err == nil {
		log.Println("Opening db", dbname)
		connstr := fmt.Sprintf("%s?_fk=true", dbname)
		db, err = sql.Open("sqlite3", connstr)
	}
	return
}

/* Returns package-local var db *sql.DB - initialized lazily */
func GetDatabase() (err error) {
	log.Println("In Get Database")
	if db == nil {
		if err = openOrCreate(getDbFileName()); err != nil {
			return err
		}

		if err = db.Ping(); err != nil {
			return err
		}

	}
	return nil
}
