import React, {Component} from 'react';
import {connect} from 'react-redux';
import {ThunkDispatch} from 'redux-thunk';
import {Table, TableCell, TableHeader, TableHeaderRow, TableRow} from '../../components/ui/Table';
import {Module, ModuleContent, ModuleHeader, ModuleHeaderButton} from '../../components/Module';
import {REDUX_STATE} from '../../ducks';
import {OrderType} from '../../ducks/exchange';
import {ActionType} from '../../ducks/types';
import {fetchUserOrders, ORDER_HISTORY_FILTERS, setOrderHistoryFilter} from '../../ducks/user';
import {AssetType} from '../../ducks/assets';

type StateProps = {
  orders: OrderType[]
  orderHistoryFilter: ORDER_HISTORY_FILTERS
  baseAsset: AssetType | undefined
  quoteAsset: AssetType | undefined
}

type DispatchProps = {
  fetchUserOrders: () => void
  setOrderHistoryFilter: (f: ORDER_HISTORY_FILTERS) => ActionType<ORDER_HISTORY_FILTERS>
}

type Props = StateProps & DispatchProps

class History extends Component<Props> {
  componentWillMount() {
    this.props.fetchUserOrders();
  }

  render() {
    const {
      setOrderHistoryFilter,
      orderHistoryFilter,
      orders,
      baseAsset,
      quoteAsset,
    } = this.props;

    const filteredOrders = orders.filter(({status}) => {
      switch (orderHistoryFilter) {
        case ORDER_HISTORY_FILTERS.ALL:
          return true;
        case ORDER_HISTORY_FILTERS.OPEN:
          return status === 'OPEN';
        default:
          return false;
      }
    });

    if (!baseAsset || !quoteAsset) return <noscript />;

    return (
      <Module className="exchange__history">
        <ModuleHeader>
          <ModuleHeaderButton
            onClick={() => setOrderHistoryFilter(ORDER_HISTORY_FILTERS.OPEN)}
            active={orderHistoryFilter === ORDER_HISTORY_FILTERS.OPEN}
          >
            Open
          </ModuleHeaderButton>
          <ModuleHeaderButton
            onClick={() => setOrderHistoryFilter(ORDER_HISTORY_FILTERS.ALL)}
            active={orderHistoryFilter === ORDER_HISTORY_FILTERS.ALL}
          >
            All
          </ModuleHeaderButton>
        </ModuleHeader>
        <ModuleContent>
          <Table className="exchange__history__table">
            <TableHeaderRow>
              <TableHeader>Block</TableHeader>
              <TableHeader>Pair</TableHeader>
              <TableHeader>Type</TableHeader>
              <TableHeader>Side</TableHeader>
              <TableHeader>{`Price (${quoteAsset.symbol})`}</TableHeader>
              <TableHeader>{`Amount (${baseAsset.symbol})`}</TableHeader>
              <TableHeader>Filled</TableHeader>
              <TableHeader>Status</TableHeader>
            </TableHeaderRow>
            <div className="exchange__history__table-content">
              {
                filteredOrders.length
                  ?
                  filteredOrders.map(this.renderRow)
                  :
                  <TableRow>
                    <TableCell>No orders to display.</TableCell>
                  </TableRow>
              }
            </div>
          </Table>
        </ModuleContent>
      </Module>
    )
  }

  renderRow = (order: OrderType): React.ReactNode => {
    const { baseAsset, quoteAsset } = this.props;

    if (!baseAsset || !quoteAsset) return <noscript />;

    return (
      <TableRow key={order.id}>
        <TableCell>{ order.created_block }</TableCell>
        <TableCell>{`${baseAsset.symbol}/${quoteAsset.symbol}`}</TableCell>
        <TableCell>LIMIT</TableCell>
        <TableCell>{ order.direction }</TableCell>
        <TableCell>
          {
            order.price
              .div(10 ** quoteAsset.decimals)
              .toFixed(Math.min(6, quoteAsset.nativeDecimals))
          }
        </TableCell>
        <TableCell>
          {
            order.quantity
              .div(10 ** baseAsset.decimals)
              .toFixed(Math.min(6, baseAsset.nativeDecimals))
          }
        </TableCell>
        <TableCell>
          {
            order.quantity_filled
              .div(10 ** baseAsset.decimals)
              .toFixed(Math.min(6, baseAsset.nativeDecimals))
          }
        </TableCell>
        <TableCell>{ order.status }</TableCell>
      </TableRow>
    );
  }
}

function mapStateToProps(state: REDUX_STATE): StateProps {
  const {
    user: {
      orders: history,
      orderHistoryFilter,
    },
    exchange: { orders, selectedMarket, markets },
    assets: { assets, symbolToAssetId },
  } = state;
  const market = markets[selectedMarket] || {};
  const { baseSymbol, quoteSymbol } = market;
  const quoteAsset = assets[symbolToAssetId[quoteSymbol]] || {};
  const baseAsset = assets[symbolToAssetId[baseSymbol]] || {};
  return {
    orders: history.map(id => orders[id]),
    orderHistoryFilter,
    quoteAsset,
    baseAsset,
  }
}

function mapDispatchToProps(dispatch: ThunkDispatch<REDUX_STATE, any, ActionType<any>>): DispatchProps {
  return {
    fetchUserOrders: () => dispatch(fetchUserOrders()),
    setOrderHistoryFilter: (filter: ORDER_HISTORY_FILTERS) => dispatch(setOrderHistoryFilter(filter)),
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(History);
