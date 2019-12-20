package order

import (
	"github.com/tendermint/dex-demo/storeutil"
	dbm "github.com/tendermint/tm-db"
	"math/big"

	"github.com/tendermint/dex-demo/embedded/store"
	"github.com/tendermint/dex-demo/types"
	"github.com/tendermint/dex-demo/types/errs"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TableKey = "order_meta"

	orderPrefix      = "order"
	openOrderPrefix  = "open_order"
	ownedOrderPrefix = "owned_order"
)

type IteratorCB func(order Order) bool

type Keeper struct {
	as  store.ArchiveStore
	cdc *codec.Codec
}

func NewKeeper(db dbm.DB, cdc *codec.Codec) Keeper {
	return Keeper{
		as:  store.NewTable(db, TableKey),
		cdc: cdc,
	}
}

func (k Keeper) OpenOrdersByMarket(mktID sdk.Uint) []Order {
	var out []Order
	k.ReverseIteratorOpenOrders(func(order Order) bool {
		if !mktID.Equal(order.MarketID) {
			return true
		}

		out = append(out, order)
		return true
	})
	return out
}

func (k Keeper) OrdersByOwner(owner sdk.AccAddress, cb IteratorCB) {
	var ownedOrders []sdk.Uint

	k.as.ReversePrefixIterator(ownerOrderIterKey(owner), func(_ []byte, v []byte) bool {
		id := sdk.NewUintFromBigInt(new(big.Int).SetBytes(v))
		ownedOrders = append(ownedOrders, id)
		return true
	})

	for _, id := range ownedOrders {
		order, err := k.Get(id)
		if err != nil {
			continue
		}

		if !cb(order) {
			return
		}
	}
}

func (k Keeper) OnOrderCreatedEvent(event types.OrderCreated) {
	order := Order{
		ID:             event.ID,
		Owner:          event.Owner,
		MarketID:       event.MarketID,
		Direction:      event.Direction,
		Price:          event.Price,
		Quantity:       event.Quantity,
		Status:         "OPEN",
		Type:           "LIMIT",
		TimeInForce:    event.TimeInForceBlocks,
		QuantityFilled: sdk.ZeroUint(),
		CreatedBlock:   event.CreatedBlock,
	}
	k.Set(order)
	k.as.Set(ownerOrderKey(order.Owner, order.ID), storeutil.SDKUintSubkey(order.ID))
}

func (k Keeper) OnFillEvent(event types.Fill) sdk.Error {
	order, err := k.Get(event.OrderID)
	if err != nil {
		return err
	}

	order.QuantityFilled = order.QuantityFilled.Add(event.QtyFilled)
	if order.Quantity.Equal(order.QuantityFilled) {
		order.Status = "FILLED"
	}

	k.Set(order)
	return nil
}

func (k Keeper) OnOrderCancelledEvent(event types.OrderCancelled) sdk.Error {
	order, err := k.Get(event.OrderID)
	if err != nil {
		return err
	}

	order.Status = "CANCELLED"
	k.Set(order)
	return nil
}

func (k Keeper) Get(id sdk.Uint) (Order, sdk.Error) {
	var order Order
	ordB := k.as.Get(orderKey(id))
	if ordB == nil {
		return order, errs.ErrNotFound("order not found")
	}
	k.cdc.MustUnmarshalBinaryBare(ordB, &order)
	return order, nil
}

func (k Keeper) Set(order Order) {
	ordB := k.cdc.MustMarshalBinaryBare(order)
	k.as.Set(orderKey(order.ID), ordB)

	if order.Status == "OPEN" {
		k.as.Set(openOrderKey(order.MarketID, order.ID), storeutil.SDKUintSubkey(order.ID))
	} else {
		k.as.Delete(openOrderKey(order.MarketID, order.ID))
	}
}

func (k Keeper) ReverseIterator(cb IteratorCB) {
	k.as.ReversePrefixIterator([]byte(orderPrefix), func(_ []byte, v []byte) bool {
		var order Order
		k.cdc.MustUnmarshalBinaryBare(v, &order)
		return cb(order)
	})
}

func (k Keeper) ReverseIteratorOpenOrders(cb IteratorCB) {
	var openOrderIDs []sdk.Uint

	k.as.ReversePrefixIterator([]byte(openOrderPrefix), func(_ []byte, v []byte) bool {
		id := sdk.NewUintFromBigInt(new(big.Int).SetBytes(v))
		openOrderIDs = append(openOrderIDs, id)
		return true
	})

	for _, id := range openOrderIDs {
		order, err := k.Get(id)
		if err != nil {
			continue
		}

		if !cb(order) {
			return
		}
	}
}

func (k Keeper) ReverseIteratorFrom(startID sdk.Uint, cb IteratorCB) {
	// Inc() below because end is exclusive
	k.as.ReverseIterator(orderKey(sdk.ZeroUint()), orderKey(startID.Add(sdk.OneUint())), func(_ []byte, v []byte) bool {
		var order Order
		k.cdc.MustUnmarshalBinaryBare(v, &order)
		return cb(order)
	})
}

func (k Keeper) OnEvent(event interface{}) error {
	switch ev := event.(type) {
	case types.OrderCreated:
		k.OnOrderCreatedEvent(ev)
	case types.OrderCancelled:
		return k.OnOrderCancelledEvent(ev)
	case types.Fill:
		return k.OnFillEvent(ev)
	}

	return nil
}

func orderKey(id sdk.Uint) []byte {
	return storeutil.PrefixKeyString(orderPrefix, storeutil.SDKUintSubkey(id))
}

func openOrderKey(marketID sdk.Uint, orderID sdk.Uint) []byte {
	return storeutil.PrefixKeyString(openOrderPrefix, storeutil.SDKUintSubkey(marketID), storeutil.SDKUintSubkey(orderID))
}

func ownerOrderKey(owner sdk.AccAddress, orderID sdk.Uint) []byte {
	return storeutil.PrefixKeyBytes(ownerOrderIterKey(owner), storeutil.SDKUintSubkey(orderID))
}

func ownerOrderIterKey(owner sdk.AccAddress) []byte {
	return storeutil.PrefixKeyString(ownedOrderPrefix, owner.Bytes())
}
