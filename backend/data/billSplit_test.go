package data

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestBillSplit_GetFullBalance(t *testing.T) {
	InitDb()
	SetupDB()
	tests := []struct {
		name       string
		wantBalance map[string]float64
		wantErr    bool
	}{
		{
			"testSurvey",
			map[string]float64{
				"A": -2.5,
				"B": 17.5,
				"C": -12.5,
				"D": -2.5,
			},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			billSplit, err := CreateBillSplit("test0")
			if err != nil {
				log.Fatal(err)
			}
			billSplit.CreateParticipant("A")
			billSplit.CreateParticipant("B")
			billSplit.CreateParticipant("C")
			billSplit.CreateParticipant("D")
			expense1, err := billSplit.CreateExpense("expense1", 10.0, "A")
			if err != nil {
				log.Fatal(err)
			}
			expense1.AddParticipant("A")
			expense1.AddParticipant("B")
			expense1.AddParticipant("C")
			expense1.AddParticipant("D")

			expense2, err := billSplit.CreateExpense("expense2", 30.0, "B")
			if err != nil {
				log.Fatal(err)
			}
			expense2.AddParticipant("A")
			expense2.AddParticipant("B")
			expense2.AddParticipant("C")

			gotBalance, err := billSplit.GetFullBalance()
			if err != nil {
				log.Fatal(err)
			}

			for k,_ := range gotBalance{
				if tt.wantBalance[k] != gotBalance[k]{
					t.Errorf("Balance() gotSurvey = %v, want %v", gotBalance[k], tt.wantBalance[k])
				}
			}
		})
	}
	Db.Close()

}

func TestBillSplit_CreateParticipant(t *testing.T) {
	InitDb()
	SetupDB()
	tests := []struct {
		name     string
		wantItem []string
		wantErr  bool
	}{
		{"test0", []string{"B", "A"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			billSplit, err := CreateBillSplit("test0")
			if err != nil {
				log.Fatal(err)
			}
			billSplit.CreateParticipant("A")
			billSplit.CreateParticipant("B")

			participants, err := billSplit.Participants()
			if err != nil {
				log.Fatal(err)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateParticipant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for idx, participant := range participants{
				if participant.Name != tt.wantItem[idx]{
					t.Errorf("CreateParticipant() gotItem = %v, want %v", participant.Name , tt.wantItem[idx])

				}
			}
		})
	}
	Db.Close()
}


func TestBillSplit_GetDebts(t *testing.T) {
	InitDb()
	SetupDB()

	tests := []struct {
		name      string
		wantDebts []Debt
		wantErr   bool
	}{
		{"test0", []Debt{
			{"C", "B", 12.5},
			{"A", "B", 2.5},
			{"D", "B", 2.5},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			billSplit, err := CreateBillSplit("test0")
			billSplit.CreateParticipant("A")
			billSplit.CreateParticipant("B")
			billSplit.CreateParticipant("C")
			billSplit.CreateParticipant("D")
			expense1, err := billSplit.CreateExpense("expense1", 10.0, "A")
			if err != nil {
				log.Fatal(err)
			}
			expense1.AddParticipant("A")
			expense1.AddParticipant("B")
			expense1.AddParticipant("C")
			expense1.AddParticipant("D")

			expense2, err := billSplit.CreateExpense("expense2", 30.0, "B")
			if err != nil {
				log.Fatal(err)
			}
			expense2.AddParticipant("A")
			expense2.AddParticipant("B")
			expense2.AddParticipant("C")

			gotDebts, err := billSplit.GetDebts()
			fmt.Println(gotDebts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDebts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDebts, tt.wantDebts) {
				t.Errorf("GetDebts() gotDebts = %v, want %v", gotDebts, tt.wantDebts)
			}
		})
	}
	Db.Close()
}

func TestBillSplit_ExpenseByUuid(t *testing.T) {
	InitDb()
	SetupDB()

	t.Run("TestBillSplit_ExpenseByUuid", func(t *testing.T) {
		billSplit, err := CreateBillSplit("test0")
		billSplit.CreateParticipant("A")

		expense, err := billSplit.CreateExpense("expense1", 10.0, "A")
		if err != nil {
			log.Fatal(err)
		}
		uuid := expense.Uuid
		gotExpense, err := billSplit.ExpenseByUuid(uuid)
		if !reflect.DeepEqual(gotExpense, expense){
			t.Errorf("gotExpense = %v, want %v", gotExpense, expense)
		}
	})
	Db.Close()
}

func TestBillSplit_CreateParticipants(t *testing.T) {
	InitDb()
	SetupDB()

	t.Run("TestBillSplit_ExpenseByUuid", func(t *testing.T) {
		billSplit, err := CreateBillSplit("test0")
		names := []string{"A", "B", "C", "D"}
		wantNames := []string{"D","C", "B", "A"}
		billSplit.CreateParticipants(names)
		if err != nil {
			log.Fatal(err)
		}
		gotParticipants, err := billSplit.Participants()
		if err != nil {
			log.Fatal(err)
		}
		gotNames := make([]string, 0)
		for _, participant := range gotParticipants{
			gotNames = append(gotNames, participant.Name)
		}
		for idx, name := range gotNames {
			if wantNames[idx] != name{
				t.Errorf("gotExpense = %v, want %v", name, wantNames[idx])
			}
		}

	})
	Db.Close()
}

func TestBillSplit_ParticipantsByName(t *testing.T) {
	InitDb()
	SetupDB()

	t.Run("TestBillSplit_ExpenseByUuid", func(t *testing.T) {
		billSplit, err := CreateBillSplit("test0")
		names := []string{"A", "B", "C", "D"}
		wantNames := []string{"C", "B", "A"}
		billSplit.CreateParticipants(names)
		p, err := billSplit.Participants()
		fmt.Println("got", p)

		if err != nil {
			log.Fatal(err)
		}
		gotParticipants, err := billSplit.ParticipantsByName([]string{"A", "B", "C"})
		if err != nil {
			log.Fatal(err)
		}
		gotNames := make([]string, 0)
		for _, participant := range gotParticipants{
			gotNames = append(gotNames, participant.Name)
		}
		for idx, name := range wantNames {
			if gotNames[idx] != name{
				t.Errorf("gotExpense = %v, want %v", name, wantNames[idx])
			}
		}

	})
	Db.Close()
}