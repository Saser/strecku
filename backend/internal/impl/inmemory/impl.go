package inmemory

import (
	"sync"

	streckuv1 "github.com/Saser/strecku/backend/gen/api/strecku/v1"
)

type Impl struct {
	mu       sync.Mutex
	users    map[string]*streckuv1.User
	stores   map[string]*streckuv1.Store
	roles    map[string]*streckuv1.Role
	products map[string]*streckuv1.Product
}

func New() *Impl {
	return &Impl{
		users:    make(map[string]*streckuv1.User),
		stores:   make(map[string]*streckuv1.Store),
		roles:    make(map[string]*streckuv1.Role),
		products: make(map[string]*streckuv1.Product),
	}
}