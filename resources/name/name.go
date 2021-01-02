package name

type Name struct {
	segments []string
	indices  map[string]int // variable name -> index into segments
}
