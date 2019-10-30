import React, { Component, ReactNode } from 'react';
import './notification-indicator.scss';

type Props = {
  children: ReactNode
  count: number
}

export default class NotificationIndicator extends Component<Props> {
  render() {
    const { children, count } = this.props;

    return (
      <div className="notification-indicator">
        {children}
        {
          count > 0
            ? (
              <div className="notification-indicator__counter">
                {count}
              </div>
            )
            : null

        }
      </div>
    )
  }
}