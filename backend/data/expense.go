package data

import (
	"log"
)

type Expense struct {
	Id        int
	Uuid      string
	Name      string
	Amount    float64
	BillSplitID   int
	PayerName  string
	CreatedAt JSONTime
}


// get posts to a thread
func (expense *Expense) ExpenseParticipants() (items []string, err error) {
	//defer db.Close()
	rows, err := Db.Query(    "SELECT p.name FROM participant_expense pe INNER JOIN participant p ON p.id = pe.participant_id WHERE pe.expense_id = $1", expense.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		var participant string
		if err = rows.Scan(&participant); err != nil {
			return
		}
		items = append(items, participant)
	}
	rows.Close()
	return
}

// get posts to a thread
func (expense *Expense) AddParticipant(name string) (err error) {
	//defer db.Close()
	participant, err := ParticipantByName(name)
	if err != nil{
		return
	}
	participantId := participant.Id
	statement := "insert into participant_expense(participant_id, expense_id) values ($1, $2) returning id, participant_id, expense_id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	_, err  = stmt.Exec(participantId, expense.Id)
	return
}


func (expense Expense) Balance() map[string] float64{
	participants, err := expense.ExpenseParticipants()
	if err != nil {
		log.Fatal(err)
	}
	payer, err := ParticipantByName(expense.PayerName)
	balance := make(map[string] float64)
	balance[payer.Name] = expense.Amount
	for _, participant := range participants {
		balance[participant] += -expense.Amount/float64(len(participants))
	}
	return balance
}