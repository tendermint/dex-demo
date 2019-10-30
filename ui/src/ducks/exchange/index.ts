import { ActionType } from '../types';
import {BatchesResponse, get, OrderJSON, PlaceOrderRequest, post} from "../../utils/fetch";
import { BigNumber as BN } from 'bignumber.js';
import {bn} from "../../utils/bn";
import {Dispatch} from "redux";

// Constants
export const SET_CHART_INTERVAL = 'app/exchange/setChartInterval';
export const SET_SPREAD_TYPE = 'app/exchange/setSpreadType';
export const SET_CHART_TYPE = 'app/exchange/setChartType';
export const SET_OPEN_BIDS_BY_MARKET_ID = 'app/exchange/setOpenBidsByMarketId';
export const SET_OPEN_ASKS_BY_MARKET_ID = 'app/exchange/setOpenAsksByMarketId';
export const SELECT_BATCH = 'app/exchange/selectBatch';
export const ADD_BATCH_BY_MARKET_ID = 'app/exchange/addBatchByMarketId';
export const ADD_ORDERS = 'app/exchange/addOrders';
export const UPDATE_DAILY_STATS_BY_MARKET_ID = 'app/exchange/updateDailyStatsByMarketId';

export enum INTERVAL {
  '1m' = '1',
  '5m' = '5',
  '15m' = '15',
  '1h' = '60',
  '1D' = 'D',
}

export enum SPREAD_TYPE {
  'ALL',
  'BID',
  'ASK',
}

export enum ORDER_SIDE {
  'buy',
  'sell',
}

export enum CHART_TYPE {
  'TradingView',
  'Depth',
  'Batch',
}

export type DayStatsType = {
  dayChange: BN
  dayChangePercentage: number
  dayHigh: BN
  dayLow: BN
  dayVolume: BN
  lastPrice: BN
  prevPrice: BN
}

export type UpdateDayStatsType = {
  dayChange?: BN
  dayChangePercentage?: number
  dayHigh?: BN
  dayLow?: BN
  dayVolume?: BN
  lastPrice?: BN
  prevPrice?: BN
}

export type BatchType = {
  date: string
  blockId: string
  marketId: string
  asks: Order[]
  bids: Order[]
  clearingPrice: BN
  clearingQuantity?: BN
}

export type MarketType = {
  quoteSymbol: string
  baseSymbol: string
  dayStats: DayStatsType
  bids: Order[]
  asks: Order[]
  batches: {
    [blockId: string]: BatchType
  }
}

export type Order = {
  marketId: string
  price: BN
  quantity: BN
}

export type OrderType = {
  created_block: string
  id: string
  market_id: string
  direction: 'BID' | 'ASK'
  owner: string
  price: BN
  quantity: BN
  quantity_filled: BN
  status: 'OPEN' | 'CANCELLED' | 'FILLED'
  time_in_force: number
}

export type ExchangeStateType = {
  selectedMarket: string
  selectedInterval: INTERVAL
  selectedSpreadType: SPREAD_TYPE
  selectedChartType: CHART_TYPE
  selectedBatch: string
  pairToMarketId: {
    [marketPair: string]: string
  },
  markets: {
    [key: string]: MarketType
  },
  orders: {
    [key: string]: OrderType
  },
}

const initialState = {
  selectedMarket: '1',
  selectedInterval: INTERVAL["1m"],
  selectedSpreadType: SPREAD_TYPE.ALL,
  selectedChartType: CHART_TYPE.TradingView,
  selectedBatch: '',
  pairToMarketId: {
    'DEMO/TEST': '1',
  },
  markets: {
    1: {
      quoteSymbol: 'DEMO',
      baseSymbol: 'TEST',
      dayStats: makeDayStats({}),
      batches: {},
      bids: [],
      asks: [],
    },
  },
  orders: {},
};

export const setChartInterval = (interval: INTERVAL): ActionType<INTERVAL> => ({
  type: SET_CHART_INTERVAL,
  payload: interval,
});

