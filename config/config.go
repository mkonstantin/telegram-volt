package config

const (
	DBUsername            = "root"
	DBPassword            = "root"
	DBHost                = "localhost"
	DBPort                = 3310
	DBDatabase            = "volt"
	MaxOpenConnections    = 30
	MaxIdleConnections    = 60
	ConnectionMaxLifeTime = 2
)

type AppConfig struct {
	DBUsername            string
	DBPassword            string
	DBHost                string
	DBPort                int
	DBDatabase            string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifeTime int
}

func GetConfig() AppConfig {
	return AppConfig{
		DBUsername:            DBUsername,
		DBPassword:            DBPassword,
		DBHost:                DBHost,
		DBPort:                DBPort,
		DBDatabase:            DBDatabase,
		MaxOpenConnections:    MaxOpenConnections,
		MaxIdleConnections:    MaxIdleConnections,
		ConnectionMaxLifeTime: ConnectionMaxLifeTime,
	}
}
