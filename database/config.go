package database

type Config struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}
