package handle

import (
	"assets/config"
	"context"

	driver "github.com/arangodb/go-driver"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type AssetsHandler struct {
	dig.In

	Log  logrus.FieldLogger
	Conf config.Conf
	DB   driver.Collection
}

type DTOIn struct {
	Proto         int    `json:"proto"`
	Quality       int    `json:"quality"`
	WalletAddress string `json:"walletAddress"`
}

type DTOOut struct {
	AssetID       string `json:"_key"`
	Proto         int    `json:"proto"`
	Quality       int    `json:"quality"`
	WalletAddress string `json:"walletAddress"`
}

func (a *AssetsHandler) HandleNew(ctx context.Context, in DTOIn) (DTOOut, error) {
	a.Log.Infoln(a.Conf.Environment)

	dtoOut := DTOOut{
		AssetID:       uuid.New().String(),
		Proto:         in.Proto,
		Quality:       in.Quality,
		WalletAddress: in.WalletAddress,
	}

	return dtoOut, nil
}

func (a *AssetsHandler) HandleNewWithDB(ctx context.Context, in DTOIn) (DTOOut, error) {
	a.Log.Infoln(a.Conf.Environment)

	dtoOut := DTOOut{
		AssetID:       uuid.New().String(),
		Proto:         in.Proto,
		Quality:       in.Quality,
		WalletAddress: in.WalletAddress,
	}

	a.DB.CreateDocument(ctx, dtoOut)

	return dtoOut, nil
}

//@unchained:routingKey /example/{id}
func (a *AssetsHandler) HandleHi(ctx context.Context, id string) (DTOOut, error) {
	a.Log.Info(id)

	var resBody DTOOut
	_, err := a.DB.ReadDocument(ctx, id, &resBody)
	if err != nil {
		a.Log.Errorln("oh no", err)
		return resBody, err
	}

	return resBody, err
}

//@unchained:routingKey /example/{id}/{hey}
//@unchained:method POST
func (a *AssetsHandler) HandleHiy(ctx context.Context, hey, id string) (string, error) {
	a.Log.Info(id)
	a.Log.Info(hey)

	return id, nil
}
