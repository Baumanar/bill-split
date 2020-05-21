package main

import (
	"bill-splitting/data"
	"log"
	"net/http"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "password"
	DB_NAME     = "test"
)

func main() {
	var app App
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	app.Initialize(DB_USER, DB_PASSWORD, DB_NAME)
	//data.SetupDB()
	//PopulateDB()
	app.SetRoutes()
	app.Run()

}



func PopulateDB(){
	billSplit, err := data.CreateBillSplit("Colloc")
	billSplit.CreateParticipant("Robin")
	billSplit.CreateParticipant("Arnaud")
	billSplit.CreateParticipant("Mallaury")



	expense1, err := billSplit.CreateExpense("Biere", 30.0, "Robin")
	if err != nil {
		log.Fatal(err)
	}
	expense1.AddParticipant("Mallaury")
	expense1.AddParticipant("Robin")
	expense1.AddParticipant("Arnaud")

	expense1, err = billSplit.CreateExpense("Pizza", 33.5, "Arnaud")
	if err != nil {
		log.Fatal(err)
	}
	expense1.AddParticipant("Mallaury")
	expense1.AddParticipant("Arnaud")
	expense1.AddParticipant("Robin")

	billSplit, err = data.CreateBillSplit("Vacances")
	billSplit.CreateParticipant("Aline")
	billSplit.CreateParticipant("Arnaud")
	billSplit.CreateParticipant("Christine")
	billSplit.CreateParticipant("Bernard")
	billSplit.CreateParticipant("Patrick")
	billSplit.CreateParticipant("Elise")

	expense1, err = billSplit.CreateExpense("Courses", 30.65, "Arnaud")
	if err != nil {
		log.Fatal(err)
	}
	expense1.AddParticipant("Elise")
	expense1.AddParticipant("Arnaud")
	expense1.AddParticipant("Patrick")

	expense2, err := billSplit.CreateExpense("Voyage", 130.20, "Elise")
	if err != nil {
		log.Fatal(err)
	}
	expense2.AddParticipant("Elise")
	expense2.AddParticipant("Patrick")
	expense2.AddParticipant("Arnaud")
	expense2.AddParticipant("Christine")
	expense2.AddParticipant("Aline")


	expense3, err := billSplit.CreateExpense("Picnic", 80.77, "Christine")
	if err != nil {
		log.Fatal(err)
	}
	expense3.AddParticipant("Elise")
	expense3.AddParticipant("Arnaud")
	expense3.AddParticipant("Christine")
	expense3.AddParticipant("Aline")
}
