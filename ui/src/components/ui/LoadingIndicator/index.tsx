import React from 'react';
import './loading.scss';

type Props = {
  className?: string
}

export const Spinner: React.FC<Props> = ({ className = '' }) => {
  return (
    <div className={`loading__spinner ${className}`}>
      <div />
      <div />
      <div />
      <div />
      <div />
      <div />
      <div />
      <div />
      <div />
      <div />
      <div />
      <div />
    </div>
  );
};
