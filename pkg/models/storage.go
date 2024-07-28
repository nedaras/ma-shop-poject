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
	Addresses      []Address
	DefaultAddress uint8
}

type Address struct {
	AddressId   uint8
	Contact     string
	CountryCode string
	Phone       string
	Country     string
	Street      string
	Region      string
	City        string
	Zipcode     string
}
