import React, {Component, ReactNode} from "react";
import './connect-wallet.scss';
import Button from "../../components/ui/Button";
import {
  CONNECT_WALLET__HARDWARE,
  CONNECT_WALLET__MOBILE,
  CONNECT_WALLET__SOFTWARE,
  CREATE_WALLET__SOFTWARE
} from "../../constants/routes";
import {RouteComponentProps, withRouter} from "react-router";

class ConnectWallet extends Component<RouteComponentProps> {

  render (): ReactNode {
    return (
      <div className="connect-wallet">
        <div className="connect-wallet__title">Connect to Your Wallet</div>
        <div className="connect-wallet__subtitle">
          <span>Don't have a wallet?</span>
          <Button
            type="link"
            onClick={() => this.props.history.push(CREATE_WALLET__SOFTWARE)}
          >
            Create a New Wallet
          </Button>
        </div>
        <div className="connect-wallet__options">
          {
            this.renderOption(
              'Hardware',
              'Ledger, Trezor, etc',
              CONNECT_WALLET__HARDWARE,
              '',
              'Recommended',
              true,
            )
          }
          {
            this.renderOption(
              'Mobile',
              'Trust Wallet, CoolWallet S, etc',
              CONNECT_WALLET__MOBILE,
              '',
              '',
              true,
            )
          }
          {
            this.renderOption(
              'Software',
              'Private Key and Seed Phrase',
              CONNECT_WALLET__SOFTWARE,
              'Not Recommended',
            )
          }
        </div>
      </div>
    );
  }

  renderOption (text: string, description: string, url: string, warning?: string, recommended?: string, disabled?: boolean): ReactNode {
    return (
      <button
        className="connect-wallet__option"
        disabled={disabled}
        onClick={() => this.props.history.push(url)}
      >
        <div className="connect-wallet__option__title">{text}</div>
        <div className="connect-wallet__option__subtitle">{description}</div>
        {
          warning && (
            <div className={"connect-wallet__option__warning"}>
              {warning}
            </div>
          )
        }
        {
          recommended && (
            <div className={"connect-wallet__option__recommended"}>
              {recommended}
            </div>
          )
        }
      </button>
    );
  }
}

export default withRouter(ConnectWallet);
