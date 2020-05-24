package data

import (
	"fmt"
	"log"
	"testing"
)

//func TestExpense_ParticipantExpense(t *testing.T) {
//	InitDb()
//	SetupDB()
//	tests := []struct {
//		name       string
//		wantSurvey []Participant
//		wantErr    bool
//	}{
//		{
//			"testSurvey",
//			Survey{
//				Id:        0,
//				Uuid:      "",
//				Name:      "testSurvey",
//				CreatedAt: time.Now()},
//			false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotSurvey, err := CreateSurvey(tt.name)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("CreateSurvey() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if gotSurvey.Name != tt.wantSurvey.Name {
//				t.Errorf("CreateSurvey() gotSurvey = %v, want %v", gotSurvey, tt.wantSurvey)
//			}
//		})
//	}
//	Db.Close()
//}

func TestExpense_Balance(t *testing.T) {
	InitDb()
	SetupDB()
	tests := []struct {
		name        string
		wantBalance map[string]float64
		wantErr     bool
	}{
		{
			"testSurvey",
			map[string]float64{
				"A": 30.0,
				"B": -10.0,
				"C": -10.0,
				"D": -10.0,
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
			expense1, err := billSplit.CreateExpense("expense1", 40.0, "A")
			if err != nil {
				fmt.Println("lol")
				log.Fatal(err)
			}
			expense1.AddParticipant("A")
			expense1.AddParticipant("B")
			expense1.AddParticipant("C")
			expense1.AddParticipant("D")

			balance := expense1.Balance()
			for k, _ := range balance {
				if balance[k] != tt.wantBalance[k] {
					t.Errorf("want %f, got %f", balance[k], tt.wantBalance[k])
				}
			}
			fmt.Println(balance)

		})
	}
	Db.Close()

}

func TestExpense_ExpenseParticipants(t *testing.T) {
	InitDb()
	//SetupDB()
	tests := []struct {
		name               string
		nameBill           string
		nameExpense        string
		participants       []string
		expenseParticipant []string
		wantErr            bool
	}{
		{
			"test0",
			"bill0",
			"expense0",
			[]string{"D", "C", "B", "A"},
			[]string{"A", "B", "C", "D"},
			false,
		},
		{
			"test1",
			"bill1",
			"expense1",
			[]string{"D", "C", "B", "A"},
			[]string{"A", "B", "C", "D", "E"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			billSplit, err := CreateBillSplit(tt.nameBill)
			if err != nil {
				log.Fatal(err)
			}
			for _, part := range tt.participants {
				billSplit.CreateParticipant(part)

			}
			expense1, err := billSplit.CreateExpense(tt.nameExpense, 10.0, tt.participants[0])
			if err != nil {
				log.Fatal(err)
			}
			for _, part := range tt.expenseParticipant {
				err := expense1.AddParticipant(part)
				if err != nil {
					if !tt.wantErr {
						t.Errorf("Got unwanted error")
					}
				}
			}
			gotParticipants, err := expense1.ExpenseParticipants()
			if err != nil {
				log.Fatal(err)
			}
			for idx, got := range gotParticipants {
				if got != tt.participants[idx] {
					t.Errorf("CreateSurvey() gotSurvey = %v, want %v", got, tt.participants[idx+1])
				}
			}

		})
	}
	Db.Close()
}

func TestExpense_AddParticipants(t *testing.T) {
	InitDb()
	SetupDB()

	t.Run("TestBillSplit_ExpenseByUuid", func(t *testing.T) {
		billSplit, err := CreateBillSplit("test0")
		names := []string{"A", "B", "C", "D"}
		wantNames := []string{"C", "B", "A"}
		billSplit.CreateParticipants(names)
		if err != nil {
			log.Fatal(err)
		}
		expense, err := billSplit.CreateExpense("testExpense", 100, "A")
		expense.AddParticipants([]string{"A", "B", "C"})
		parts, err := expense.ExpenseParticipants()
		if err != nil {
			log.Fatal(err)
		}
		for idx, name := range parts {
			if wantNames[idx] != name {
				t.Errorf("gotExpense = %v, want %v", name, wantNames[idx])
			}
		}

	})
	Db.Close()
}
