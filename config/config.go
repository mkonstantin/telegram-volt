package config

type AppConfig struct {
	Username              string
	Password              string
	Host                  string
	Port                  int
	Database              string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifeTime int
	Version               string
}
