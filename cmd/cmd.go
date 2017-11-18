package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/iancoleman/strcase"

	"gogit/database/migrations"
	"gogit/database/seeders"
)

var (
	create  *string
	command *string
)

// Basic flag declarations are available for string, integer, and boolean options.
func init() {
	command = flag.String("run", "", "The command")
	create = flag.String("create", "development", "a string")
}

func main() {

	flag.Parse()

	if strings.EqualFold(*command, "") {
		fmt.Println("No command specified")
		os.Exit(1)
	}
	fmt.Println("loading ", *command)

	switch *command {
	case "migrate":
		{
			migrations.Migrate()
			fmt.Println("Database migrated")
		}
	case "db:seed":
		{
			seeders.Seed()
			fmt.Println("Seed finished")
		}
	case "migrate:rollback":
		{
			// migrations.Rollback()
			// fmt.Println("Database Rolled back")
		}
	case "make:migration":
		{
			timestamp := time.Now().UnixNano()
			filename := fmt.Sprintf("../database/migrations/%v_%v.go", timestamp, strcase.ToSnake(*create))
			content := fmt.Sprintf(`package migrations

import "log"
import "gogit/database"

type %v struct{}

func (%v) Migrate() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Exec("")
	if e != nil {
		log.Println(e)
	}

}

func (%v) Rollback() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Exec("")
	if e != nil {
		log.Println(e)
	}

}`, strcase.ToCamel(*create), strcase.ToCamel(*create), strcase.ToCamel(*create))
			f, e := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0777)
			if e != nil {
				fmt.Println(e)
				os.Exit(1)
			}
			defer f.Close()
			f.WriteString(content)
		}
	}

}