import React, { Component } from 'react';
import { BigNumber as BN } from 'bignumber.js';

type PropTypes = {
  className?: string
  value: string | number | BN
  decimals?: number
  displayDecimals?: number
  formatAsCurrency?: boolean
}

export default class Numeral extends Component<PropTypes> {
  render () {
    const {
      className = '',
      value,
      decimals = 0,
      displayDecimals,
      formatAsCurrency,
    } = this.props;

    const outputDecimals = typeof displayDecimals === 'number'
      ? displayDecimals
      : decimals;

    let val = value;
    val = new BN(val);
    val = val.dividedBy(10 ** decimals);

    if (val.isNaN()) {
      val = '-'
    } else if (formatAsCurrency) {
      val = val.toFormat(outputDecimals);
    } else {
      val = val.toFixed(outputDecimals);
    }

    return (
      <span
        className={`numeral ${className}`}
      >
        { val }
      </span>
    )
  }
}