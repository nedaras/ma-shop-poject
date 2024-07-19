package storage

import (
	"errors"
	"log"
	"math"
	"nedas/shop/pkg/models"
	"os"

	"github.com/gocql/gocql"
)

type Cassandra struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

func NewCassandra() (*Cassandra, error) {
	address := os.Getenv("CASSANDRA_ADDRESS")
	if address == "" {
		log.Fatal("CASSANDRA_ADDRESS not found in .env")
	}

	keyspace := os.Getenv("CASSANDRA_KEYSPACE")
	if keyspace == "" {
		log.Fatal("CASSANDRA_KEYSPACE not found in .env")
	}

	cluster := gocql.NewCluster(address)
	cluster.Keyspace = keyspace

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &Cassandra{
		cluster: cluster,
		session: session,
	}, nil
}

func (c *Cassandra) AddUser(user models.User) error {
	assert(user.UserID != "", "user id is empty")
	assert(user.Email != "", "user email is empty")

	query := c.session.Query(
		"INSERT INTO users(user_id, email) VALUES (?, ?)",
		user.UserID,
		user.Email,
	)

	if err := query.Exec(); err != nil {
		return &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}

	return nil
}

func (c *Cassandra) RemoveUser(userId string) error {
	assert(userId != "", "user id is empty")

	query := c.session.Query(
		"DELETE FROM users WHERE user_id = ?",
		userId,
	)

	return query.Exec()
}

func (c *Cassandra) GetUser(userId string) (models.User, error) {
	assert(userId != "", "user id is empty")

	query := c.session.Query(
		"SELECT * FROM users WHERE user_id = ?",
		userId,
	)

	user := models.User{}
	iter := query.Iter()

	if iter.NumRows() == 0 {
		err := iter.Close()
		if err != nil {
			if errors.Is(err, gocql.ErrNotFound) {
				err = ErrNotFound
			}
			return models.User{}, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
		}
		return models.User{}, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: ErrNotFound}
	}

	ok := iter.Scan(&user.UserID, &user.Email)
	if err := iter.Close(); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			err = ErrNotFound
		}
		return models.User{}, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}

	if !ok {
		return models.User{}, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: ErrNotFound}
	}

	return user, nil
}

func (c *Cassandra) GetProducts(userId string) ([]models.Product, error) {
	assert(userId != "", "user id is empty")

	query := c.session.Query(
		"SELECT * FROM products WHERE user_id = ?",
		userId,
	)

	iter := query.Iter()

	if iter.NumRows() == 0 {
		if err := iter.Close(); err != nil {
			return []models.Product{}, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
		}
		return []models.Product{}, nil
	}

	products := make([]models.Product, iter.NumRows())
	for i := range iter.NumRows() {
		product := &products[i]
		if ok := iter.Scan(&product.UserID, &product.ProductId, &product.Size, &product.Amount); !ok {
			err := iter.Close()
			if err != nil {
				panic("no err and not ok!!!!")
			}
			return []models.Product{}, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
		}
	}

	if err := iter.Close(); err != nil {
		return []models.Product{}, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}

	return products, nil
}

func (c *Cassandra) GetProductAmount(userId string, tid string, mid string, size string) (uint8, error) {
	assert(userId != "", "user id is empty")
	assert(tid != "", "thread id is empty")
	assert(mid != "", "mid is empty")
	assert(size != "", "size is empty")

	var amount uint8

	productId := tid + ":" + mid
	query := c.session.Query(
		"SELECT amount FROM products WHERE user_id = ? AND product_id = ? AND size = ?",
		userId,
		productId,
		size,
	)

	iter := query.Iter()
	if iter.NumRows() == 0 {
		if err := iter.Close(); err != nil {
			if errors.Is(err, gocql.ErrNotFound) {
				err = ErrNotFound
			}
			return 0, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
		}
		return 0, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: ErrNotFound}
	}

	if ok := iter.Scan(&amount); !ok {
		err := iter.Close()
		if err != nil {
			panic("no err and not ok!!!!")
		}
		return 0, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}

	if err := iter.Close(); err != nil {
		return 0, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}

	return amount, nil
}

func (c *Cassandra) AddProduct(userId string, tid string, mid string, size string) error {
	assert(userId != "", "user id is empty")
	assert(tid != "", "thread id is empty")
	assert(mid != "", "mid is empty")
	assert(size != "", "size is empty")

	// todo: it should be like increment but u know would not err if not item
	productId := tid + ":" + mid
	query := c.session.Query(
		"INSERT INTO products (user_id, product_id, size, amount) VALUES (?, ?, ?, 1) IF NOT EXISTS",
		userId,
		productId,
		size,
	)

	applied, err := query.ScanCAS(nil, nil, nil, nil)
	if err != nil {
		return &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}

	if !applied {
		return &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: ErrAlreadySet}
	}

	return nil
}

func (c *Cassandra) IncreaseProduct(userId string, tid string, mid string, size string) (uint8, error) {
	assert(userId != "", "user id is empty")
	assert(tid != "", "thread id is empty")
	assert(mid != "", "mid is empty")
	assert(size != "", "size is empty")

	amount, err := c.GetProductAmount(userId, tid, mid, size)
	if err != nil {
		return 0, err
	}

	if amount == math.MaxUint8 {
		return math.MaxUint8, nil
	}

	productId := tid + ":" + mid
	query := c.session.Query(
		"UPDATE products SET amount = ? WHERE user_id = ? AND product_id = ? AND size = ?",
		amount+1,
		userId,
		productId,
		size,
	)

	if err := query.Exec(); err != nil {
		return 0, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}
	return amount + 1, nil
}

func (c *Cassandra) DecreaseProduct(userId string, tid string, mid string, size string) (uint8, error) {
	assert(userId != "", "user id is empty")
	assert(tid != "", "thread id is empty")
	assert(mid != "", "mid is empty")
	assert(size != "", "size is empty")

	amount, err := c.GetProductAmount(userId, tid, mid, size)
	if err != nil {
		return 0, err
	}

	var query *gocql.Query
	productId := tid + ":" + mid

	if amount == 1 {
		query = c.session.Query(
			"DELETE FROM products WHERE user_id = ? AND product_id = ? AND size = ?",
			userId,
			productId,
			size,
		)
	} else {
		query = c.session.Query(
			"UPDATE products SET amount = ? WHERE user_id = ? AND product_id = ? AND size = ?",
			amount-1,
			userId,
			productId,
			size,
		)
	}

	if err := query.Exec(); err != nil {
		return 0, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}
	return amount - 1, nil
}

func (c *Cassandra) DeleteProduct(userId string, tid string, mid string, size string) error {
	assert(userId != "", "user id is empty")
	assert(tid != "", "thread id is empty")
	assert(mid != "", "mid is empty")
	assert(size != "", "size is empty")

	productId := tid + ":" + mid
	query := c.session.Query(
		"DELETE FROM products WHERE user_id = ? AND product_id = ? AND size = ?",
		userId,
		productId,
		size,
	)

	if err := query.Exec(); err != nil {
		return &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}

	return nil
}

func (c *Cassandra) Close() {
	c.session.Close()
}

func assert(ok bool, msg string) {
	if !ok {
		panic(msg)
	}
}
