import React, { Component, ReactNode } from 'react';
import { connect } from 'react-redux';
import { ThunkDispatch } from "redux-thunk";
import {Switch, Route, Redirect} from 'react-router-dom';
import AppSidebar from '../../components/AppSiderbar';
import AppHeader from '../../components/AppHeader';
import Wallet from "../Wallet";
import Exchange from '../Exchange';
import Login from "../Login";
import {
  CONNECT_WALLET,
  CONNECT_WALLET__SOFTWARE,
  CREATE_WALLET__SOFTWARE,
  CONFIRM_SEEDPHRASE_BACKUP__SOFTWARE,
  EXCHANGE,
  WALLET
} from '../../constants/routes';
import {checkLogin, login} from '../../ducks/user';
import { ActionType } from "../../ducks/types";
import { REDUX_STATE } from "../../ducks";
import './main.scss';
import {Spinner} from "../../components/ui/LoadingIndicator";

type StateProps = {
  isLoggedIn?: boolean
}

type DispatchProps = {
  login: () => void
  checkLogin: () => void
}

type PropsType = StateProps & DispatchProps

class Main extends Component<PropsType> {
  componentWillMount() {
    this.props.checkLogin();
  }

  render() {
    const { isLoggedIn } = this.props;

    if (typeof isLoggedIn === 'undefined') {
      return (
        <div className="app app--loading">
          <Spinner />
        </div>
      );
    }

    return (
      <div className="app">
        <AppSidebar />
        <div className="app__body">
          <AppHeader />
          { isLoggedIn ? this.renderRoutes() : this.renderAuthRoutes() }
        </div>
      </div>
    );
  }

  renderAuthRoutes(): ReactNode {
    return (
      <div className="app__content">
        <Switch>
          <Route path={CONNECT_WALLET__SOFTWARE} component={Login} exact />
          <Redirect to={CONNECT_WALLET__SOFTWARE} />
        </Switch>
      </div>
    );
  }

  renderRoutes(): ReactNode {
    return (
      <div className="app__content">
        <Switch>
          <Route path={EXCHANGE} component={Exchange} exact />
          <Route path={WALLET} component={Wallet} />
          <Redirect to="/exchange" />
        </Switch>
      </div>
    );
  }
}

function mapStateToProps(state: REDUX_STATE): StateProps {
  return {
    isLoggedIn: state.user.isLoggedIn,
  };
}

function mapDispatchToProps(dispatch: ThunkDispatch<REDUX_STATE, any, ActionType<any>>): DispatchProps {
  return {
    login: () => dispatch(login('password')),
    checkLogin: () => dispatch(checkLogin()),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Main);
