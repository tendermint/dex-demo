import React, { Component, ReactNode } from 'react';
import { connect } from 'react-redux';
import { REDUX_STATE } from '../../ducks';
import './exchange-header.scss';
import Numeral from "../Numeral";
import {AssetType} from "../../ducks/assets";
import {DayStatsType} from "../../ducks/exchange";

type StatePropTypes = {
  baseAsset?: AssetType
  quoteAsset?: AssetType
  dayStats?: DayStatsType
};

type DispatchPropTypes = {}

type PropTypes = StatePropTypes & DispatchPropTypes;

class ExchangeHeader extends Component<PropTypes> {
  render() {
    const {
      quoteAsset,
      baseAsset,
      dayStats,
    } = this.props;


    if (!quoteAsset || !baseAsset || !dayStats) return <div />;

    const quoteSymbol = quoteAsset.symbol;
    const baseSymbol = baseAsset.symbol;

    const {
      lastPrice,
      prevPrice,
      dayChange,
      dayLow,
      dayHigh,
      dayChangePercentage,
    } = dayStats;

    const changePercentage = lastPrice.minus(prevPrice).div(prevPrice).toNumber();
    let lastPriceModifier = '';

    if (changePercentage < 0) {
      lastPriceModifier = 'negative';
    } else if (changePercentage > 0) {
      lastPriceModifier = 'positive';
    }

    return (
      <div className="exchange-header">
        { this.renderItem('Trading Pair', `${baseSymbol}/${quoteSymbol}`) }
        {
          this.renderItem(
            'Last Price',
            <Numeral
              value={lastPrice}
              decimals={quoteAsset.decimals}
              displayDecimals={quoteAsset.nativeDecimals}
              formatAsCurrency
            />,
            lastPriceModifier,
            lastPrice.isZero(),
          )
        }
        {
          this.renderItem(
            '24H Change',
            <div>
              <Numeral
                value={dayChange}
                decimals={quoteAsset.decimals}
                displayDecimals={quoteAsset.nativeDecimals}
                formatAsCurrency
              />
              <span>(</span>
              <span>{(dayChangePercentage * 100).toFixed(2)}</span>
              <span>%)</span>
            </div>,
            dayChange.isPositive() ? 'positive' : 'negative',
            dayChange.isZero(),
          )
        }
        {
          this.renderItem(
            '24H High',
            <Numeral
              value={dayHigh}
              decimals={quoteAsset.decimals}
              displayDecimals={quoteAsset.nativeDecimals}
              formatAsCurrency
            />,
            '',
            dayHigh.isZero(),
          )
        }
        {
          this.renderItem(
            '24H Low',
            <Numeral
              value={dayLow}
              decimals={quoteAsset.decimals}
              displayDecimals={quoteAsset.nativeDecimals}
              formatAsCurrency
            />,
            '',
            dayLow.isZero(),
          )
        }
      </div>
    )
  }

  renderItem(label: string, value: ReactNode, modifier: string = '', isLoading: boolean = false): ReactNode {
    if (isLoading) {
     return (
       <div className={`exchange-header__item exchange-header__item--loading`}>
         <div className="exchange-header__item__label">{label}</div>
         <div className="exchange-header__item__value" />
       </div>
     );
    }

    return (
      <div className={`exchange-header__item exchange-header__item--${modifier}`}>
        <div className="exchange-header__item__label">{label}</div>
        <div className="exchange-header__item__value">{value}</div>
      </div>
    )
  }
}

function mapStateToProps(state: REDUX_STATE) {
  const {
    exchange: {
      selectedMarket,
      markets,
    },
    assets: { assets, symbolToAssetId },
  } = state;
  const market = markets[selectedMarket] || { dayStats: {} };
  const baseAsset = assets[symbolToAssetId[market.baseSymbol]];
  const quoteAsset = assets[symbolToAssetId[market.quoteSymbol]];

  return {
    dayStats: market.dayStats,
    baseAsset,
    quoteAsset,
  }
}

export default connect(mapStateToProps)(ExchangeHeader)
