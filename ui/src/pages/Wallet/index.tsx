import React, { Component } from 'react';
import SidePanel from "./SidePanel";
import WalletTable from "./WalletTable";
import {Route, Switch} from "react-router";
import {WALLET, TRANSFER} from "../../constants/routes";
import './style/wallet.scss';
import Transfer from "../Transfer";

class Wallet extends Component {
  render() {
    return (
      <div className="wallet">
        <SidePanel />
        <Switch>
          <Route path={WALLET} component={WalletTable} exact />
          <Route path={TRANSFER} component={Transfer} exact />
        </Switch>
      </div>
    )
  }
}

export default Wallet;