export const setChartType = (type: CHART_TYPE): ActionType<CHART_TYPE> => ({
  type: SET_CHART_TYPE,
  payload: type,
});

export const setSpreadType = (type: SPREAD_TYPE): ActionType<SPREAD_TYPE> => ({
  type: SET_SPREAD_TYPE,
  payload: type,
});

export const updateDailyStatsByMarketId = (extended: UpdateDayStatsType, marketId: string): ActionType<UpdateDailyActionPayload> => ({
  type: UPDATE_DAILY_STATS_BY_MARKET_ID,
  payload: {
    extended,
    marketId,
  },
});

type SetOpenBidsByMarketIdPayload = {
  marketId: string
  bids: OrderJSON[]
}
export const setOpenBidsByMarketId = (bids: OrderJSON[], marketId: string): ActionType<SetOpenBidsByMarketIdPayload> => ({
  type: SET_OPEN_BIDS_BY_MARKET_ID,
  payload: {
    marketId,
    bids,
  },
});

type SetOpenAsksByMarketIdPayload = {
  marketId: string
  asks: OrderJSON[]
}
export const setOpenAsksByMarketId = (asks: OrderJSON[], marketId: string): ActionType<SetOpenAsksByMarketIdPayload> => ({
  type: SET_OPEN_ASKS_BY_MARKET_ID,
  payload: {
    marketId,
    asks,
  },
});

export const addBatchByMarketId = (batch: BatchType): ActionType<BatchType> => ({
  type: ADD_BATCH_BY_MARKET_ID,
  payload: batch,
});

export const fetchBatchByMarketId = (marketId: string, blockId?: string) => async (dispatch: Dispatch) => {
  const resp = await get(`/markets/${marketId}/batches${blockId ? '/' + blockId : ''}`);
  if (resp.status === 404) {
    return;
  }
  const json: BatchesResponse = await resp.json();

  const batch: BatchType = {
    marketId,
    clearingPrice: bn(json.clearing_price),
    bids: json.bids.map(([ price, quantity ]) => ({
      marketId,
      price: bn(price),
      quantity: bn(quantity),
    })),
    asks: json.asks.map(([ price, quantity ]) => ({
      marketId,
      price: bn(price),
      quantity: bn(quantity),
    })),
    date: json.block_time,
    blockId: json.block_number,
  };

  dispatch(addBatchByMarketId(batch));
};

export const placeOrder = (r: PlaceOrderRequest) => (): Promise<Response> => {
  return post('/exchange/orders', r);
};

export const addOrders = (orders: OrderType[]): ActionType<OrderType[]> => ({
  type: ADD_ORDERS,
  payload: orders,
});

export const selectBatch = (blockId: string): ActionType<string> => ({
  type: SELECT_BATCH,
  payload: blockId,
});

export default function exchangeReducer(state: ExchangeStateType = initialState, action: ActionType<any>): ExchangeStateType {
  switch (action.type) {
    case SET_CHART_INTERVAL:
      return handleSetChartInterval(state, action);
    case SET_SPREAD_TYPE:
      return {
        ...state,
        selectedSpreadType: action.payload,
      };
    case SET_CHART_TYPE:
      return {
        ...state,
        selectedChartType: action.payload,
      };
    case SELECT_BATCH:
      return {
        ...state,
        selectedBatch: action.payload,
        selectedChartType: action.payload
          ? CHART_TYPE.Batch
          : CHART_TYPE.TradingView,
      };
    case SET_OPEN_BIDS_BY_MARKET_ID:
      return handleSetOpenBidsByMarketId(state, action);
    case SET_OPEN_ASKS_BY_MARKET_ID:
      return handleSetOpenAsksByMarketId(state, action);
    case ADD_BATCH_BY_MARKET_ID:
      return handleAddBatchByMarketId(state, action);
    case ADD_ORDERS:
      return handleAddOrders(state, action);
    case UPDATE_DAILY_STATS_BY_MARKET_ID:
      return handleUpdateDailyStatsByMarketId(state, action);
    default:
      return state;
  }
}

