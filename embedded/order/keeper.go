package order

import (
	dbm "github.com/tendermint/tm-db"

	"github.com/tendermint/dex-demo/types"
	"github.com/tendermint/dex-demo/types/errs"
	"github.com/tendermint/dex-demo/types/store"

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

func (k Keeper) OpenOrdersByMarket(mktID store.EntityID) []Order {
	var out []Order
	k.ReverseIteratorOpenOrders(func(order Order) bool {
		if !mktID.Equals(order.MarketID) {
			return true
		}

		out = append(out, order)
		return true
	})
	return out
}

func (k Keeper) OrdersByOwner(owner sdk.AccAddress, cb IteratorCB) {
	var ownedOrders []store.EntityID

	k.as.ReversePrefixIterator(ownerOrderIterKey(owner), func(_ []byte, v []byte) bool {
		id := store.NewEntityIDFromBytes(v)
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
		QuantityFilled: sdk.NewUint(0),
		CreatedBlock:   event.CreatedBlock,
	}
	k.Set(order)
	k.as.Set(ownerOrderKey(order.Owner, order.ID), order.ID.Bytes())
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

func (k Keeper) Get(id store.EntityID) (Order, sdk.Error) {
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
		k.as.Set(openOrderKey(order.MarketID, order.ID), order.ID.Bytes())
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
	var openOrderIDs []store.EntityID

	k.as.ReversePrefixIterator([]byte(openOrderPrefix), func(_ []byte, v []byte) bool {
		id := store.NewEntityIDFromBytes(v)
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

func (k Keeper) ReverseIteratorFrom(startID store.EntityID, cb IteratorCB) {
	// Inc() below because end is exclusive
	k.as.ReverseIterator(orderKey(store.NewEntityID(0)), orderKey(startID.Inc()), func(_ []byte, v []byte) bool {
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

func orderKey(id store.EntityID) []byte {
	return store.PrefixKeyString(orderPrefix, id.Bytes())
}

func openOrderKey(marketID store.EntityID, orderID store.EntityID) []byte {
	return store.PrefixKeyString(openOrderPrefix, marketID.Bytes(), orderID.Bytes())
}

func ownerOrderKey(owner sdk.AccAddress, orderID store.EntityID) []byte {
	return store.PrefixKeyBytes(ownerOrderIterKey(owner), orderID.Bytes())
}

func ownerOrderIterKey(owner sdk.AccAddress) []byte {
	return store.PrefixKeyString(ownedOrderPrefix, owner.Bytes())
}
