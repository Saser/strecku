package database

import (
	"net/url"
	"testing"

	"github.com/ory/dockertest/v3"
)

const (
	user     = "strecku"
	password = "password"
	dbName   = "strecku"
)

type Container struct {
	res *dockertest.Resource
}

func NewContainer() (*Container, error) {
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
	return &Container{
		res: res,
	}, nil
}

func NewContainerT(t *testing.T) *Container {
	c, err := NewContainer()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := c.Cleanup(); err != nil {
			t.Fatal(err)
		}
	})
	return c
}

func (c *Container) Cleanup() error {
	return c.res.Close()
}

func (c *Container) ConnString() string {
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
