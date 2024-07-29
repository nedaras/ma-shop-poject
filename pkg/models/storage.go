package models

type StorageProduct struct {
	UserID    string
	ProductId string
	Size      string
	Amount    uint8
}

type StorageUser struct {
	UserID         string
	Email          string
	DefaultAddress string
}

type Address struct {
	AddressId   string
	Contact     string
	CountryCode string
	Phone       string
	Country     string
	Street      string
	Region      string
	City        string
	Zipcode     string
}
