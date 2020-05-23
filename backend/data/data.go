package data

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

var (
	DB_USER     = Getenv("DB_USER", "postgres")
	DB_PASSWORD     = Getenv("DB_PASSWORD", "password")
	DB_NAME     = Getenv("DB_NAME", "test_bill")
)


func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		fmt.Println(key, fallback)
		return fallback
	}
	fmt.Println(key, value)
	return value
}



type JSONTime time.Time

func (t JSONTime)MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("02 January 2006"))
	return []byte(stamp), nil
}


var Db *sql.DB

func InitDb() {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	var err error
	Db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	if err = Db.Ping(); err != nil {
		log.Fatal(err)
	}
}

// create a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// hash plaintext with SHA-1
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}


// Create a new survey
func CreateBillSplit(name string) (survey BillSplit, err error) {
	//defer db.Close()
	statement := "insert into billsplit (uuid, name, created_at) values ($1, $2, $3) returning id, uuid, name, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(createUUID(), name, time.Now()).Scan(&survey.Id, &survey.Uuid, &survey.Name, &survey.CreatedAt)
	if err != nil {
		return
	}
	return
}

// Get all threads in the database and returns it
func BillSplits() (billSplits []BillSplit, err error) {
	//defer db.Close()
	rows, err := Db.Query("SELECT id, uuid, name, created_at FROM billsplit ORDER BY created_at DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		conv := BillSplit{}
		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Name, &conv.CreatedAt); err != nil {
			return
		}
		billSplits = append(billSplits, conv)
	}
	rows.Close()
	return
}

// Get a thread by the UUID
func BillSplitByUUID(uuid string) (billSplit BillSplit, err error) {
	err = Db.QueryRow("SELECT id, uuid, name, created_at FROM billsplit WHERE uuid = $1", uuid).
		Scan(&billSplit.Id, &billSplit.Uuid, &billSplit.Name, &billSplit.CreatedAt)
	return
}

// Get a thread by the UUID
func BillSplitByID(id int) (billSplit BillSplit, err error) {
	err = Db.QueryRow("SELECT id, uuid, name, created_at FROM billsplit WHERE id = $1", id).
		Scan(&billSplit.Id, &billSplit.Uuid, &billSplit.Name, &billSplit.CreatedAt)
	return
}


// Get a thread by the UUID
func BillSplitByName(name string) (billSplit BillSplit, err error) {
	err = Db.QueryRow("SELECT id, uuid, name, created_at FROM billsplit WHERE name = $1", name).
		Scan(&billSplit.Id, &billSplit.Uuid, &billSplit.Name, &billSplit.CreatedAt)
	return
}


// get posts to a thread
func ExpenseByUuid(name string) (expense Expense, err error) {
	err = Db.QueryRow("SELECT e.id, e.uuid, e.name, e.amount, e.billsplit_id, p.name, e.created_at FROM expense e INNER JOIN participant p ON e.participant_id = p.id where e.uuid = $1", name).
		Scan(&expense.Id, &expense.Uuid, &expense.Name,  &expense.Amount, &expense.BillSplitID, &expense.PayerName,&expense.CreatedAt)
	return
}

// Get a thread by the UUID
func ParticipantByUUID(uuid string) (participant Participant, err error) {
	err = Db.QueryRow("SELECT id, uuid, name, billsplit_id, created_at FROM participant WHERE uuid = $1", uuid).
		Scan(&participant.Id, &participant.Uuid, &participant.Name, &participant.BillSplitID, &participant.CreatedAt)
	return
}

// Get a thread by the UUID
func ParticipantByName(uuid string) (participant Participant, err error) {
	err = Db.QueryRow("SELECT id, uuid, name, billsplit_id, created_at FROM participant WHERE name = $1", uuid).
		Scan(&participant.Id, &participant.Uuid, &participant.Name, &participant.BillSplitID, &participant.CreatedAt)
	return
}

// Get a thread by the UUID
func ParticipantByID(id int) (participant Participant, err error) {
	err = Db.QueryRow("SELECT id, uuid, name, billsplit_id, created_at FROM participant WHERE id = $1", id).
		Scan(&participant.Id, &participant.Uuid, &participant.Name, &participant.BillSplitID, &participant.CreatedAt)
	return
}

// Delete all Participant from database
func ParticipantDeleteAll() (err error) {
	//defer db.Close()
	statement := "delete  from participant"
	_, err = Db.Exec(statement)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// Delete all surveys from database
func ExpenseDeleteAll() (err error) {
	statement := "delete from expense"
	_, err = Db.Exec(statement)
	if err != nil {
		log.Fatal(err)
	}
	return
}


// Delete all users from database
func BillSplitDeleteAll() (err error) {
	statement := "delete from billsplit"
	_, err = Db.Exec(statement)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// Delete all users from database
func ParticipantExpenseDeleteAll() (err error) {
	statement := "delete from participant_expense"
	_, err = Db.Exec(statement)
	if err != nil {
		log.Fatal(err)
	}
	return
}


func SetupDB() {
	ParticipantExpenseDeleteAll()
	ExpenseDeleteAll()
	ParticipantDeleteAll()
	BillSplitDeleteAll()
}

