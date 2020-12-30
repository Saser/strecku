package database

import "testing"

func TestNewContainer(t *testing.T) {
	c, err := NewContainer()
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Cleanup(); err != nil {
		t.Error(err)
	}
}

func TestNewContainerT(t *testing.T) {
	_ = NewContainerT(t)
}
