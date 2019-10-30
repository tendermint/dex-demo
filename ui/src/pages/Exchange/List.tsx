import React, {Component, ReactNode} from 'react';
import { Dispatch} from "redux";
import { connect } from 'react-redux';
import {
  Module,
  ModuleHeader,
  ModuleHeaderButton,
  ModuleContent,
} from "../../components/Module";
import {
  Table,
  TableHeaderRow,
  TableHeader,
  TableRow,
  TableCell,
} from "../../components/ui/Table";
import {REDUX_STATE} from "../../ducks";
import {ActionType} from "../../ducks/types";
import {BatchType, selectBatch} from "../../ducks/exchange";
import {AssetType} from "../../ducks/assets";
import {getHHMMSS} from "../../utils/date-util";
import {findQuantityOverPrice, findQuantityUnderPrice} from "../../utils/exchange-utils";

enum ListTab {
  batch = 'batch',
  order = 'order',
}

type StateProps = {
  batches: {
    [blockId: string]: BatchType
  }
  baseAsset?: AssetType
  quoteAsset?: AssetType
}

type DispatchProps = {
  selectBatch: (blockId: string) => void
}

type PropTypes = StateProps & DispatchProps

type State = {
  currentTab: ListTab
}

class List extends Component<PropTypes, State> {
  state = {
    currentTab: ListTab.batch,
  };

  render() {
    return (
      <Module className="exchange__list">
        <ModuleHeader>
          <ModuleHeaderButton
            active={this.state.currentTab === ListTab.batch}
            onClick={() => this.setState({ currentTab: ListTab.batch })}
          >
            Batch
          </ModuleHeaderButton>
        </ModuleHeader>
        <ModuleContent>
          <Table>
            <TableHeaderRow>
              <TableHeader>Price</TableHeader>
              <TableHeader>Amount</TableHeader>
              <TableHeader>Time</TableHeader>
            </TableHeaderRow>
            { this.renderContent() }
          </Table>
        </ModuleContent>
      </Module>
    );
  }

  renderContent (): ReactNode {
    switch (this.state.currentTab) {
      case ListTab.batch:
        return (
          <div className="exchange__list__table-content">
            {
              Object.entries(this.props.batches)
                .reverse()
                .map(([_, batch]) => {
                  return this.renderRow(batch);
                })
            }
          </div>
        );
      default:
        return null;
    }
  }

  renderRow(batch: BatchType): React.ReactNode {
    const { baseAsset, quoteAsset, selectBatch} = this.props;

    if (!baseAsset || !quoteAsset) return null;

    const {
      blockId,
      clearingPrice,
      bids,
      asks,
      date,
    } = batch;

    const as = findQuantityUnderPrice(asks, clearingPrice);
    const ad = findQuantityOverPrice(bids, clearingPrice);

    const clearingQuantity = as.isLessThan(ad) ? as : ad;

    return (
      <TableRow
        key={blockId}
        onClick={() => selectBatch(blockId)}
        tabIndex={0}
      >
        <TableCell>
          {
            clearingPrice.div(10 ** quoteAsset.decimals)
              .toFixed(Math.min(quoteAsset.nativeDecimals, 4))
          }
        </TableCell>
        <TableCell>
          {
            clearingQuantity.div(10 ** baseAsset.decimals)
              .toFixed(Math.min(baseAsset.nativeDecimals, 6))
          }
        </TableCell>
        <TableCell>{getHHMMSS(date)}</TableCell>
      </TableRow>
    )
  }
}

function mapStateToProps(state: REDUX_STATE): StateProps {
  const {
    exchange: { selectedMarket, markets },
    assets: { assets, symbolToAssetId },
  } = state;
  const {
    batches = {},
    quoteSymbol = '',
    baseSymbol = '',
  } = markets[selectedMarket] || {};

  const baseAsset = assets[symbolToAssetId[baseSymbol]];
  const quoteAsset = assets[symbolToAssetId[quoteSymbol]];

  return {
    batches,
    baseAsset,
    quoteAsset,
  }
}

function mapDispatchToProps(dispatch: Dispatch<ActionType<any>>): DispatchProps {
  return {
    selectBatch: (blockId: string) => dispatch(selectBatch(blockId)),
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(List)