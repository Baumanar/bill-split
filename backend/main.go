package main

import (
	"flag"
	"fmt"
	"github.com/Baumanar/bill-split/backend/data"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "password"
	DB_NAME     = "test"
)

func main() {
	var flagvar bool
	flag.BoolVar(&flagvar, "demo", false, "run the backend with demo data")
	flag.Parse()
	var app App
	app.Initialize(DB_USER, DB_PASSWORD, DB_NAME)
	if flagvar{
		// Reset Database
		fmt.Println("Running in demo mode")
		data.SetupDB()
		// Populate the databse
		PopulateDB()
	}

	app.SetRoutes()
	app.Run()

}


