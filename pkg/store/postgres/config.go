package postgres

type config struct {
	driver   string
	host     string
	port     int
	user     string
	password string
	name     string
	ssl      string
}

func NewConfig(
	driver string,
	host string,
	port int,
	user string,
	password string,
	name string,
	ssl string,
) *config {
	if host == "" {
		host = "localhost"
	}

	if port == 0 {
		port = 5432
	}

	if ssl == "" {
		ssl = "disable"
	}

	return &config{
		driver:   driver,
		host:     host,
		port:     port,
		user:     user,
		password: password,
		name:     name,
		ssl:      ssl,
	}
}
