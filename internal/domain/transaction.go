package domain

type TransactionID = int

type Transaction struct {
	ID     TransactionID
	From   UserID
	To     UserID
	Amount int
}

type TransactionDirection int

const (
	Sent TransactionDirection = iota
	Received
)

type UserTransaction struct {
	OtherUser UserName
	Amount    int
	Direction TransactionDirection
}
