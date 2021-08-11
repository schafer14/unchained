package config

// Conf ...
type Conf struct {
	Environment string `conf:"default:development"`

	Database struct {
		User     string   `conf:"default:root"`
		Password string   `conf:"default:openSesame,noprint"`
		Hosts    []string `conf:"default:http://0.0.0.0:8529"`
		Name     string   `conf:"default:unchained"`
	}
}
