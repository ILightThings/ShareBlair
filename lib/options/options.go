package options

type UserFlags struct {
	Target   string
	Threads  int
	Verbose  bool
	User     string
	Domain   string
	Password string
	Hash     string
	Port     int
}