function handleSetChartInterval(state: ExchangeStateType, action: ActionType<INTERVAL>): ExchangeStateType {
  const { payload } = action;
  const { markets, selectedMarket } = state;
  const market = markets[selectedMarket];

  if (!market) return state;

  return {
    ...state,
    selectedInterval: payload,
  }
}

function handleSetOpenBidsByMarketId(state: ExchangeStateType, action: ActionType<SetOpenBidsByMarketIdPayload>): ExchangeStateType {
  const payload = action.payload;
  const marketId = payload.marketId;
  const { markets } = state;
  const market = markets[marketId];

  if (!market) return state;

  return {
    ...state,
    markets: {
      ...markets,
      [marketId]: {
        ...market,
        bids: payload.bids.map((b) => ({
          marketId: marketId,
          price: bn(b.price),
          quantity: bn(b.quantity),
        })),
      }
    },
  }
}

function handleSetOpenAsksByMarketId(state: ExchangeStateType, action: ActionType<SetOpenAsksByMarketIdPayload>): ExchangeStateType {
  const payload = action.payload;
  const marketId = payload.marketId;
  const { markets } = state;
  const market = markets[marketId];

  if (!market) return state;

  return {
    ...state,
    markets: {
      ...markets,
      [marketId]: {
        ...market,
        asks: payload.asks.map((a) => ({
          marketId: marketId,
          price: bn(a.price),
          quantity: bn(a.quantity),
        })),
      },
    },
  }
}

const MAX_BATCHES = 50;

function handleAddBatchByMarketId(state: ExchangeStateType, action: ActionType<BatchType>): ExchangeStateType {
  const payload = action.payload;
  const marketId = payload.marketId;
  const { markets } = state;
  const market = markets[marketId];

  if (!market) return state;

  const batchesClone = { ...market.batches };
  batchesClone[payload.blockId] = payload;
  const blockIds = Object.keys(batchesClone).sort();
  const slicedBatches = blockIds.slice(Math.max(blockIds.length - MAX_BATCHES, 0)).reduce((acc: {[blockId: string]: BatchType}, curr) => {
    acc[curr] = batchesClone[curr];
    return acc;
  }, {});

  return {
    ...state,
    markets: {
      ...markets,
      [marketId]: {
        ...market,
        batches: slicedBatches,
      }
    },
  }
}

function handleAddOrders(state: ExchangeStateType, action: ActionType<OrderType[]>): ExchangeStateType {
  const payload = action.payload;

  return {
    ...state,
    orders: payload.reduce((orders: {[k: string]: OrderType}, order) => {
      orders[order.id] = {
        ...order,
        price: bn(order.price),
        quantity: bn(order.quantity),
        quantity_filled: bn(order.quantity_filled),
      };
      return orders;
    }, {}),
  };
}

type UpdateDailyActionPayload = {
  extended: UpdateDayStatsType
  marketId: string
}

function handleUpdateDailyStatsByMarketId(state: ExchangeStateType, action: ActionType<UpdateDailyActionPayload>): ExchangeStateType {
  const {
    marketId,
    extended,
  } = action.payload;
  const {
    markets: {
      [marketId]: market,
    }
  } = state;

  if (!market) return state;

  return {
    ...state,
    markets: {
      ...state.markets,
      [marketId]: {
        ...market,
        dayStats: {
          ...market.dayStats,
          ...extended,
        },
      },
    },
  };
}

function makeDayStats (extended: UpdateDayStatsType): DayStatsType {
  return {
    dayChange: bn(extended.dayChange || 0),
    dayChangePercentage: extended.dayChangePercentage || 0,
    dayHigh: bn(extended.dayHigh || 0),
    dayLow: bn(extended.dayLow || 0),
    dayVolume: bn(extended.dayVolume || 0),
    lastPrice: bn(extended.lastPrice || 0),
    prevPrice: bn(extended.prevPrice || 0),
  }
}
