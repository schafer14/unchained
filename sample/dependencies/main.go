package dependencies

import (
	"github.com/sirupsen/logrus"
)

func ProvideLogger() logrus.FieldLogger {
	return logrus.New()
}
