package dockerssh

type Options struct {
	Name    string
	Port    string
	Volume  string
	Verbose bool
}

type Service struct {
	Name    string
	Connect string
	Volume  string
	State   string
}
