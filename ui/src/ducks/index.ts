import { combineReducers } from 'redux';
import general, { GeneralStateType } from './general';
import exchange, { ExchangeStateType } from './exchange';
import user, { UserStateType } from './user';
import assets, { AssetStateType } from './assets';

export type REDUX_STATE = {
  general: GeneralStateType
  exchange: ExchangeStateType
  user: UserStateType
  assets: AssetStateType
}

export default combineReducers({
  general,
  exchange,
  user,
  assets,
});
