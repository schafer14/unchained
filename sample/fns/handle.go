package handle

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type AssetsHandler struct {
	dig.In

	Log logrus.FieldLogger
}

type DTOIn struct {
	Proto         int    `json:"proto"`
	Quality       int    `json:"quality"`
	WalletAddress string `json:"walletAddress"`
}

type DTOOut struct {
	AssetID       string `json:"assetId"`
	Proto         int    `json:"proto"`
	Quality       int    `json:"quality"`
	WalletAddress string `json:"walletAddress"`
}

func (a *AssetsHandler) HandleNew(ctx context.Context, in DTOIn) (DTOOut, error) {
	a.Log.Warnln(in)

	dtoOut := DTOOut{
		AssetID:       uuid.New().String(),
		Proto:         in.Proto,
		Quality:       in.Quality,
		WalletAddress: in.WalletAddress,
	}

	return dtoOut, nil
}
