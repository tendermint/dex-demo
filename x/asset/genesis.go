package asset

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/dex-demo/types/store"
	"github.com/tendermint/dex-demo/x/asset/types"
)

type GenesisState struct {
	Assets []types.Asset `json:"assets"`
}

func NewGenesisState(assets []types.Asset) GenesisState {
	return GenesisState{Assets: assets}
}

func ValidateGenesis(data GenesisState) error {
	currentId := store.ZeroEntityID

	for _, asset := range data.Assets {
		currentId = currentId.Inc()
		if !currentId.Equals(asset.ID) {
			return errors.New("Invalid Asset: ID must monotonically increase.")
		}
		if asset.Name == "" {
			return errors.New("Invalid Asset: Must specify a name.")
		}
		if asset.Symbol == "" {
			return errors.New("Invalid Asset: Must specify a symbol.")
		}
		if asset.TotalSupply.IsZero() {
			return errors.New("Invalid Asset: Must specify a non-zero total supply.")
		}
	}

	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Assets: []types.Asset{
			{
				ID:                store.NewEntityID(1),
				Name:              "UEX Staking Token",
				Symbol:            "UEX",
				CirculatingSupply: sdk.NewUintFromString("40000000000000000000000000"),
				TotalSupply:       sdk.NewUintFromString("1000000000000000000000000000"),
			},
			{
				ID:                store.NewEntityID(2),
				Name:              "Test Token",
				Symbol:            "TEST",
				CirculatingSupply: sdk.NewUintFromString("40000000000000000000000000"),
				TotalSupply:       sdk.NewUintFromString("1000000000000000000000000000"),
			},
		},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, asset := range data.Assets {
		keeper.Inject(ctx, asset)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var assets []types.Asset
	k.Iterator(ctx, func(asset types.Asset) bool {
		assets = append(assets, asset)
		return true
	})
	return GenesisState{Assets: assets}
}
