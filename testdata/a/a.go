package a

const (
	// aa const
	aa = "aaaaaa"
)

// c cccc
const c = "cccc"

type person struct{}

// Say hello
func (p *person) Say() string {
	return "hello"
}

// Say hello
func (p person) Age() int {
	return 5
}

// Add godoc
// line 2
// line 3
func Add(m, n int) int {
	return m + n
}
