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
	Admins                []string
}

func (a *AppConfig) IsAdmin(telegramName string) bool {
	if len(a.Admins) == 0 {
		return false
	}

	for _, admin := range a.Admins {
		if admin == telegramName {
			return true
		}
	}
	
	return false
}