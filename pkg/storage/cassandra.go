package storage

import (
	"errors"
	"math"
	"nedas/shop/pkg/models"
	"nedas/shop/pkg/utils"

	"github.com/gocql/gocql"
)

type Cassandra struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

func NewCassandra() (*Cassandra, error) {
	address := utils.Getenv("CASSANDRA_ADDRESS")
	keyspace := utils.Getenv("CASSANDRA_KEYSPACE")

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

func (c *Cassandra) AddUser(user models.StorageUser) error {
	utils.Assert(user.UserID != "", "user id is empty")
	utils.Assert(user.Email != "", "user email is empty")

	query := c.session.Query(
		"INSERT INTO users(user_id, email, addresses, default_address) VALUES (?, ?, {}, 0)",
		user.UserID,
		user.Email,
	)

	if err := query.Exec(); err != nil {
		return &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}

	return nil
}

func (c *Cassandra) RemoveUser(userId string) error {
	// todo: implement
	panic("not implemented")
}

func (c *Cassandra) GetUser(userId string) (models.StorageUser, error) {
	utils.Assert(userId != "", "user id is empty")

	query := c.session.Query(
		"SELECT user_id, email, addresses, default_address FROM users WHERE user_id = ?",
		userId,
	)

	user := models.StorageUser{}
	iter := query.Iter()

	if iter.NumRows() == 0 {
		err := iter.Close()
		if err != nil {
			if errors.Is(err, gocql.ErrNotFound) {
				err = ErrNotFound
			}
			return models.StorageUser{}, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
		}
		return models.StorageUser{}, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: ErrNotFound}
	}

	var addresses []map[string]any
	ok := iter.Scan(&user.UserID, &user.Email, &addresses, &user.DefaultAddress)
	if err := iter.Close(); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			err = ErrNotFound
		}
		return models.StorageUser{}, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}

	if !ok {
		return models.StorageUser{}, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: ErrNotFound}
	}

	if len(addresses) > 0 {
		user.Addresses = make([]models.Address, len(addresses))
		for i, address := range addresses {
			user.Addresses[i] = models.Address{
				AddressId:   uint8(address["address_id"].(int8)),
				Contact:     address["contact"].(string),
				CountryCode: address["country_code"].(string),
				Phone:       address["phone"].(string),
				Country:     address["country"].(string),
				Street:      address["street"].(string),
				Region:      address["region"].(string),
				City:        address["city"].(string),
				Zipcode:     address["zipcode"].(string),
			}
		}
	}

	return user, nil
}

func (c *Cassandra) GetProducts(userId string) ([]models.StorageProduct, error) {
	utils.Assert(userId != "", "user id is empty")

	query := c.session.Query(
		"SELECT * FROM products WHERE user_id = ?",
		userId,
	)

	iter := query.Iter()

	if iter.NumRows() == 0 {
		if err := iter.Close(); err != nil {
			return []models.StorageProduct{}, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
		}
		return []models.StorageProduct{}, nil
	}

	products := make([]models.StorageProduct, iter.NumRows())
	for i := range iter.NumRows() {
		product := &products[i]
		if ok := iter.Scan(&product.UserID, &product.ProductId, &product.Size, &product.Amount); !ok {
			err := iter.Close()
			if err != nil {
				panic("no err and not ok!!!!")
			}
			return []models.StorageProduct{}, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
		}
	}

	if err := iter.Close(); err != nil {
		return []models.StorageProduct{}, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}

	return products, nil
}

func (c *Cassandra) GetProductAmount(userId string, tid string, mid string, size string) (uint8, error) {
	utils.Assert(userId != "", "user id is empty")
	utils.Assert(tid != "", "thread id is empty")
	utils.Assert(mid != "", "mid is empty")
	utils.Assert(size != "", "size is empty")

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

func (c *Cassandra) AddProduct(userId string, tid string, mid string, size string) (uint8, error) {
	utils.Assert(userId != "", "user id is empty")
	utils.Assert(tid != "", "thread id is empty")
	utils.Assert(mid != "", "mid is empty")
	utils.Assert(size != "", "size is empty")

	productId := tid + ":" + mid
	var query *gocql.Query

	amount, err := c.GetProductAmount(userId, tid, mid, size)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			query = c.session.Query(
				"INSERT INTO products (user_id, product_id, size, amount) VALUES (?, ?, ?, 1)",
				userId,
				productId,
				size,
			)
		} else {
			return 0, err
		}
	} else {
		query = c.session.Query(
			"UPDATE products SET amount = ? WHERE user_id = ? AND product_id = ? AND size = ?",
			amount+1,
			userId,
			productId,
			size,
		)
	}

	if err := query.Exec(); err != nil {
		return 0, &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}
	return amount + 1, nil
}

func (c *Cassandra) IncreaseProduct(userId string, tid string, mid string, size string) (uint8, error) {
	utils.Assert(userId != "", "user id is empty")
	utils.Assert(tid != "", "thread id is empty")
	utils.Assert(mid != "", "mid is empty")
	utils.Assert(size != "", "size is empty")

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
	utils.Assert(userId != "", "user id is empty")
	utils.Assert(tid != "", "thread id is empty")
	utils.Assert(mid != "", "mid is empty")
	utils.Assert(size != "", "size is empty")

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
	utils.Assert(userId != "", "user id is empty")
	utils.Assert(tid != "", "thread id is empty")
	utils.Assert(mid != "", "mid is empty")
	utils.Assert(size != "", "size is empty")

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

func (c *Cassandra) AddAddress(userId string, address models.Address, isDefault bool) error {
	utils.Assert(userId != "", "user id is empty")
	utils.Assert(address.Contact != "", "address contact is empty")
	utils.Assert(address.CountryCode != "", "address country code is empty")
	utils.Assert(address.Phone != "", "address phone is empty")
	utils.Assert(address.Country != "", "address country is empty")
	utils.Assert(address.Street != "", "address street is empty")
	utils.Assert(address.Region != "", "address region is empty")
	utils.Assert(address.City != "", "address city is empty")
	utils.Assert(address.Zipcode != "", "address zipcode is empt")

	cqlAddress := map[string]any{
		"address_id":   address.AddressId,
		"contact":      address.Contact,
		"country_code": address.CountryCode,
		"phone":        address.Phone,
		"country":      address.Country,
		"street":       address.Street,
		"region":       address.Region,
		"city":         address.City,
		"zipcode":      address.Zipcode,
	}

	var query *gocql.Query
	if isDefault {

	} else {
		query = c.session.Query(
			"UPDATE users SET addresses = addresses + ? WHERE user_id = ?",
			[]map[string]any{cqlAddress},
			userId,
		)
	}

	if err := query.Exec(); err != nil {
		// todo: way to handle not found
		return &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}

	return nil
}

func (c *Cassandra) Close() {
	c.session.Close()
}
