import {CandleStickJSON} from "../fetch";
import {bn} from "../bn";
import store from "../../store";
import {updateDailyStatsByMarketId} from "../../ducks/exchange";

type CandleCacheParams = {
  marketId: string
}

export default class CandleCache {
  marketId: string;

  map: {
    [iso: string]: CandleStickJSON
  } = {};

  list: string[] = [];

  updateDailyTimeout: any | null = null;

  constructor (params: CandleCacheParams) {
    this.marketId = params.marketId;
  }

  at (iso: string): CandleStickJSON | null {
    return this.map[iso] || null;
  }

  add (candle: CandleStickJSON) {
    // @ts-ignore
    this.list = [ ...new Set([...this.list, candle.date]) ];
    this.map[candle.date] = candle;
    this.updateDaily();
  }

  remove (iso: string) {
    delete this.map[iso];
    this.list = this.list.filter(dateString => dateString !== iso);
    this.updateDaily();
  }

  first (): CandleStickJSON {
    const first = this.list.reduce((earliest, iso) => {
      if (!earliest || iso < earliest) {
        return iso;
      }
      return earliest;
    }, '');

    return this.map[first];
  }

  last (): CandleStickJSON {
    const last = this.list.reduce((latest, iso) => {
      if (!latest || iso > latest) {
        return iso;
      }
      return latest;
    }, '');

    return this.map[last];
  }

  lastTwo (): CandleStickJSON[] {
    if (this.list.length < 2) return [];

    let last = this.list[0] > this.list[1]
      ? this.list[0]
      : this.list[1];
    let secondLast = this.list[0] < this.list[1]
      ? this.list[0]
      : this.list[1];

    this.list.forEach(iso => {
      if (iso > last) {
        secondLast = last;
        last = iso;
      } else if (iso > secondLast) {
        secondLast = iso;
      }
    });

    const bar1 = this.map[secondLast];
    const bar2 = this.map[last];

    if (!bar1 || !bar2) {
      return [];
    }

    return [bar1, bar2];
  }

  high (): CandleStickJSON {
    let highest: CandleStickJSON = this.map[this.list[0]];

    this.list.forEach(iso => {
      const candle = this.map[iso];

      if (!highest || bn(candle.high).isGreaterThan(highest.high)) {
        highest = candle;
      }
    });

    return highest;
  }

  low (): CandleStickJSON {
    let lowest: CandleStickJSON = this.map[this.list[0]];

    this.list.forEach(iso => {
      const candle = this.map[iso];

      if (!lowest || bn(candle.low).isLessThan(lowest.low)) {
        lowest = candle;
      }
    });

    return lowest;
  }

  updateDaily () {
    if (this.updateDailyTimeout) {
      clearTimeout(this.updateDailyTimeout);
    }

    this.updateDailyTimeout = setTimeout(() => {
      const current = new Date().getTime();

      this.list = this.list.filter((dateString) => {
        const time = new Date(dateString).getTime();
        const isWithin24Hour = current - time <= 24 * 60 * 60 * 1000;

        if (!isWithin24Hour) {
          delete this.map[dateString];
        }

        return isWithin24Hour;
      });

      const first = this.first();
      const h = this.high();
      const l = this.low();

      const open = first ? first.open : '0';
      const [secondLast, last]= this.lastTwo();
      const close = last ? last.close : '0';
      const secondClose = secondLast ? secondLast.close : '0';
      const high = h ? h.high : '0';
      const low = l ? l.low : '0';

      const extended = {
        lastPrice: bn(close),
        dayHigh: bn(high),
        dayLow: bn(low),
        dayChange: bn(close).minus(bn(open)),
        dayChangePercentage: bn(close).minus(bn(open)).dividedBy(bn(open)).toNumber(),
        prevPrice: bn(secondClose),
      };

      store.dispatch(updateDailyStatsByMarketId(extended, this.marketId));

      this.updateDailyTimeout = null;
    }, 250);
  }
}