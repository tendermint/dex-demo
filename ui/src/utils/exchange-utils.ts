import {BatchType, Order} from "../ducks/exchange";
import BigNumber, { BigNumber as BN } from 'bignumber.js';
import {bn} from "./bn";
import store from "../store";

type ReduceDepthFromOrdersReturnTypes = {
  max: BN
  depths: Order[],
}
export const reduceDepthFromOrders = (orders: Order[], decimals: number, nativeDecimals: number): ReduceDepthFromOrdersReturnTypes => {
  let max = bn(0);
  const depths: {[p: string]: BN} = {};
  const { assets: { assets } } = store.getState();

  orders.forEach((order) => {
    const key = order.price.div(10 ** decimals).toFixed(Math.min(4, nativeDecimals));
    depths[key] = depths[key] || bn(0);
    depths[key] = depths[key].plus(order.quantity);
  });

  const ret = Object
    .entries(depths)
    .map(([ price, quantity ]) => {
      const bnPrice = bn(price);
      const total = bnPrice.multipliedBy(quantity);
      max = BN.max(max, total);
      return {
        marketId: orders[0].marketId,
        price: bnPrice,
        quantity: quantity,
      }
    });

  return {
    max,
    depths: ret,
  };
};

type EstimateBatchReturnTypes = {
  clearingPrice: BN
  bidRation: BN
  askRation: BN
}
export function estimateBatch(bids: Order[], asks: Order[], decimals: number, nativeDecimals: number): EstimateBatchReturnTypes {
  const { depths: bidDepth } = reduceDepthFromOrders(bids, decimals, nativeDecimals);
  const { depths: askDepth } = reduceDepthFromOrders(asks, decimals, nativeDecimals);

  const priceOrder: BN[] = bidDepth
    .map(({ price }) => price)
    .concat(askDepth.map(({ price }) => price))
    .sort((a: BN, b: BN) => a.minus(b).toNumber());

  let found = false;
  let clearingPrice = bn(0);
  let clearingPriceIndex = 0;

  priceOrder.forEach((price: BN, i: number) => {
    if (found) return;

    const ad = aggregateDemand(bidDepth, price);
    const as = aggregateSupply(askDepth, price);

    if (!as.isZero()) {
      if (!ad.isZero()) {
        clearingPrice = price;
        clearingPriceIndex = i;
      }
    }

    if (as.isGreaterThan(ad)) {
      found = true;
      const lastPrice = priceOrder[i - 1];

      if (lastPrice && !clearingPrice.isZero()) {
        const lastAs = aggregateSupply(askDepth, lastPrice);
        const lastAd = aggregateDemand(bidDepth, lastPrice);
        // const currentAs = aggregateSupply(askDepth, clearingPrice);
        // const currentAd = aggregateDemand(bidDepth, clearingPrice);

        if (!lastAs.isZero() && !lastAd.isZero()) {
          clearingPrice = lastPrice;
        }
      }
    }
  });

  let bidRation: BN;
  let askRation: BN;

  if (clearingPrice.isZero()) {
    bidRation = bn(0);
    askRation = bn(0);
  } else {
    const rAd = aggregateDemand(bidDepth, clearingPrice);
    const rAs = aggregateSupply(askDepth, clearingPrice);
    bidRation = rAd.isGreaterThanOrEqualTo(rAs)
      ? rAs.div(rAd)
      : bn(1);
    askRation = rAs.isGreaterThanOrEqualTo(rAd)
      ? rAd.div(rAs)
      : bn(1);
  }

  return {
    clearingPrice,
    bidRation,
    askRation,
  }
}

export function aggregateSupply(list: Order[], stopPrice: BN): BN {
  return list.reduce((acc: BN, order) => {
    return stopPrice.isGreaterThanOrEqualTo(order.price)
      ? acc.plus(order.quantity)
      : acc;
  }, bn(0))
}

export function aggregateDemand(list: Order[], stopPrice: BN): BN {
  return list.reduce((acc: BN, order) => {
    return stopPrice.isLessThanOrEqualTo(order.price)
      ? acc.plus(order.quantity)
      : acc;
  }, bn(0))
}

export function sortOrders(list: Order[], descending: boolean = false): Order[] {
  return list.concat().sort((a, b) => {
    if (a.price.isLessThan(b.price)) return descending ? 1 : -1;
    if (a.price.isGreaterThan(b.price)) return descending ? -1 : 1;
    return 0;
  });
}

export function findQuantityUnderPrice(orders: Order[], price: BigNumber): BigNumber {
  return orders.reduce((acc, order) => {
    if (order.price.isLessThanOrEqualTo(price)) {
      return acc.plus(order.quantity);
    }
    return acc;
  }, bn(0));
}

export function findQuantityOverPrice(orders: Order[], price: BigNumber): BigNumber {
  return orders.reduce((acc, order) => {
    if (order.price.isGreaterThanOrEqualTo(price)) {
      return acc.plus(order.quantity);
    }
    return acc;
  }, bn(0));
}