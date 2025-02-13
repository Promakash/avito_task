package domain

type MerchID = int
type MerchName = string

type Merch struct {
	ID    MerchID
	Name  MerchName
	Price int
}
