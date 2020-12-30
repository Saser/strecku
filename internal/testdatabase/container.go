package testdatabase

import (
	"net/url"

	"github.com/ory/dockertest/v3"
)

const (
	user     = "strecku"
	password = "password"
	dbName   = "strecku"
)

type container struct {
	res *dockertest.Resource
}

func newContainer() (*container, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}
	res, err := pool.Run("postgres", "13.1", []string{
		"POSTGRES_USER=" + user,
		"POSTGRES_PASSWORD=" + password,
		"POSTGRES_DB=" + dbName,
	})
	if err != nil {
		return nil, err
	}
	return &container{
		res: res,
	}, nil
}

func (c *container) Cleanup() error {
	return c.res.Close()
}

func (c *container) ConnString() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   c.res.GetHostPort("5432/tcp"),
		Path:   dbName,
	}
	q := url.Values{}
	q.Add("sslmode", "disable")
	u.RawQuery = q.Encode()
	return u.String()
}
