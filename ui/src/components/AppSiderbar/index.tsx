import React, { Component } from 'react';
import { withRouter, RouteComponentProps } from 'react-router-dom';
import Item from './app-sidebar-item';
// import Logo from '../../assets/icons/ironman.svg';
import Candlestick from '../../assets/icons/candlestick.svg';
import CandlestickBlue from '../../assets/icons/candlestick-blue.svg';
import Wallet from '../../assets/icons/wallet.svg';
import WalletBlue from '../../assets/icons/wallet-blue.svg';
import Tendermint from '../../assets/icons/tendermint.svg';
import {
  HOME,
  WALLET,
  EXCHANGE,
  DEPOSIT,
  WITHDRAWAL,
  CONNECT_WALLET,
  CONNECT_WALLET__HARDWARE,
  CONNECT_WALLET__MOBILE,
  CONNECT_WALLET__SOFTWARE,
  CREATE_WALLET__SOFTWARE,
  CONFIRM_SEEDPHRASE_BACKUP__SOFTWARE,
  TRANSFER,
} from '../../constants/routes';
import "./app-sidebar.scss";

class AppSidebar extends Component<RouteComponentProps> {
  render() {
    const {
      location: {
        pathname,
      },
    } = this.props;

    const isWalletSelected = [
      WALLET,
      DEPOSIT,
      WITHDRAWAL,
      CONNECT_WALLET,
      CONNECT_WALLET__HARDWARE,
      CONNECT_WALLET__MOBILE,
      CONNECT_WALLET__SOFTWARE,
      CREATE_WALLET__SOFTWARE,
      CONFIRM_SEEDPHRASE_BACKUP__SOFTWARE,
      TRANSFER,
    ].includes(pathname);
    const isExchangeSelected = [EXCHANGE].includes(pathname);

    return (
      <div className="app-sidebar">
        <Item
          imageUrl={ Tendermint }
        />

        <Item
          imageUrl={isExchangeSelected ? CandlestickBlue : Candlestick}
          onClick={() => this.props.history.push(EXCHANGE)}
          selected={isExchangeSelected}
          hoverable
        />
        <Item
          imageUrl={isWalletSelected ? WalletBlue : Wallet}
          onClick={() => this.props.history.push(WALLET)}
          selected={isWalletSelected}
          hoverable
        />
      </div>
    )
  }
}

export default withRouter(AppSidebar);
