package models

type Product struct {
	Title    string
	Price    float64
	Image    string
	PathName string
	Mid      string
	ThreadId string
	Slug     string
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
