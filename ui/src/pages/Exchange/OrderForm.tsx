import React, {ChangeEvent, Component, ReactNode} from 'react';
import {connect} from "react-redux";
import {ThunkDispatch} from "redux-thunk";
import {Module, ModuleContent, ModuleHeader, ModuleHeaderButton} from '../../components/Module';
import Input from '../../components/ui/Input';
import Button from "../../components/ui/Button";
import { PlaceOrderRequest } from "../../utils/fetch";
import {ORDER_SIDE, placeOrder} from "../../ducks/exchange";
import {REDUX_STATE} from "../../ducks";
import {ActionType} from "../../ducks/types";
import {Spinner} from "../../components/ui/LoadingIndicator";
import {bn} from "../../utils/bn";
import {AssetType} from "../../ducks/assets";

type StateProps = {
  selectedMarket: string
  baseAsset: AssetType | undefined
  quoteAsset: AssetType | undefined
}

type DispatchProps = {
  placeOrder: (o: PlaceOrderRequest) => Promise<any>
}

type Props = StateProps & DispatchProps

type State = {
  bidPrice?: string
  bidAmount?: string
  bidTotal?: string
  askPrice?: string
  askAmount?: string
  askTotal?: string
  isPlacingBid: boolean
  isPlacingAsk: boolean
  askErrorMessage: string
  bidErrorMessage: string
}

class OrderForm extends Component<Props, State> {
  state = {
    bidPrice: '',
    bidAmount: '',
    bidTotal: '',
    askPrice: '',
    askAmount: '',
    askTotal: '',
    isPlacingBid: false,
    isPlacingAsk: false,
    askErrorMessage: '',
    bidErrorMessage: '',
  };

  isValid(type: ORDER_SIDE): boolean {
    const {
      bidPrice = 0,
      bidAmount = 0,
      bidTotal = 0,
      askPrice = 0,
      askAmount = 0,
      askTotal = 0,
      isPlacingBid,
      isPlacingAsk,
    } = this.state;

    if (isPlacingBid || isPlacingAsk) return false;

    switch (type) {
      case ORDER_SIDE.buy:
        return bidPrice > 0 && bidAmount > 0 && bidTotal > 0;
      case ORDER_SIDE.sell:
        return askPrice > 0 && askAmount > 0 && askTotal > 0;
      default:
        return false;
    }
  }

  onPriceChange = (e: ChangeEvent<HTMLInputElement>, type: ORDER_SIDE) => {
    const price = e.currentTarget.value;
    const {
      bidAmount = '0',
      askAmount = '0',
    } = this.state;
    const {
      quoteAsset,
    } = this.props;

    if (!quoteAsset) return;

    switch (type) {
      case ORDER_SIDE.buy:
        this.setState({
          bidPrice: price,
          bidTotal: bn(price).multipliedBy(bn(bidAmount)).toFixed(quoteAsset.nativeDecimals),
        });
        return;
      case ORDER_SIDE.sell:
        this.setState({
          askPrice: price,
          askTotal: bn(price).multipliedBy(bn(askAmount)).toFixed(quoteAsset.nativeDecimals),
        });
        return;
      default:
        return;
    }
  };

  onAmountChange = (e: ChangeEvent<HTMLInputElement>, type: ORDER_SIDE) => {
    const amount = e.currentTarget.value;
    const { bidPrice = '0', askPrice = '0' } = this.state;
    const { quoteAsset } = this.props;

    if (!quoteAsset) return;

    switch (type) {
      case ORDER_SIDE.buy:
        this.setState({
          bidAmount: amount,
          bidTotal: bn(bidPrice)
            .multipliedBy(bn(amount))
            .toFixed(quoteAsset.nativeDecimals),
        });
        return;
      case ORDER_SIDE.sell:
        this.setState({
          askAmount: amount,
          askTotal: bn(askPrice)
            .multipliedBy(bn(amount))
            .toFixed(quoteAsset.nativeDecimals),
        });
        return;
      default:
        return;
    }
  };

