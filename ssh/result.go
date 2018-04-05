package ssh

// Result is the struct that represents row in the table
type Result struct {
	User string
	Host string
	Port string
	Err  error
}
