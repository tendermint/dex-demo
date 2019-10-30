import { ActionType } from '../types';

enum THEME {
  'dark',
  'light'
}

enum LANGUAGE {
  'en',
  'en-us',
  'zh-cn',
  'zh-hk'
}

export type GeneralStateType = {
  theme: THEME
  language: LANGUAGE
}

const initialState = {
  theme: THEME.dark,
  language: LANGUAGE.en,
};

export default function generalReducer(state: GeneralStateType = initialState, action: ActionType<any>): GeneralStateType {
  switch (action.type) {
    default:
      return state;
  }
}
