import { Dispatch } from "redux";
import { ActionType } from '../types';
import {post, get, GetUserOrderResponse, BalanceResponse} from "../../utils/fetch";
import {addOrders, OrderType} from "../exchange";
import BigNumber from "bignumber.js";
import {bn} from "../../utils/bn";

export const ADD_USER_ORDERS = 'app/user/addUserOrders';
export const ADD_USER_ADDRESS = 'app/user/addUserAddress';
export const SET_ORDER_HISTORY_FILTER = 'app/user/setOrderHistoryFilter';
export const SET_BALANCE = 'app/user/setBalance';
export const SET_LOGIN = 'app/user/setLogin';

export enum ORDER_HISTORY_FILTERS {
  'ALL',
  'OPEN',
}

export type BalanceType = {
  assetId: string
  locked: BigNumber
  unlocked: BigNumber
}

export type UserStateType = {
  orderHistoryFilter: ORDER_HISTORY_FILTERS
  orders: string[]
  transactions: string[]
  balances: {
    [assetId: string]: BalanceType
  }
  address: string
  isLoggedIn?: boolean
}

const initialState = {
  orderHistoryFilter: ORDER_HISTORY_FILTERS.OPEN,
  orders: [],
  transactions: [],
  balances: {},
  address: '',
  isLoggedIn: undefined,
};

export const setOrderHistoryFilter = (filter: ORDER_HISTORY_FILTERS): ActionType<ORDER_HISTORY_FILTERS> => ({
  type: SET_ORDER_HISTORY_FILTER,
  payload: filter,
});

export const addUserOrders = (orders: OrderType[]): ActionType<OrderType[]> => ({
  type: ADD_USER_ORDERS,
  payload: orders,
});

export const login = (password: string) => async (dispatch: Dispatch): Promise<Response> => {
  const resp = await post('/auth/login', { username: 'dex-demo', password });

  if (resp.status === 204) {
    const addrRes = await get('/auth/me');
    const addrJSON: { address: string} = await addrRes.json();
    dispatch(setAddress(addrJSON.address));
    dispatch({ type: '%INIT' });
  }

  return resp;
};

export const setLogin = (payload: boolean): ActionType<boolean> => ({
  type: SET_LOGIN,
  payload,
});

export const setBalance  = (payload: BalanceType): ActionType<BalanceType> => ({
  type: SET_BALANCE,
  payload,
});

export const setAddress = (payload: string): ActionType<string> => ({
  type: ADD_USER_ADDRESS,
  payload,
});

export const checkLogin = () => async (dispatch: Dispatch) => {
  const resp = await get('/user/balances');

  switch (resp.status) {
    case 401:
      return dispatch(setLogin(false));
    case 200:
      const addrRes = await get('/auth/me');
      const addrJSON: { address: string} = await addrRes.json();
      dispatch({ type: '%INIT' });
      dispatch(setAddress(addrJSON.address));
      dispatch(setLogin(true));
  }
};

export const fetchUserOrders = () => async (dispatch: Dispatch<ActionType<OrderType[]>>) => {
  try {
    const resp = await get('/user/orders');
    const json: GetUserOrderResponse = await resp.json();
    dispatch(addOrders(json.orders || []));
    dispatch(addUserOrders(json.orders || []));
  } catch (e) {
    console.log(e);
  }
};

export const addUserAddress = (address: string) => ({
  type: ADD_USER_ADDRESS,
  payload: address,
});

export const fetchBalance = () => async (dispatch: Dispatch<ActionType<BalanceType>>) => {
  try {
    const resp = await get('/user/balances');
    const json: BalanceResponse = await resp.json();

    json.balances.forEach(balance => {
      dispatch(setBalance({
        assetId: balance.asset_id,
        locked: bn(balance.at_risk),
        unlocked: bn(balance.liquid),
      }))
    })
  } catch (e) {
    console.log(e);
  }
};

export default function userReducer(state: UserStateType = initialState, action: ActionType<any>): UserStateType {
  switch (action.type) {
    case SET_ORDER_HISTORY_FILTER:
      return handleSetOrderHistoryFilter(state, action);
    case ADD_USER_ORDERS:
      return handleAddOrders(state, action);
    case SET_BALANCE:
      return handleSetBalance(state, action);
    case ADD_USER_ADDRESS:
      return {
        ...state,
        address: action.payload,
      };
    case SET_LOGIN:
      return {
        ...state,
        isLoggedIn: action.payload,
      };
    case '%INIT':
      return {
        ...state,
        isLoggedIn: true,
      };
    default:
      return state;
  }
}

function handleAddOrders(state: UserStateType, action: ActionType<OrderType[]>): UserStateType {
  return {
    ...state,
    orders: action.payload.map(({ id }) => id),
  }
}

function handleSetOrderHistoryFilter(state: UserStateType, action: ActionType<ORDER_HISTORY_FILTERS>): UserStateType {
  return {
    ...state,
    orderHistoryFilter: action.payload,
  };
}

function handleSetBalance(state: UserStateType, action: ActionType<BalanceType>): UserStateType {
  const { assetId } = action.payload;
  return {
    ...state,
    balances: {
      ...state.balances,
      [assetId]: action.payload,
    },
  };
}