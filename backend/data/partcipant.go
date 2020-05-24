package data

type Participant struct {
	Id          int
	Uuid        string
	Name        string
	BillSplitID int
	CreatedAt   JSONTime
}
