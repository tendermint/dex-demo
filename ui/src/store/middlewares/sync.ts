import {ActionType} from "../../ducks/types";
import {Dispatch, Store} from "redux";
import {REDUX_STATE} from "../../ducks";
import {
  addBatchByMarketId,
  fetchBatchByMarketId,
  SET_CHART_INTERVAL,
  setOpenAsksByMarketId,
  setOpenBidsByMarketId,
} from "../../ducks/exchange";
import {get, OrderbookResponse} from "../../utils/fetch";
import {fetchBalance, fetchUserOrders} from "../../ducks/user";
import {ThunkDispatch} from "redux-thunk";

let watchDepthTimeout: any;
let watchOrderHistoryTimeout: any;
let watchWithdrawalsTimeout: any;
let watchBalanceTimeout: any;
let watchDepositStatusTimeout: any;
let watchBatchTimeout: any;

const syncMiddleware = (store: Store) => {
  return (next: Dispatch) => (action: ActionType<any>) => {
    const { dispatch, getState } = store;

    next(action);

    switch(action.type) {
      case '%INIT':
        handleFetchBook(getState, action, dispatch);
        handleOrderHistory(getState, action, dispatch);
        handleBalance(getState, action, dispatch);
        handleBatch(getState, action, dispatch);
        return;
      case SET_CHART_INTERVAL:
        return;
      default:
        return;
    }
  }
}

export default syncMiddleware;

function handleFetchBook(getState: () => REDUX_STATE, action: ActionType<any>, dispatch: Dispatch) {
  if (!watchDepthTimeout) {
    watchDepthTimeout = setTimeout(getDepth, 0);
  }

  const { exchange: { selectedMarket } } = getState();

  async function getDepth() {
    try {
      const resp = await get(`/markets/${1}/book`);
      const json: OrderbookResponse = await resp.json();

      dispatch(setOpenBidsByMarketId(json.bids, selectedMarket));
      dispatch(setOpenAsksByMarketId(json.asks, selectedMarket));

      watchDepthTimeout = setTimeout(getDepth, 500);
    } catch (e) {
      watchDepthTimeout = setTimeout(getDepth, 500);
    }

  }
}

function handleBatch(getState: () => REDUX_STATE, action: ActionType<any>, dispatch: ThunkDispatch<REDUX_STATE, any, ActionType<any>>) {
  if (!watchBatchTimeout) {
    watchBatchTimeout = setTimeout(getBatch, 0);
  }

  function getBatch() {
    const {
      exchange: { selectedMarket },
    } = getState();

    if (selectedMarket) {
      dispatch(fetchBatchByMarketId(selectedMarket));
    }
    watchBatchTimeout = setTimeout(getBatch, 2000);
  }
}

function handleOrderHistory(getState: () => REDUX_STATE, action: ActionType<any>, dispatch: ThunkDispatch<REDUX_STATE, any, ActionType<any>>) {
  if (!watchOrderHistoryTimeout) {
    watchOrderHistoryTimeout = setTimeout(getOrderHistory, 0);
  }

  async function getOrderHistory() {
    await dispatch(fetchUserOrders());
    watchOrderHistoryTimeout = setTimeout(getOrderHistory, 2000);
  }
}

function handleBalance(getState: () => REDUX_STATE, action: ActionType<any>, dispatch: ThunkDispatch<REDUX_STATE, any, ActionType<any>>) {
  if (!watchBalanceTimeout) {
    watchBalanceTimeout = setTimeout(getDaily, 0);
  }

  async function getDaily() {
    await dispatch(fetchBalance());
    watchBalanceTimeout = setTimeout(getDaily, 2000);
  }
}
