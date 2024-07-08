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

	// Any returned error should be of type [*NikeAPIError].
	AddUser(user models.User) error

	// Any returned error should be of type [*NikeAPIError].
	RemoveUser(userId string) error

	// Any returned error should be of type [*NikeAPIError].
	GetUser(userId string) (models.User, error)

	// Any returned error should be of type [*NikeAPIError].
  GetProducts(userId string) ([]models.Product, error)

	// Any returned error should be of type [*NikeAPIError].
	IncreaseProduct(userId string, tid string, mid string) error

	// Any returned error should be of type [*NikeAPIError].
	DecreaseProduct(userId string, tid string, mid string) error

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
