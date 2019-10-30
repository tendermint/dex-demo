import React, {Component, ReactNode} from "react";
import {Dispatch} from "redux";
import {connect} from "react-redux";
import c from 'classnames';
import {Module, ModuleContent, ModuleHeader, ModuleHeaderButton,} from "../../../components/Module";
import {Table, TableCell, TableHeader, TableHeaderRow, TableRow,} from "../../../components/ui/Table";
import {REDUX_STATE} from "../../../ducks";
import {Order, setSpreadType, SPREAD_TYPE} from "../../../ducks/exchange";
import {ActionType} from "../../../ducks/types";
import {reduceDepthFromOrders, sortOrders} from '../../../utils/exchange-utils';
import BigNumber, {BigNumber as BN} from 'bignumber.js';
import {AssetType} from "../../../ducks/assets";
import DepthTableRow from "./DepthTableRow";
import {bn} from "../../../utils/bn";
import Tooltip from "../../../components/ui/Tooltip";

type StatePropTypes = {
  selectedSpreadType: SPREAD_TYPE
  max: BN
  bids: Order[]
  asks: Order[]
  baseAsset: AssetType | undefined
  quoteAsset: AssetType | undefined
  lastPrice: BigNumber | undefined
  prevPrice: BigNumber | undefined
}

type DispatchPropTypes = {
  setSpreadType: (t: SPREAD_TYPE) => void,
}

type PropTypes = StatePropTypes & DispatchPropTypes

class Index extends Component<PropTypes> {
  render() {
    return (
      <Module className="exchange__depth">
        { this.renderHeader() }
        <ModuleContent>
          <Table>
            <TableHeaderRow>
              <TableHeader>Price</TableHeader>
              <TableHeader>Amount</TableHeader>
              <TableHeader>Total</TableHeader>
            </TableHeaderRow>
            <div className="exchange__depth__content">
              { this.renderAsks() }
              { this.renderSpreadRow() }
              { this.renderBuys() }
            </div>
          </Table>
        </ModuleContent>
      </Module>
    );
  }

  renderSpreadRow (): ReactNode {
    const {
      lastPrice,
      prevPrice,
      quoteAsset,
    } = this.props;

    if (!quoteAsset || !lastPrice || lastPrice.isZero() || !prevPrice || prevPrice.isZero()) {
      return (
        <div className="exchange__depth__spread exchange__depth__spread--loading">
          <TableRow>
            <TableCell>&nbsp;</TableCell>
            <TableCell>&nbsp;</TableCell>
          </TableRow>
        </div>
      )
    }

    const changePercentage = lastPrice.minus(prevPrice).div(prevPrice).toNumber();

    return (
      <div
        className={c("exchange__depth__spread", {
          'exchange__depth__spread--up': changePercentage > 0,
          'exchange__depth__spread--down': changePercentage < 0,
        })}
      >
        <TableRow>
          <Tooltip content="Last Clearing Price">
            <TableCell>
              {
                lastPrice
                  .div(10 ** quoteAsset.decimals)
                  .toFixed(Math.min(quoteAsset.nativeDecimals, 4))
              }
            </TableCell>
          </Tooltip>
          <Tooltip content="Since Last Tick">
            <TableCell>{(changePercentage * 100).toFixed(2)}%</TableCell>
          </Tooltip>
        </TableRow>
      </div>
    );
  }

  renderBuys (): ReactNode {
    const { bids, selectedSpreadType } = this.props;

    const sorted = sortOrders(bids, true);

    return (
      <div
        className={c('exchange__depth__group exchange__depth__buy', {
          'exchange__depth__group--hidden': selectedSpreadType === SPREAD_TYPE.ASK,
        })}
      >
        { sorted.map(({ price, quantity }) => this.renderRow(price, quantity, SPREAD_TYPE.BID)) }
      </div>
    )
  }

  renderAsks(): ReactNode {
    const { asks, selectedSpreadType } = this.props;

    const sorted = sortOrders(asks);

    return (
      <div
        className={c('exchange__depth__group exchange__depth__sell', {
          'exchange__depth__group--hidden': selectedSpreadType === SPREAD_TYPE.BID,
        })}
      >
        { sorted.map(({ price, quantity }) => this.renderRow(price, quantity, SPREAD_TYPE.ASK)) }
      </div>
    )
  }

  renderRow(price: BN, quantity: BN, side: SPREAD_TYPE): ReactNode {
    const {
      max,
      quoteAsset,
      baseAsset,
    } = this.props;

    if (!baseAsset || !quoteAsset) return null;

    return (
      <DepthTableRow
        key={`${price.toString()}}`}
        price={price}
        quantity={quantity}
        max={max}
        side={side}
        baseAsset={baseAsset}
        quoteAsset={quoteAsset}
      />
    );
  }

  renderHeader(): ReactNode {
    const { selectedSpreadType, setSpreadType } = this.props;

    return (
      <ModuleHeader>
        <ModuleHeaderButton
          onClick={() => setSpreadType(SPREAD_TYPE.ALL)}
          active={selectedSpreadType === SPREAD_TYPE.ALL}
        >
          All
        </ModuleHeaderButton>
        <ModuleHeaderButton
          onClick={() => setSpreadType(SPREAD_TYPE.BID)}
          active={selectedSpreadType === SPREAD_TYPE.BID}
        >
          BID
        </ModuleHeaderButton>
        <ModuleHeaderButton
          onClick={() => setSpreadType(SPREAD_TYPE.ASK)}
          active={selectedSpreadType === SPREAD_TYPE.ASK}
        >
          ASK
        </ModuleHeaderButton>
      </ModuleHeader>
    );
  }
}

function mapStateToProps(state: REDUX_STATE): StatePropTypes {
  const {
    exchange: { selectedMarket, markets, selectedSpreadType },
    assets: { assets, symbolToAssetId },
  } = state;

  const market = markets[selectedMarket];
  const { bids = [], asks = [], baseSymbol = '', quoteSymbol = '' } = market || {};
  const quoteAsset = assets[symbolToAssetId[quoteSymbol]] || {};
  const baseAsset = assets[symbolToAssetId[baseSymbol]] || {};
  const { max: bidMax, depths: bidDepth } = reduceDepthFromOrders(bids, quoteAsset.decimals, quoteAsset.nativeDecimals);
  const { max: askMax, depths: askDepth } = reduceDepthFromOrders(asks, quoteAsset.decimals, quoteAsset.nativeDecimals);

  return {
    selectedSpreadType: selectedSpreadType,
    bids: bidDepth,
    asks: askDepth,
    max: BN.max(bidMax, askMax).multipliedBy(1.25),
    quoteAsset,
    baseAsset,
    lastPrice: market.dayStats ? market.dayStats.lastPrice : bn(0),
    prevPrice: market.dayStats ? market.dayStats.prevPrice : bn(0),
  }
}

function mapDispatchToProps(dispatch: Dispatch<ActionType<any>>): DispatchPropTypes {
  return {
    setSpreadType: (type: SPREAD_TYPE) => dispatch(setSpreadType(type)),
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Index);
