import {bn} from "../bn";
import {LibrarySymbolInfo} from "../../../public/charting_library/charting_library.min";
import {CandlesResponse, CandleStickJSON, get} from "../fetch";
import store from "../../store";
import CandleCache from "./candle-cache";

type RawCandleType = {
  time: string,
  open: string,
  close: string,
  high: string,
  low: string,
  volume: string,
}

type TVCandleType = {
  time: number,
  open: number,
  close: number,
  high: number,
  low: number,
  volume: number,
}

type SubscriberType = {
  callback?: (tvCandle: TVCandleType) => void
  recentCandles: CandleCache
}

class ChartDataProvider {
  pulseTimeout: any | null = null;
  subscribers: {
    [subscribeUID: string]: SubscriberType
  } = {};

  async start () {
    if (!this.pulseTimeout) {
      await this.scanSubscribers();
    } else {
      console.log('Chart data pulse is already started.');
    }
  }

  stop () {
    if (this.pulseTimeout) {
      clearTimeout(this.pulseTimeout);
      this.pulseTimeout = null;
    }
  }

  getSubscriberAt (subscriberId: string): SubscriberType {
    const sub = this.subscribers[subscriberId];

    if (sub) return sub;

    const { baseSymbol, quoteSymbol } = deserializeUID(subscriberId);
    const { exchange: { pairToMarketId } } = store.getState();
    const marketId = pairToMarketId[`${baseSymbol}/${quoteSymbol}`];

    const newSub: SubscriberType = {
      recentCandles: new CandleCache({ marketId }),
    };

    this.subscribers[subscriberId] = newSub;

    return newSub;
  }

  async fetchCandles (marketId: string, baseSymbol: string, quoteSymbol: string, resolution: string, limit: number = 2000, toTimestamp: number, fromTimestamp?: number): Promise<RawCandleType[]> {
    const json: CandlesResponse = await fetchCandlesFromUEX(marketId, toTimestamp, fromTimestamp);
    // const json: CandlesResponse = await fetchCandlesFromCryptoCompare(baseSymbol, quoteSymbol, toTimestamp);

    if (!json || !json.candles || !json.candles.length) {
      return [];
    }

    const UID = `${baseSymbol}/${quoteSymbol}_${resolution}`;
    const sub = this.getSubscriberAt(UID);

    return json.candles.map((bar: CandleStickJSON) => {
      sub.recentCandles.add(bar);

      return {
        time: bar.date,
        open: String(bar.open),
        close: String(bar.close),
        high: String(bar.high),
        low: String(bar.low),
        volume: '0',
      };
    });
  }

  scanSubscribers = async () => {
    const promises = Object
      .entries(this.subscribers)
      .map(async ([ subscribeUID, sub ]) => {
        const { baseSymbol, quoteSymbol, resolution } = deserializeUID(subscribeUID);
        const {
          exchange: { pairToMarketId },
          assets: { assets, symbolToAssetId },
        } = store.getState();
        const marketId = pairToMarketId[`${baseSymbol}/${quoteSymbol}`];
        const asset = assets[symbolToAssetId[quoteSymbol]];

        if (!asset) return;

        const currentTimestamp = new Date().getTime();
        const last = sub.recentCandles.last();
        const lastTimestamp = last && new Date(last.date).getTime();

        const rawCandles = await this.fetchCandles(
          marketId,
          baseSymbol,
          quoteSymbol,
          resolution,
          2,
          currentTimestamp,
          lastTimestamp,
        );

        const tvCandles = formatTVCandles(rawCandles, 1 / (10 ** asset.decimals));

        tvCandles.forEach(candle => {
          if (lastTimestamp <= candle.time) {
            if (sub.callback) sub.callback(candle);
          }
        });
      });

    await Promise.all(promises);
    this.pulseTimeout = setTimeout(this.scanSubscribers, 1000);
  };

  subscribe (subscribeUID: string, callback: (candle: TVCandleType) => void) {
    this.subscribers[subscribeUID] = this.getSubscriberAt(subscribeUID);
    this.subscribers[subscribeUID].callback = callback;
    this.start();
  }

  unsubscribe (subscribeUID: string) {
    if (this.subscribers[subscribeUID]) {
      delete this.subscribers[subscribeUID];
    }

    if (!Object.keys(this.subscribers).length) {
      this.stop();
    }
  }
}

const chartDataProvider = new ChartDataProvider();

export default chartDataProvider;

export function formatTVCandles (rawData: RawCandleType[], mutiplier: number = 1): TVCandleType[] {
  return rawData.map(raw => {
    return {
      time: new Date(raw.time).getTime(),
      open: bn(raw.open).multipliedBy(mutiplier).toNumber(),
      close: bn(raw.close).multipliedBy(mutiplier).toNumber(),
      high: bn(raw.high).multipliedBy(mutiplier).toNumber(),
      low: bn(raw.low).multipliedBy(mutiplier).toNumber(),
      volume: bn(raw.volume).multipliedBy(mutiplier).toNumber(),
    };
  });
}

export function serializeUID (symbolInfo: LibrarySymbolInfo, resolution: string): string {
  const { name } = symbolInfo;
  return `${name}_${resolution}`;
}

export function deserializeUID (subscriberUID: string): { baseSymbol: string, quoteSymbol: string, resolution: string} {
  const [ name, resolution ] = subscriberUID.split('_');
  const [ baseSymbol, quoteSymbol ] = name.split('/');
  return { baseSymbol, quoteSymbol, resolution };
}

async function fetchCandlesFromUEX (marketId: string, toTimestamp: number, fromTimestamp?: number): Promise<CandlesResponse> {
  const iso = fromTimestamp
    ? new Date(fromTimestamp).toISOString()
    : new Date(toTimestamp).toISOString();
  const resp = fromTimestamp
    ? await get(`/markets/${marketId}/candles?granularity=1m&start=${iso}&end=${new Date(toTimestamp).toISOString()}`)
    : await get(`/markets/${marketId}/candles?granularity=1m&end=${iso}`);
  return await resp.json() as CandlesResponse;
}

async function fetchCandlesFromCryptoCompare (baseSymbol: string, quoteSymbol: string, toTimestamp: number): Promise<CandlesResponse> {
  const resp = await fetch(`https://min-api.cryptocompare.com/data/histominute?fsym=${baseSymbol}&tsym=${quoteSymbol}&toTs=${toTimestamp / 1000}&limit=2000&api_key=9d2c1563297d190bf26d5ac5d0e375c0f3441689d30e6e78d11533ace062f363`);
  const raw: {
    Response: string
    Data: {
      close: number
      high: number
      low: number
      open: number
      time: number
      volumefrom: number
      volumeto: number
    }[]
  } = await resp.json();

  if (raw.Response !== 'Success') {
    return {
      pair: `${baseSymbol}/${quoteSymbol}`,
      market_id: '1',
      candles: [],
    };
  }

  return {
    pair: `${baseSymbol}/${quoteSymbol}`,
    market_id: '1',
    candles: raw.Data.map(d => {
      return {
        open: bn(d.open).multipliedBy(10 ** 18).toFixed(0),
        close: bn(d.close).multipliedBy(10 ** 18).toFixed(0),
        high: bn(d.high).multipliedBy(10 ** 18).toFixed(0),
        low: bn(d.low).multipliedBy(10 ** 18).toFixed(0),
        date: new Date(d.time * 1000).toISOString(),
      }
    }),
  };
}
