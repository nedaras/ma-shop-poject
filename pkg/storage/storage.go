package storage

import (
	"errors"
	"nedas/shop/pkg/models"
)

var (
	ErrAlreadySet = errors.New("row already is set")
	ErrNotFound   = errors.New("row not found")
)

type Storage interface {

	// Any returned error should be of type [*StorageError].
	AddUser(user models.StorageUser) error

	// Any returned error should be of type [*StorageError].
	RemoveUser(userId string) error

	// Any returned error should be of type [*StorageError].
	GetUser(userId string) (models.StorageUser, error)

	// Any returned error should be of type [*StorageError].
	GetProducts(userId string) ([]models.StorageProduct, error)

	// Any returned error should be of type [*StorageError].
	GetProductAmount(userId string, tid string, mid string, size string) (uint8, error)

	// Any returned error should be of type [*StorageError].
	AddProduct(userId string, tid string, mid string, size string) (uint8, error)

	// Any returned error should be of type [*StorageError].
	IncreaseProduct(userId string, tid string, mid string, size string) (uint8, error)

	// Any returned error should be of type [*StorageError].
	DecreaseProduct(userId string, tid string, mid string, size string) (uint8, error)

	// Any returned error should be of type [*StorageError].
	DeleteProduct(userId string, tid string, mid string, size string) error

	// Any returned error should be of type [*StorageError].
	AddAddress(userId string, address models.Address) error

	// Any returned error should be of type [*StorageError].
	DeleteAddress(userId string, addressId string) error

	// Any returned error should be of type [*StorageError].
	GetAddresses(userId string) ([]models.Address, error)

	// Any returned error should be of type [*StorageError].
	GetAddress(userId string, addressId string) (models.Address, error)

	Close()
}

type StorageError struct {
	Provider  string
	Execution string
	Err       error
}

func (e *StorageError) Error() string {
	return e.Provider + " '" + e.Execution + "': " + e.Err.Error()
}

func (e *StorageError) Unwrap() error {
	return e.Err
}
