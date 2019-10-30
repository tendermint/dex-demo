import React, {Component, ReactNode} from "react";
import {RouteComponentProps, withRouter} from "react-router";
import {connect} from "react-redux";
import Button from "../../components/ui/Button";
import "./login.scss";
import Input from "../../components/ui/Input";
import {ThunkDispatch} from "redux-thunk";
import {REDUX_STATE} from "../../ducks";
import {ActionType} from "../../ducks/types";
import {login} from "../../ducks/user";
import {CREATE_WALLET__SOFTWARE, WALLET} from "../../constants/routes";

type DispatchProps = {
  login: (pw: string) => Promise<Response>
}

type Props = DispatchProps & RouteComponentProps;

type State = {
  password: string
  errorMessage: string,
  isLoggingIn: boolean,
}

class Login extends Component<Props, State> {
  state = {
    password: '',
    errorMessage: '',
    isLoggingIn: false,
  };

  login = async () => {
    const { password } = this.state;

    this.setState({
      isLoggingIn: true,
      errorMessage: '',
    });

    const resp = await this.props.login(password);

    if (resp.status !== 204) {
      this.setState({
        isLoggingIn: false,
        errorMessage: await resp.text(),
      });
    } else {
      this.setState({
        isLoggingIn: false,
      });
      this.props.history.push(WALLET);
    }
  };

  render (): ReactNode {
    const { isLoggingIn } = this.state;

    return (
      <div className="connect-wallet">
        <div className="login">
          <div className="login__title">
            Please Login
          </div>
          <Input
            label="Password"
            type="password"
            onChange={e => this.setState({
              password: e.target.value,
              errorMessage: '',
            })}
            onKeyPress={e => {
              if (e.key === 'Enter') {
                e.stopPropagation();
                this.login();
              }
            }}
            value={this.state.password}
            autoFocus
          />
          <div className="login__error-message">
            { this.state.errorMessage }
          </div>
          <div className="login__actions">
            <Button
              type="primary"
              onClick={this.login}
              disabled={!this.state.password || isLoggingIn}
              loading={isLoggingIn}
            >
              { isLoggingIn ? 'Logging In' : 'Login' }
            </Button>
          </div>
        </div>
      </div>
    )
  }
}

function mapDispatchToProps (dispatch: ThunkDispatch<REDUX_STATE, any, ActionType<any>>): DispatchProps {
  return {
    login: (pw: string) => dispatch(login(pw)),
  }
}

export default withRouter(
  connect(null, mapDispatchToProps)(Login)
);
