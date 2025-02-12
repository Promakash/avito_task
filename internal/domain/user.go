package domain

type UserID = int
type UserName = string
type UserHashPass = []byte

type User struct {
	ID             UserID
	Name           UserName
	HashedPassword UserHashPass
	Info           UserInfo
}

type UserInfo struct {
	Coins        int
	Transactions []UserTransaction
	Inventory    []Inventory
}