  onTotalChange = (e: ChangeEvent<HTMLInputElement>, type: ORDER_SIDE) => {
    const total = e.currentTarget.value;
    const { bidPrice = '0', askPrice = '0' } = this.state;
    const { baseAsset } = this.props;

    if (!baseAsset) return;

    switch (type) {
      case ORDER_SIDE.buy:
        this.setState({
          bidAmount: !bn(bidPrice).isZero()
            ? bn(total).div(bn(bidPrice)).toFixed(baseAsset.nativeDecimals)
            : '0',
          bidTotal: total,
        });
        return;
      case ORDER_SIDE.sell:
        this.setState({
          askAmount: !bn(askPrice).isZero()
            ? bn(total).div(bn(askPrice)).toFixed(baseAsset.nativeDecimals)
            : '0',
          askTotal: total,
        });
        return;
      default:
        return;
    }
  };

  bid = async () => {
    if (!this.isValid(ORDER_SIDE.buy)) return;

    const { baseAsset, quoteAsset, selectedMarket } = this.props;
    const { bidAmount, bidPrice } = this.state;

    if (!baseAsset || !quoteAsset) return;

    this.setState({ isPlacingBid: true });

    try {
      const resp = await this.props.placeOrder({
        market_id: selectedMarket,
        direction: 'BID',
        price: bn(bidPrice)
          .multipliedBy(10 ** quoteAsset.decimals)
          .toFixed(0),
        quantity: bn(bidAmount)
          .multipliedBy(10 ** baseAsset.decimals)
          .toFixed(0),
        type: 'LIMIT',
        time_in_force: 100,
      });
      const json = await resp.json();

      if (json.error) {
        if (json.error.message) {
          throw new Error(json.error.message);
        }

        if ((/insufficient account funds/gi).test(json.error)) {
          throw new Error('Insufficient Fund');
        }

        console.error(json.error);
        throw new Error('See console for error');
      }

      this.setState({
        isPlacingBid: false,
        bidAmount: '',
        bidPrice: '',
        bidTotal: '',
      });
    } catch (e) {
      this.setState({
        isPlacingBid: false,
        bidErrorMessage: e.message,
      });
    }

  };

  ask = async () => {
    if (!this.isValid(ORDER_SIDE.sell)) return;

    const { baseAsset, quoteAsset, selectedMarket } = this.props;
    const { askAmount, askPrice } = this.state;

    if (!baseAsset || !quoteAsset) return;

    this.setState({ isPlacingAsk: true });

    try {
      const resp = await this.props.placeOrder({
        market_id: selectedMarket,
        direction: 'ASK',
        price: bn(askPrice)
          .multipliedBy(10 ** quoteAsset.decimals)
          .toFixed(0),
        quantity: bn(askAmount)
          .multipliedBy(10 ** baseAsset.decimals)
          .toFixed(0),
        type: 'LIMIT',
        time_in_force: 100,
      });
      const json = await resp.json();

      if (json.error) {
        if (json.error.message) {
          throw new Error(json.error.message);
        }

        if ((/insufficient account funds/gi).test(json.error)) {
          throw new Error('Insufficient Fund');
        }

        console.error(json.error);
        throw new Error('See console for error');
      }

      this.setState({
        isPlacingAsk: false,
        askAmount: '',
        askPrice: '',
        askTotal: '',
      });
    } catch (e) {
      console.log(e);
      this.setState({
        isPlacingAsk: false,
        askErrorMessage: e.message,
      });
    }

  };

  render() {
    return (
      <Module className="exchange__order-form">
        <ModuleHeader>
          <ModuleHeaderButton active>Limit Order</ModuleHeaderButton>
          {/*<ModuleHeaderButton>Market Order</ModuleHeaderButton>*/}
        </ModuleHeader>
        <ModuleContent className="exchange__order-form__content">
          { this.renderBuy() }
          <div className="exchange__order-form__content__divider" />
          { this.renderSell() }
        </ModuleContent>
      </Module>
    );
  }

