package models

type StorageProduct struct {
	UserID    string
	ProductId string
	Size      string
	Amount    uint8
}

type StorageUser struct {
	UserID string
	Email  string
}
