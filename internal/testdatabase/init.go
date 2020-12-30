package testdatabase

import "sync"

var (
	mu               sync.Mutex
	defaultContainer *container
)

func Init() {
	mu.Lock()
	defer mu.Unlock()
	if defaultContainer != nil {
		return
	}
	var err error
	defaultContainer, err = newContainer()
	if err != nil {
		panic(err)
	}
}

func Cleanup() {
	mu.Lock()
	defer mu.Unlock()
	if defaultContainer == nil {
		return
	}
	if err := defaultContainer.Cleanup(); err != nil {
		panic(err)
	}
	defaultContainer = nil
}
