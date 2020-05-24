package data

import (
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

type BillSplit struct {
	Id        int
	Uuid      string
	Name      string
	CreatedAt JSONTime
}

//func (billSplit *BillSplit) CreatedAtDate() string {
//	return billSplit.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
//}

// get posts to a thread
func (billSplit *BillSplit) Participants() (items []Participant, err error) {
	//defer db.Close()
	rows, err := Db.Query("SELECT id, uuid, name, billsplit_id, created_at FROM participant where billsplit_id = $1 ORDER BY created_at DESC", billSplit.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Participant{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Name, &post.BillSplitID, &post.CreatedAt); err != nil {
			return
		}
		items = append(items, post)
	}
	rows.Close()
	return
}

// get posts to a thread
func (billSplit *BillSplit) Expenses() (items []Expense, err error) {
	//defer db.Close()
	rows, err := Db.Query("SELECT e.id, e.uuid, e.name, e.amount, e.billsplit_id, p.name, e.created_at FROM expense e INNER JOIN participant p ON e.participant_id = p.id where e.billSplit_id = $1 ORDER BY created_at DESC", billSplit.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Expense{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Name, &post.Amount, &post.BillSplitID, &post.PayerName, &post.CreatedAt); err != nil {
			return
		}
		items = append(items, post)
	}
	rows.Close()
	return
}

// get posts to a thread
func (billSplit *BillSplit) ExpenseByUuid(name string) (expense Expense, err error) {
	err = Db.QueryRow("SELECT e.id, e.uuid, e.name, e.amount, e.billsplit_id, p.name, e.created_at FROM expense e INNER JOIN participant p ON e.participant_id = p.id where e.uuid = $1 and e.billsplit_id = $2", name, billSplit.Id).
		Scan(&expense.Id, &expense.Uuid, &expense.Name, &expense.Amount, &expense.BillSplitID, &expense.PayerName, &expense.CreatedAt)
	return
}

// get posts to a thread
func (billSplit *BillSplit) ParticipantByName(name string) (participant Participant, err error) {
	err = Db.QueryRow("SELECT id, uuid, name, created_at FROM participant WHERE name = $1 and billsplit_id= $2", name, billSplit.Id).
		Scan(&participant.Id, &participant.Uuid, &participant.Name, &participant.CreatedAt)
	return
}

// get posts to a thread
func (billSplit *BillSplit) ParticipantsByName(names []string) (items []Participant, err error) {
	//defer db.Close()
	sqlStr := "SELECT id, uuid, name, billsplit_id, created_at FROM participant where billsplit_id = $1 and name in (?" + strings.Repeat(",?", len(names)-1) + ") ORDER BY created_at DESC"
	sqlStr = ReplaceSQL(sqlStr, "?", 2)

	args := make([]interface{}, len(names)+1)
	args[0] = billSplit.Id
	for i, id := range names {
		args[i+1] = id
	}
	rows, err := Db.Query(sqlStr, args...)
	// (?` + strings.Repeat(",?", len(args)-1) + `)`

	if err != nil {
		return
	}
	for rows.Next() {
		post := Participant{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Name, &post.BillSplitID, &post.CreatedAt); err != nil {
			return
		}
		items = append(items, post)
	}
	rows.Close()
	return
}

// Create a new item to a survey
func (billSplit *BillSplit) CreateParticipant(name string) (participant Participant, err error) {
	//defer db.Close()
	statement := "insert into participant (uuid, name, billsplit_id, created_at) values ($1, $2, $3, $4) returning id, uuid, name, billSplit_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(createUUID(), name, billSplit.Id, time.Now()).Scan(&participant.Id, &participant.Uuid, &participant.Name, &participant.BillSplitID, &participant.CreatedAt)
	if err != nil {
		return
	}
	return
}

// ReplaceSQL replaces the instance occurrence of any string pattern with an increasing $n based sequence
func ReplaceSQL(old, searchPattern string, startCount int) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := startCount; m <= (tmpCount + startCount - 1); m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

// Create a new item to a survey
func (billSplit *BillSplit) CreateParticipants(names []string) (err error) {

	sqlStr := "insert into participant(uuid, name, billsplit_id, created_at) VALUES "
	vals := []interface{}{}

	for _, row := range names {
		sqlStr += "(?, ?, ?, ?),"
		vals = append(vals, createUUID(), row, billSplit.Id, time.Now())
	}
	//trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	//Replacing ? with $n for postgres
	sqlStr = ReplaceSQL(sqlStr, "?", 1)

	//prepare the statement
	stmt, _ := Db.Prepare(sqlStr)

	//format all vals at once
	_, err = stmt.Exec(vals...)
	return
}

// Create a new item to a survey
func (billSplit *BillSplit) CreateExpense(name string, amount float64, participantName string) (expense Expense, err error) {
	//defer db.Close()
	participant, err := billSplit.ParticipantByName(participantName)
	if err != nil {
		return
	}
	_, err = Db.Exec("insert into expense (uuid, name, amount, billsplit_id, participant_id, created_at) values ($1, $2, $3, $4, $5, $6)", createUUID(), name, amount, billSplit.Id, participant.Id, time.Now())
	statement := "SELECT e.id, e.uuid, e.name, e.amount, e.billsplit_id, p.name, e.created_at FROM expense e INNER JOIN participant p ON e.participant_id = p.id where e.name = $1 and e.billsplit_id = $2 "
	if err != nil {
		return
	}
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = Db.QueryRow(statement, name, billSplit.Id).Scan(&expense.Id, &expense.Uuid, &expense.Name, &expense.Amount, &expense.BillSplitID, &expense.PayerName, &expense.CreatedAt)
	if err != nil {
		return
	}
	return
}

// Create a new item to a survey
func (billSplit *BillSplit) CreateExpenseParticipants(uuid string, participantNames []string) (err error) {
	expense, err := billSplit.ExpenseByUuid(uuid)
	for _, participant := range participantNames {
		err := expense.AddParticipant(participant)
		if err != nil {
			log.Fatal()
		}
	}
	if err != nil {
		return
	}
	return
}

// Create a new item to a survey
func (billSplit *BillSplit) GetFullBalance() (fullBalance map[string]float64, err error) {
	expenses, err := billSplit.Expenses()
	fullBalance = make(map[string]float64)
	participants, err := billSplit.Participants()
	for _, participant := range participants {
		fullBalance[participant.Name] = 0
	}
	for _, expense := range expenses {
		balanceExpense := expense.Balance()
		for k, v := range balanceExpense {
			fullBalance[k] += v
		}
	}
	if err != nil {
		return
	}
	return
}

type Debt struct {
	Debtor   string
	Creditor string
	Amount   float64
}

func (billSplit *BillSplit) GetDebts() (debts []Debt, err error) {
	debts = make([]Debt, 0)
	balance, err := billSplit.GetFullBalance()
	if err != nil {
		return
	}

	type kv struct {
		Key   string
		Value float64
	}
	var sortedBalance []kv
	for k, v := range balance {
		sortedBalance = append(sortedBalance, kv{k, v})
	}

	sort.Slice(sortedBalance, func(i, j int) bool {
		return sortedBalance[i].Value < sortedBalance[j].Value
	})

	i := 0
	j := len(sortedBalance) - 1
	var debt float64
	for i < j {
		debt = math.Min(-(sortedBalance[i].Value), math.Abs(sortedBalance[j].Value))

		sortedBalance[i].Value += debt
		sortedBalance[j].Value -= debt

		debts = append(debts, Debt{
			Debtor:   sortedBalance[i].Key,
			Creditor: sortedBalance[j].Key,
			Amount:   debt,
		})

		if sortedBalance[i].Value == 0 {
			i++
		}
		if sortedBalance[j].Value == 0 {
			j--
		}
	}
	return
}
