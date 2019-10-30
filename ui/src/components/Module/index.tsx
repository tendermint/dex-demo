import React, {MouseEvent} from 'react';
import './module.scss';
import Button from "../ui/Button";

type PropTypes = {
  className?: String
  children?: React.ReactNode
}

type ClickablePropTypes = PropTypes & {
  onClick?: (event: MouseEvent) => void
  active?: boolean
}

export const Module: React.FC<PropTypes> = props => {
  const {
    className = '',
    children,
  } = props;

  return (
    <div className={`module ${className}`}>
      { children }
    </div>
  );
};

export const ModuleHeader: React.FC<PropTypes> = props => {
  const {
    className = '',
    children,
  } = props;

  return (
    <div className={`module__header ${className}`}>
      { children }
    </div>
  );
};

export const ModuleContent: React.FC<PropTypes> = props => {
  const {
    className = '',
    children,
  } = props;

  return (
    <div className={`module__content ${className}`}>
      { children }
    </div>
  );
};

export const ModuleHeaderButton: React.FC<ClickablePropTypes> = props => {
  const {
    className = '',
    children,
    onClick,
    active,
  } = props;

  return (
    <Button
      className={`module__header__button ${className}`}
      onClick={onClick}
      active={active}
    >
      { children }
    </Button>
  );
};