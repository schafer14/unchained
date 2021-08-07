package handle

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

//@unchained:routingKey /assets
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

//@unchained:routingKey
//@unchained:method Post
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

//@unchained:routingKey
//@unchained:method Get
func (a *AssetsHandler) HandleQuery(ctx context.Context, in DTOIn) (string, error) {
	a.Log.Warnln(in)

	return "hey", nil
}

//@unchained:routingKey /thing/{id}
func (a *AssetsHandler) HandleHi(ctx context.Context, id string, in DTOIn) (string, error) {
	a.Log.Warnln(in)
	a.Log.Info(id)

	return id, nil
}

//@unchained:routingKey /thing/{id}/{hey}
func (a *AssetsHandler) HandleHiy(ctx context.Context, hey, id string, in DTOIn) (string, error) {
	a.Log.Warnln(in)
	a.Log.Info(id)
	a.Log.Info(hey)

	return id, nil
}

type GroupsHandler struct {
	dig.In

	Log logrus.FieldLogger
}

func (a *GroupsHandler) HandleMember(ctx context.Context) (string, error) {
	a.Log.Warnln("no args")

	return "Banner", nil
}
