package inmemory

import (
	"sync"

	streckuv1 "github.com/Saser/strecku/backend/gen/api/v1"
)

type Impl struct {
	mu     sync.Mutex
	users  map[string]*streckuv1.User
	stores map[string]*streckuv1.Store
	roles  map[string]*streckuv1.Role
}

func New() *Impl {
	return &Impl{
		users:  make(map[string]*streckuv1.User),
		stores: make(map[string]*streckuv1.Store),
		roles:  make(map[string]*streckuv1.Role),
	}
}
