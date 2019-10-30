import React, { Component } from 'react';
import './table.scss';

type PropTypes = {
  className?: String
  children?: React.ReactNode
  onClick?: () => void
  tabIndex?: number
}

export class Table extends Component<PropTypes>{
  render() {
    const {
      className = '',
      children,
    } = this.props;

    return (
      <div className={`table ${className}`}>
        {children}
      </div>
    );
  }
}

export class TableHeaderRow extends Component<PropTypes> {
  render() {
    const {
      children,
      className = '',
    } = this.props;

    return (
      <div className={`table__header-row ${className}`}>
        { children }
      </div>
    )
  }
}

export class TableHeader extends Component<PropTypes> {
  render() {
    const {
      children,
      className = '',
    } = this.props;

    return (
      <div className={`table__header ${className}`}>
        { children }
      </div>
    )
  }
}

export class TableRow extends Component<PropTypes> {
  render() {
    const {
      children,
      className = '',
      onClick,
      tabIndex,
    } = this.props;

    return (
      <div
        className={`table__row ${className}`}
        onClick={onClick}
        onKeyPress={() => {
          if (typeof tabIndex !== "undefined") {
            onClick && onClick()
          }
        }}
        tabIndex={tabIndex}
      >
        { children }
      </div>
    )
  }
}

export class TableCell extends Component<PropTypes> {
  render() {
    const {
      children,
      className = '',
    } = this.props;

    return (
      <div className={`table__cell ${className}`}>
        { children }
      </div>
    )
  }
}