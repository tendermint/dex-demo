import React, { Component } from 'react';
import SearchBox from '../SeachBox';
import ExchangeHeader from '../ExchangeHeader';
import NetworkDropdown from '../NetworkDropdown';
import {Route, Switch} from 'react-router';
import "./app-header.scss";
import {EXCHANGE} from "../../constants/routes";

class AppHeader extends Component {
  render() {
    return (
      <div className="app-header">
        <div className="app-header__content">
          <Switch>
            <Route path={EXCHANGE} component={ExchangeHeader} />
          </Switch>
        </div>
      </div>
    )
  }
}

export default AppHeader;
