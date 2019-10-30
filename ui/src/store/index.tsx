import thunk from "redux-thunk";
import sync from "./middlewares/sync";
import {ActionType} from "../ducks/types";
import {SET_OPEN_ASKS_BY_MARKET_ID, SET_OPEN_BIDS_BY_MARKET_ID} from "../ducks/exchange";
import {applyMiddleware, createStore} from "redux";
import rootReducer from "../ducks";

const middleswares: any[] = [
  thunk,
  sync,
];

if (process.env.NODE_ENV === 'development') {
  const { createLogger } = require('redux-logger');
  middleswares.push(
    createLogger({
      predicate: (_: any, action: ActionType<any>) => {
        const blacklist: string[] = [
          SET_OPEN_BIDS_BY_MARKET_ID,
          SET_OPEN_ASKS_BY_MARKET_ID,
        ];

        return !blacklist.includes(action.type);
      },
      collapsed: true,
      diff: true,
    })
  );
}

const store = createStore(
  rootReducer,
  applyMiddleware(...middleswares)
);

export default store;
