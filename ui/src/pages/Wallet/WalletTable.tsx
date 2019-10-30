import React, {Component, ReactNode} from 'react';
import {connect} from 'react-redux';
import c from 'classnames';
import {Module, ModuleContent, ModuleHeader, ModuleHeaderButton,} from "../../components/Module";
import {REDUX_STATE} from "../../ducks";
import {AssetType} from "../../ducks/assets";
import {RouteComponentProps, withRouter} from "react-router";
import {Dispatch} from "redux";
import {BalanceType} from "../../ducks/user";
import BalanceTable from "./BalanceTable";
import CopyIcon from "../../components/ui/CopyIcon";
import Button from "../../components/ui/Button";
import {TRANSFER} from "../../constants/routes";
import QRIcon from "../../components/QRIcon";

enum Tabs {
  Balances = 'Balances',
  Deposits = 'Deposits',
  Withdrawals = 'Withdrawals',
}

type StateProps = {
  balances: {
    [assetId: string]: BalanceType
  }
  assets: {
    [k: string]: AssetType
  }
  address: string
}

type DispatchProps = {
}

type Props = StateProps & DispatchProps & RouteComponentProps

type State = {
  currentTab: Tabs
}

class WalletTable extends Component<Props, State> {
  state = {
    currentTab: Tabs.Balances,
  };

  render() {
    const { address } = this.props;
    const { currentTab } = this.state;

    return (
      <div className="wallet__content">
        <div className="wallet__content__header">
            My Wallet
        </div>
        <div className="wallet__content__subheader">
          <div className="wallet__content__subheader__address-group">
            <div className="wallet__content__subheader__address-group__label">Address</div>
            <div
              className={c('wallet__content__subheader__address-group__value', {
                'wallet__content__subheader__address-group__value--loading': !address,
              })}
            >
              {address}
              { !!address && <CopyIcon copyText={address} /> }
              { !!address && <QRIcon text={address} /> }
            </div>
          </div>
          <div className="wallet__content__subheader__actions">
            <Button
              type="primary"
              onClick={() => this.props.history.push(TRANSFER)}
            >
              Transfer
            </Button>
          </div>
        </div>
        <Module className="wallet__content__table">
          <ModuleHeader>
            <ModuleHeaderButton
              onClick={() => this.setState({ currentTab: Tabs.Balances})}
              active={currentTab === Tabs.Balances}
            >
              Balances
            </ModuleHeaderButton>
          </ModuleHeader>
          <ModuleContent>
            { this.renderTable() }
          </ModuleContent>
        </Module>
      </div>
    )
  }

  renderTable (): ReactNode {
    switch (this.state.currentTab) {
      case Tabs.Balances:
        return <BalanceTable />;
      default:
        return null;
    }
  }
}

function mapStateToProps (state: REDUX_STATE): StateProps {
  const {
    user: {
      address,
      balances,
    }
  } = state;
  return {
    balances,
    address,
    assets: state.assets.assets,
  }
}

function mapDispatchToProps (dispatch: Dispatch): DispatchProps {
  return {}
}

export default withRouter(
  connect(mapStateToProps, mapDispatchToProps)(WalletTable)
);