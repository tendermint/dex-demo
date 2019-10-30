import React from 'react';
import Chart from './Chart';
import OrderForm from './OrderForm';
import History from './History';
import Depth from './Depth';
import List from './List';
import './style/exchange.scss';

const Exchange: React.FC = () => {
  return (
    <div className="exchange">
      <Chart />
      <Depth />
      <List />
      <History />
      <OrderForm />
    </div>
  )
}

export default Exchange;
