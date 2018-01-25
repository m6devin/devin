package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/iancoleman/strcase"

	"gogit/database"
	"gogit/database/migrations"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type migration struct {
	ID    uint
	Name  string
	Batch uint
}

// MakeMigration create new migration file
func MakeMigration(create *string) {
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("./database/migrations/%v_%v.go", timestamp, strcase.ToSnake(*create))
	content := fmt.Sprintf(`package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) Migrate%v() (e error) {
    db := database.NewPGInstance()
    defer db.Close()
    _, e = db.Exec("")

    return
}

// Rollback the database to previous version
func (Migration) Rollback%v() (e error) {
    db := database.NewPGInstance()
    defer db.Close()
    _, e = db.Exec("")

    return
}`, strcase.ToCamel(*create), strcase.ToCamel(*create))
	f, e := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0777)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	defer f.Close()
	f.WriteString(content)
}

// checkMigrationTable will check existance of migrations table and create it if it dosen't exist.
func checkMigrationTable() {
	db := database.NewPGInstance()
	defer db.Close()

	db.Exec(`CREATE SCHEMA IF NOT EXISTS public;`)
	db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
    id serial NOT NULL,
    name varchar(255) NOT NULL,
    batch integer,
    CONSTRAINT migrations_pkey PRIMARY KEY (id)
    )`)
}

func Migrate() {
	checkMigrationTable()

	db := database.NewPGInstance()
	defer db.Close()

	files, e := ioutil.ReadDir("./database/migrations/")
	if e != nil {
		log.Fatalln("Error on loading migrations directory")
	}

	for _, v := range files {
		var mg migration
		filename := strings.TrimSuffix(v.Name(), ".go")
		db.Model(&mg).Where("name LIKE ?", filename).First()

		if mg.ID != 0 {
			//This file already migrated
			continue
		}

		name := strcase.ToCamel(filename)
		m := migrations.Migration{}
		val := reflect.ValueOf(m)
		f := val.MethodByName("Migrate" + name)
		if !f.IsValid() {
			continue
		}

		rets := f.Call(nil)
		if len(rets) == 0 {
			continue
		}

		if rets[0].Interface() != nil {
			continue
		}

		// save migrated file to DB
		mg.Name = filename
		db.Insert(&mg)
	}

	db.Exec("update migrations set batch=coalesce((select max(batch) from migrations) , 0)+1 where batch is null;")

}
