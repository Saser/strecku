package testdatabase

// check verifies that the container has been initialized. It is not
// thread-safe, and the caller must lock using mu before calling check.
func check() {
	if defaultContainer == nil {
		panic("no container -- Init() must be called before any operations involving the container can be done")
	}
}
