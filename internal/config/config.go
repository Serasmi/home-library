package config

type AppConfig struct {
	Host string
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
}

type Config struct {
	App AppConfig
	DB  DBConfig
}
