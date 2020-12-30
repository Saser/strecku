package testdatabase

import "testing"

func TestNewContainer(t *testing.T) {
	c, err := newContainer()
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Cleanup(); err != nil {
		t.Error(err)
	}
}
