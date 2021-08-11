package dependencies

import (
	"assets/config"
	"context"
	"fmt"
	"os"
	"time"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ProvideLogger ...
func ProvideLogger() logrus.FieldLogger {
	return logrus.New()
}

// ProvideConf ...
func ProvideConf() (config.Conf, error) {
	var cfg config.Conf

	if err := conf.Parse(os.Args[1:], "", &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage("", &cfg)
			if err != nil {
				return cfg, errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			os.Exit(0)
		case conf.ErrVersionWanted:
			version, err := conf.VersionString("", &cfg)
			if err != nil {
				return cfg, errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			os.Exit(0)
		}
		return cfg, errors.Wrap(err, "parsing config")
	}

	return cfg, nil
}

// ProvideDatabase ...
func ProvideDatabase(cfg config.Conf) (driver.Collection, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: cfg.Database.Hosts,
	})
	if err != nil {
		return nil, errors.Wrap(err, "opening database")
	}

	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(cfg.Database.User, cfg.Database.Password),
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating database client")
	}

	db, err := c.Database(ctx, cfg.Database.Name)
	if err != nil {
		return nil, errors.Wrap(err, "connecting to database")
	}

	return db.Collection(ctx, "assets")
}
