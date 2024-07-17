package storage

import (
	"errors"
	"math"
	"nedas/shop/pkg/models"

	"github.com/gocql/gocql"
)

// thread safe
type Cassandra struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

// todo: env
func NewCassandra() (*Cassandra, error) {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "ma_shop"

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
		if ok := iter.Scan(&product.UserID, &product.ProductId, &product.Amount); !ok {
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

func (c *Cassandra) GetProductAmount(userId string, tid string, mid string) (uint8, error) {
	assert(userId != "", "user id is empty")
	assert(tid != "", "thread id is empty")
	assert(mid != "", "mid id is empty")

	var amount uint8

	productId := tid + ":" + mid
	query := c.session.Query(
		"SELECT amount FROM products WHERE user_id = ? AND product_id = ?",
		userId,
		productId,
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

func (c *Cassandra) AddProduct(userId string, tid string, mid string) error {
	assert(userId != "", "user id is empty")
	assert(tid != "", "thread id is empty")
	assert(mid != "", "mid is empty")

	productId := tid + ":" + mid
	query := c.session.Query(
		"INSERT INTO products (user_id, product_id, amount) VALUES (?, ?, 1) IF NOT EXISTS",
		userId,
		productId,
	)

	applied, err := query.ScanCAS(nil, nil, nil)
	if err != nil {
		return &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}

	if !applied {
		return &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: ErrAlreadySet}
	}

	return nil
}

func (c *Cassandra) IncreaseProduct(userId string, tid string, mid string) (uint8, error) {
	assert(userId != "", "user id is empty")
	assert(tid != "", "thread id is empty")
	assert(mid != "", "mid is empty")

	amount, err := c.GetProductAmount(userId, tid, mid)
	if err != nil {
		return 0, err
	}

	if amount == math.MaxUint8 {
		return math.MaxUint8, nil
	}

	productId := tid + ":" + mid
	query := c.session.Query(
		"UPDATE products SET amount = ? WHERE user_id = ? AND product_id = ?",
		amount+1,
		userId,
		productId,
	)

	if err := query.Exec(); err != nil {
		return 0, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}
	return amount + 1, nil
}

func (c *Cassandra) DecreaseProduct(userId string, tid string, mid string) (uint8, error) {
	assert(userId != "", "user id is empty")
	assert(tid != "", "thread id is empty")
	assert(mid != "", "mid is empty")

	amount, err := c.GetProductAmount(userId, tid, mid)
	if err != nil {
		return 0, err
	}

	var query *gocql.Query
	productId := tid + ":" + mid

	if amount == 1 {
		query = c.session.Query(
			"DELETE FROM products WHERE user_id = ? AND product_id = ?",
			userId,
			productId,
		)
	} else {
		query = c.session.Query(
			"UPDATE products SET amount = ? WHERE user_id = ? AND product_id = ?",
			amount-1,
			userId,
			productId,
		)
	}

	if err := query.Exec(); err != nil {
		return 0, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}
	return amount - 1, nil
}

func (c *Cassandra) DeleteProduct(userId string, tid string, mid string) error {
	assert(userId != "", "user id is empty")
	assert(tid != "", "thread id is empty")
	assert(mid != "", "mid is empty")

	productId := tid + ":" + mid
	query := c.session.Query(
		"DELETE FROM products WHERE user_id = ? AND product_id = ?",
		userId,
		productId,
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
