import React, {Component, ReactNode} from "react";
import {RouteComponentProps, withRouter} from "react-router";
import {connect} from "react-redux";
import Input from "../../components/ui/Input";
import Button from "../../components/ui/Button";
import Checkbox from "../../components/ui/Checkbox";
import "./create-wallet.scss";
import {CONFIRM_SEEDPHRASE_BACKUP__SOFTWARE, CONNECT_WALLET, CREATE_WALLET__SOFTWARE} from "../../constants/routes";

type StateProps = {

}

type DispatchProps = {

}

type Props = StateProps & DispatchProps & RouteComponentProps

type State = {
  password: string
  confirmPassword: string
  hasAcknowledged: boolean
}

class CreateWallet extends Component<Props, State> {
  state = {
    password: '',
    confirmPassword: '',
    hasAcknowledged: false,
  };

  isValid = () => {
    const {
      password,
      confirmPassword,
      hasAcknowledged,
    } = this.state;

    if (!password || !confirmPassword) {
      return false;
    }

    return password === confirmPassword && hasAcknowledged;
  };

  render(): ReactNode {
    const {
      password,
      confirmPassword,
      hasAcknowledged,
    } = this.state;

    return (
      <div className="connect-wallet">
        <div className="connect-wallet__title">Create a New Wallet</div>
        <div className="connect-wallet__subtitle">
          <span>Already have a wallet?</span>
          <Button
            type="link"
            onClick={() => this.props.history.push(CONNECT_WALLET)}
          >
            Connect to Your Wallet
          </Button>
        </div>
        <div className="create-wallet">
          <div className="create-wallet__title">Please set a password</div>
          <div className="create-wallet__info-text">
            We will encrypt your secrets with your password, and store them locally on your device. Please note that we do not have the power to recover your secrets without your password, and you MUST create your own secured back up.
          </div>
          <div className="create-wallet__form">
            <Input
              label="Password"
              type="password"
              value={password}
              onChange={e => this.setState({ password: e.target.value })}
              autoFocus
            />
            <Input
              label="Confirm Password"
              type="password"
              value={confirmPassword}
              onChange={e => this.setState({ confirmPassword: e.target.value })}
            />
          </div>
          <div className="create-wallet__acknowledgement">
            <Checkbox
              onChange={e => this.setState({ hasAcknowledged: e.target.checked })}
              checked={hasAcknowledged}
            >
              I understand that I can not reset or recover my password, and I will create my own back up securely.
            </Checkbox>

          </div>
          <div className="create-wallet__actions">
            <Button
              type="primary"
              disabled={!this.isValid()}
              onClick={() => this.props.history.push(CONFIRM_SEEDPHRASE_BACKUP__SOFTWARE)}
            >
              Set Password
            </Button>
          </div>
        </div>
      </div>
    );
  }
}

export default withRouter(
  connect()(CreateWallet)
);
