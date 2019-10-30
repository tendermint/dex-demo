import React, {Component, ReactNode} from "react";
import {
  ChartingLibraryWidgetOptions,
  IChartingLibraryWidget
} from "../../../public/charting_library/charting_library.min";
import {baseWidgetOption} from "../../utils/trading-view-util";
import {connect} from "react-redux";
import {REDUX_STATE} from "../../ducks";
import {INTERVAL} from "../../ducks/exchange";

import './trading-view.scss';

type StateProps = {
  marketSymbol: string
  selectedInterval: INTERVAL
}

type OwnProps = {
  className?: string
}

type Props = StateProps & OwnProps

class TradingView extends Component<Props> {
  widget: IChartingLibraryWidget | null = null;

  componentDidMount() {
    this.getWidget();
  }

  componentWillUpdate(nextProps: Props) {
    const widget = this.getWidget();
    if (nextProps.selectedInterval !== this.props.selectedInterval) {
      if (widget) {
        widget
          .chart()
          .setResolution(nextProps.selectedInterval, () => {
            // console.log('Resolution changed');
          });
      }
    }
  }

  getWidget(): IChartingLibraryWidget | null {
    if (this.widget) {
      return this.widget;
    }

    const {
      TradingView: tv,
    } = window as any;

    if (!tv) {
      return null;
    }

    const widgetOptions: ChartingLibraryWidgetOptions = {
      ...baseWidgetOption,
      symbol: 'DEMO/TEST',
      interval: '1',
      // debug: true,
    };

    const widget = (window as any).tvWidget = new tv.widget(widgetOptions);
    this.widget = widget;

    return widget;
  }

  render (): ReactNode {
    const { className = '' } = this.props;

    return (
      <div
        id="tv-container"
        className={`trading-view ${className}`}
      />
    )
  }
}

function mapStateToProps (state: REDUX_STATE): StateProps {
  const {
    exchange: {
      selectedMarket,
      markets,
      selectedInterval,
    },
  } = state;

  const market = markets[selectedMarket] || {};

  return {
    marketSymbol: `${market.baseSymbol}/${market.quoteSymbol}`,
    selectedInterval,
  };
}

export default connect(mapStateToProps)(TradingView);
