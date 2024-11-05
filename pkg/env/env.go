package env

// Mode represents the environment in which the application is running
type Mode string

const (
	Local Mode = "local"
	Test  Mode = "test"
	Dev   Mode = "dev"
	Prod  Mode = "prod"
)

func (e Mode) String() string {
	return string(e)
}

func (e Mode) Validate() bool {
	switch e {
	case Local, Test, Dev, Prod:
		return true
	default:
		return false
	}
}
