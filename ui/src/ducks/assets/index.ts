import {ActionType} from "../types";
import {FundingSources} from "../../constants/clearinghouse";

export type AssetType = {
  symbol: string
  name: string
  decimals: number
  nativeDecimals: number
  sources: FundingSources[]
  chainId: string
}

export type ChainType = {
  id: string
  name: string
  depositFinality: number
}

export type AssetStateType = {
  symbolToAssetId: {
    [k: string]: string
  }
  assets: {
    [k: string]: AssetType
  }
  chains: {
    [k: string]: ChainType
  }
}

const initialState: AssetStateType = {
  symbolToAssetId: {
    BTC: '3',
    DEMO: '1',
    TEST: '2',
    USD: '4',
  },
  assets: {
    '4': {
      symbol: 'USD',
      name: 'US Dollar',
      decimals: 18,
      nativeDecimals: 2,
      sources: [],
      chainId: 'USD',
    },
    '3': {
      symbol: 'BTC',
      name: 'Bitcoin',
      decimals: 18,
      nativeDecimals: 8,
      sources: [],
      chainId: 'BTC',
    },
    '1': {
      symbol: 'DEMO',
      name: 'DEX Demo Token',
      decimals: 18,
      nativeDecimals: 4,
      sources: [],
      chainId: '',
    },
    '2': {
      symbol: 'TEST',
      name: 'Test Token',
      decimals: 18,
      nativeDecimals: 4,
      sources: [],
      chainId: '',
    },
  },
  chains: {
    ETH: {
      id: 'ETH',
      name: 'Rinkeby - Ethereum Testnet',
      depositFinality: 1,
    },
  },
};

export default function assetReducer(state: AssetStateType = initialState, action: ActionType<any>) {
  switch (action.type) {
    default:
      return state;
  }
}