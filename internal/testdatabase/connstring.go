package testdatabase

func ConnString() string {
	mu.Lock()
	defer mu.Unlock()
	check()
	return defaultContainer.ConnString()
}