  renderBuy(): ReactNode {
    const { baseAsset, quoteAsset} = this.props;
    const { bidPrice, bidAmount, bidTotal, isPlacingBid, bidErrorMessage } = this.state;

    if (!baseAsset || !quoteAsset) return null;

    return (
      <div className="exchange__order-form__content__buy">
        <Input
          type="number"
          min="0"
          label="Price"
          placeholder="0.00"
          suffix={quoteAsset.symbol}
          onChange={e => this.onPriceChange(e, ORDER_SIDE.buy)}
          value={bidPrice}
          step="0.01"
        />
        <Input
          type="number"
          min="0"
          label="Amount"
          placeholder="0.00"
          suffix={baseAsset.symbol}
          onChange={e => this.onAmountChange(e, ORDER_SIDE.buy)}
          value={bidAmount}
          step="0.01"
        />
        <Input
          type="number"
          min="0"
          label="Total"
          placeholder="0.00"
          suffix={quoteAsset.symbol}
          onChange={e => this.onTotalChange(e, ORDER_SIDE.buy)}
          value={bidTotal}
          step="0.01"
        />
        <div className="input">
          <div className="input__label" />
          <div className="input__wrapper exchange__order-form__content__selectors">
            <div className="exchange__order-form__content__selectors__selector">25%</div>
            <div className="exchange__order-form__content__selectors__selector">50%</div>
            <div className="exchange__order-form__content__selectors__selector">75%</div>
            <div className="exchange__order-form__content__selectors__selector">Max</div>
          </div>
        </div>
        <div className="exchange__order-form__content__footer">
          { !!bidErrorMessage && (
            <div className="exchange__order-form__error-message">
              { bidErrorMessage }
            </div>
          )}
          <Button
            type="buy"
            disabled={!this.isValid(ORDER_SIDE.buy)}
            onClick={this.bid}
          >
            { isPlacingBid ? 'Placing Order...' : `Buy ${baseAsset.symbol}` }
          </Button>
        </div>
        { isPlacingBid && this.renderSpinner() }
      </div>
    )
  }

  renderSell() {
    const { baseAsset, quoteAsset} = this.props;
    const { askPrice, askAmount, askTotal, isPlacingAsk, askErrorMessage } = this.state;

    if (!baseAsset || !quoteAsset) return;

    return (
      <div className="exchange__order-form__content__sell">
        <Input
          type="number"
          min="0"
          label="Price"
          placeholder="0.00"
          suffix={quoteAsset.symbol}
          value={askPrice}
          onChange={e => this.onPriceChange(e, ORDER_SIDE.sell)}
          step="0.01"
        />
        <Input
          type="number"
          label="Amount"
          placeholder="0.00"
          suffix={baseAsset.symbol}
          min="0"
          value={askAmount}
          onChange={e => this.onAmountChange(e, ORDER_SIDE.sell)}
          step="0.01"
        />
        <Input
          type="number"
          label="Total"
          placeholder="0.00"
          suffix={quoteAsset.symbol}
          min="0"
          value={askTotal}
          onChange={e => this.onTotalChange(e, ORDER_SIDE.sell)}
          step="0.01"
        />
        <div className="input">
          <div className="input__label" />
          <div className="input__wrapper exchange__order-form__content__selectors">
            <div className="exchange__order-form__content__selectors__selector">25%</div>
            <div className="exchange__order-form__content__selectors__selector">50%</div>
            <div className="exchange__order-form__content__selectors__selector">75%</div>
            <div className="exchange__order-form__content__selectors__selector">Max</div>
          </div>
        </div>
        <div className="exchange__order-form__content__footer">
          { !!askErrorMessage && (
            <div className="exchange__order-form__error-message">
              { askErrorMessage }
            </div>
          )}
          <Button
            type="sell"
            disabled={!this.isValid(ORDER_SIDE.sell)}
            onClick={this.ask}
          >
            { isPlacingAsk ? 'Placing Order...' : `Sell ${baseAsset.symbol}` }
          </Button>
        </div>
        { isPlacingAsk && this.renderSpinner() }
      </div>
    )
  }

  renderSpinner(): ReactNode {
    return (
      <div className="exchange__order-form__spinner">
        <Spinner />
      </div>
    )
  }
}

function mapStateToProps (state: REDUX_STATE) {
  const {
    exchange: { selectedMarket, markets },
    assets: { assets, symbolToAssetId },
  } = state;

  const market = markets[selectedMarket] || {};
  const { baseSymbol, quoteSymbol } = market;
  const quoteAsset = assets[symbolToAssetId[quoteSymbol]] || {};
  const baseAsset = assets[symbolToAssetId[baseSymbol]] || {};

  return {
    selectedMarket,
    quoteAsset,
    baseAsset,
  }
}

function mapDispatchToProps(dispatch: ThunkDispatch<REDUX_STATE, PlaceOrderRequest, ActionType<any>>): DispatchProps {
  return {
    placeOrder: (order: PlaceOrderRequest) => dispatch(placeOrder(order)),
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(OrderForm);
