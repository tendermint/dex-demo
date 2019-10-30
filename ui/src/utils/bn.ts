import { BigNumber as BN } from 'bignumber.js';



export const bn = (num: number | string | BN | undefined): BN => {
  if (!num) return new BN(0);
  return new BN(num);
};
