import {OrderType} from '../ducks/exchange';

export const BASE_API: string = process.env.BASE_API || '/api/v1';

let CSRF_TOKEN: string = (window as any).CSRF_TOKEN;

const getCSRFToken = async (): Promise<string> => {
  if (CSRF_TOKEN) return CSRF_TOKEN;

  const resp = await fetch(`${BASE_API}/auth/csrf_token`);
  CSRF_TOKEN = await resp.text();
  return CSRF_TOKEN;
};

export const post = async (url: string, body: object): Promise<Response> => {
  const token = await getCSRFToken();
  return fetch(`${BASE_API}${url}`, {
    method: 'POST',
    headers: {
      'X-CSRF-Token': token,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(body),
  });
};

export const get = async (url: string): Promise<Response> => {
  // const token = await getCSRFToken();
  return fetch(`${BASE_API}${url}`);
};

export type OrderJSON = {
  price: string
  quantity: string
}

export type PlaceOrderRequest = {
  market_id: string
  direction: 'BID' | 'ASK'
  price: string
  quantity: string
  type: 'MARKET' | 'LIMIT'
  time_in_force: number
}

export type CandleStickJSON = {
  date: string
  open: string
  close: string
  high: string
  low: string
}

export type CandlesResponse = {
  pair: string
  market_id: string
  candles: CandleStickJSON[]
}

export type GetUserOrderResponse = {
  next_id: string
  orders: OrderType[]
}

export type OrderbookResponse = {
  asks: OrderJSON[]
  bids: OrderJSON[]
  block_number: string
  market_id: string
}

type Balance = {
  asset_id: string
  at_risk: string
  liquid: string
  name: string
  symbol: string
}

export type BalanceResponse = {
  balances: Balance[]
}

// asset_id: "2"
// beneficiary: "0x7f751422e3ffcae90ae74049ff8aa5f3bc47335d"
// burn_id: "4"
// initiated_block: 639
// merkle_leaf: "0x5c71919851268e73587da9f985c8d493fa52afcb2d094b2fc10dd3bfd627ef96"
// merkle_proof: null
// merkle_root: null
// owner: "cosmos1j689jv788xfhmvm27pgz0f7uvxjxz8tmuw2yqf"
// quantity: "10000000000000000"
export type GetWithdrawalResponse = {
  owner: string,
  withdrawals: WithdrawlResponse[] | null
}

export type WithdrawlResponse = {
  asset_id: string
  beneficiary: string
  burn_id: string
  initiated_block: string
  merkle_leaf: string
  merkle_proof: string | null
  merkle_root: string | null
  owner: string
  quantity: string
}
// asks: (13) [Array(2), Array(2), Array(2), Array(2), Array(2), Array(2), Array(2), Array(2), Array(2), Array(2), Array(2), Array(2), Array(2)]
// bids: (12) [Array(2), Array(2), Array(2), Array(2), Array(2), Array(2), Array(2), Array(2), Array(2), Array(2), Array(2), Array(2)]
// block_number: "49523"
// block_time: "2019-07-12T07:45:35.934611487Z"
// clearing_price: "273780000000000000000"
// market_id: "1"
export type BatchesResponse = {
  block_number: string
  block_time: string
  clearing_price: string
  market_id: string
  asks: [string, string][]
  bids: [string, string][]
}