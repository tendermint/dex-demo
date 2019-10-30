import React, { Component, MouseEvent } from 'react';
import c from 'classnames';
import './button.scss';

type PropTypes = {
  className?: string
  onClick?: (e: MouseEvent) => void
  disabled?: boolean
  active?: boolean
  children: React.ReactNode
  type?: 'buy' | 'sell' | 'primary' | 'secondary' | 'secondary-reverse' | 'link' | 'alert'
  loading?: boolean
}

const TYPE_TO_CLASS: {[k: string]: string} = {
  buy: 'button--buy-btn',
  sell: 'button--sell-btn',
  primary: 'button--primary',
  secondary: 'button--secondary',
  'secondary-reverse': 'button--secondary-reverse',
  link: 'button--link',
  alert: 'button--alert',
  default: '',
};

class Button extends Component<PropTypes> {
  onClick = (e: MouseEvent): void => {
    const { onClick } = this.props;

    // Clear focus state when click outside
    if (e.screenX > 0 || e.screenY > 0) {
      (e.target as HTMLButtonElement).blur();
    }

    if (onClick) {
      onClick(e);
    }
  };

  render() {
    const {
      className = '',
      children,
      disabled,
      active,
      type = 'default',
      loading,
    } = this.props;

    return (
      <button
        className={c('button', className, TYPE_TO_CLASS[type], {
          'button--active': active,
          'button--loading': loading,
        })}
        onClick={this.onClick}
        disabled={disabled}
      >
        {children}
      </button>
    )
  }
}

export default Button;
