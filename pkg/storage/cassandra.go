package storage

import (
	"errors"
	"nedas/shop/pkg/models"

	"github.com/gocql/gocql"
)

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
		"INSERT INTO users(user_id, email) VALUES (?, ?) IF NOT EXISTS",
		user.UserID,
		user.Email,
	)

	ok, err := query.ScanCAS(nil, nil)

	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			err = ErrNotFound
		}
		return &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: err}
	}

	if !ok {
		return &StorageError{Provider: "CASSANDRA", Execution: query.Statement(), Err: ErrAlreadySet}
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

func (c *Cassandra) Close() {
	c.session.Close()
}

func assert(ok bool, msg string) {
	if !ok {
		panic(msg)
	}
}
