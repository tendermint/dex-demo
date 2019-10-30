import React, {Component} from 'react';
import * as am4core from "@amcharts/amcharts4/core";
import * as am4charts from "@amcharts/amcharts4/charts";
import {connect} from "react-redux";
import {REDUX_STATE} from "../../ducks";
import {AssetType} from "../../ducks/assets";
import {BatchType, Order} from "../../ducks/exchange";
import {bn} from "../../utils/bn";
import {reduceDepthFromOrders, sortOrders} from "../../utils/exchange-utils";

const COLOR_WHITE = am4core.color('#fff');
const COLOR_BLUE = am4core.color('#3084c3');
const COLOR_GREEN = am4core.color('#53b987');
const COLOR_RED = am4core.color('#eb4d5c');

type TooltipType = {
  fontSize: Number
  background: {
    fill: am4core.Color,
    fillOpacity: Number,
  }
}

type BatchChartData = {
  price: string
  value: string
}

type OwnProps = {
  batchId: string
}

type StateProps = {
  baseAsset?: AssetType
  quoteAsset?: AssetType
  bids: Order[]
  asks: Order[]
  batch?: BatchType
}

type Props = OwnProps & StateProps

class BatchExecutionChart extends Component<Props> {
  chart?: am4charts.XYChart;
  bidSeries?: am4charts.StepLineSeries;
  askSeries?: am4charts.StepLineSeries;
  priceAxis?: am4charts.ValueAxis;
  amountAxis?: am4charts.ValueAxis;
  clearingPriceSeries?: am4charts.LineSeries;

  componentDidMount() {
    this.hydrateChartWithData();
  }

  shouldComponentUpdate(nextProps: Props): boolean {
    const {
      batch,
    } = this.props;

    const {
      batch: nextBatch,
    } = nextProps;

    if (!batch || !nextBatch) {
      return false;
    }

    return batch.marketId !== nextBatch.marketId || batch.blockId !== nextBatch.blockId;
  }

  componentDidUpdate() {
    this.hydrateChartWithData();
  }

  componentWillUnmount() {
    if (this.chart) {
      this.chart.dispose();
    }
  }

  readyChart() {
    const { quoteAsset } = this.props;

    if (!this.chart && quoteAsset) {
      const chart = am4core.create('batchchart-wrapper', am4charts.XYChart);

      const amountAxis = chart.xAxes.push(new am4charts.ValueAxis());
      const priceAxis = chart.yAxes.push(new am4charts.ValueAxis());

      const priceTooltip = priceAxis.tooltip;
      const amountTooltip = amountAxis.tooltip;
      const bidSeries = chart.series.push(new am4charts.StepLineSeries());
      const askSeries = chart.series.push(new am4charts.StepLineSeries());
      const clearingPriceSeries = chart.series.push(new am4charts.LineSeries());
      const cursor = new am4charts.XYCursor();

      chart.cursor = cursor;
      chart.numberFormatter.numberFormat = `#.${Array(Math.min(quoteAsset.nativeDecimals, 4)).fill('0').join('')}`;
      // chart.scrollbarY = new am4charts.XYChartScrollbar();
      // (chart.scrollbarX as any).series.push(askSeries);

      // Configure Tooltip on X Axis
      (priceTooltip as TooltipType).background.fill = COLOR_WHITE;
      (priceTooltip as TooltipType).background.fillOpacity = .2;
      (priceTooltip as TooltipType).fontSize = 12;

      // Configure Tooltip on Y Axis
      (amountTooltip as TooltipType).background.fill = COLOR_WHITE;
      (amountTooltip as TooltipType).background.fillOpacity = .2;
      (amountTooltip as TooltipType).fontSize = 12;

      // Configure Cursor
      cursor.lineX.stroke = COLOR_WHITE;
      cursor.lineX.strokeWidth = 1;
      cursor.lineX.strokeOpacity = 0.5;
      cursor.lineY.stroke = COLOR_WHITE;
      cursor.lineY.strokeWidth = 1;
      cursor.lineY.strokeOpacity = 0.5;

      // Configure X Axis
      amountAxis.renderer.minGridDistance = 50;
      amountAxis.renderer.grid.template.strokeOpacity = .05;
      amountAxis.renderer.grid.template.stroke = COLOR_WHITE;
      amountAxis.renderer.grid.template.strokeWidth = 1;
      amountAxis.renderer.labels.template.fill = COLOR_WHITE;
      amountAxis.renderer.labels.template.opacity = .5;
      amountAxis.renderer.labels.template.fontSize = 10;
      amountAxis.renderer.labels.template.fontFamily = 'Roboto';
      amountAxis.renderer.labels.template.fontWeight = '300';
      amountAxis.min = 0;
      amountAxis.strictMinMax = true;

      // Configure Y Axis
      priceAxis.renderer.minGridDistance = 50;
      priceAxis.renderer.grid.template.strokeOpacity = .05;
      priceAxis.renderer.grid.template.stroke = COLOR_WHITE;
      priceAxis.renderer.grid.template.strokeWidth = 1;
      priceAxis.renderer.labels.template.fill = COLOR_WHITE;
      priceAxis.renderer.labels.template.opacity = .5;
      priceAxis.renderer.labels.template.fontSize = 10;
      priceAxis.renderer.labels.template.fontFamily = 'Roboto';
      priceAxis.renderer.labels.template.fontWeight = '300';
      priceAxis.strictMinMax = true;

      bidSeries.strokeWidth = 2;
      bidSeries.stroke = COLOR_GREEN;
      bidSeries.fill = COLOR_GREEN;
      bidSeries.fillOpacity = 0.1;

      askSeries.strokeWidth = 2;
      askSeries.stroke = COLOR_RED;
      askSeries.fill = COLOR_RED;
      askSeries.fillOpacity = 0.1;
      askSeries.baseAxis = priceAxis;

      clearingPriceSeries.strokeWidth = 1;
      clearingPriceSeries.stroke = COLOR_BLUE;
      clearingPriceSeries.strokeDasharray = '4 2';

      this.chart = chart;
      this.bidSeries = bidSeries;
      this.askSeries = askSeries;
      this.clearingPriceSeries = clearingPriceSeries;
      this.priceAxis = priceAxis;
      this.amountAxis = amountAxis;
    }
  }

  hydrateChartWithData() {
    const {quoteAsset, baseAsset, batch} = this.props;

    if (!quoteAsset || !baseAsset || !batch) return;

    this.readyChart();
    if (!this.chart || !this.bidSeries || !this.askSeries || !this.priceAxis || !this.amountAxis || !this.clearingPriceSeries) {
      return;
    }

    const bids: BatchChartData[] = [];
    const asks: BatchChartData[] = [];

    let accumBids = bn(0);
    let accumAsks = bn(0);

    const sortedBids = sortOrders(this.props.bids, true);

    sortedBids.forEach((b, i) => {
      const valueN =  b.price
        .multipliedBy(b.quantity.div(10 ** baseAsset.decimals));

      if (i === 0) {
        bids.push({
          price: b.price
            .toFixed(quoteAsset.nativeDecimals),
          value: bn(0).toFixed(quoteAsset.nativeDecimals),
        });
      }

      bids.push({
        price: b.price
          .toFixed(quoteAsset.nativeDecimals),
        value: valueN
          .plus(accumBids)
          .toFixed(quoteAsset.nativeDecimals),
      });
      accumBids = valueN.plus(accumBids);
    });

    this.bidSeries.data = bids;
    this.bidSeries.dataFields.valueX = 'value';
    this.bidSeries.dataFields.valueY= 'price';


    const sortedAsks = sortOrders(this.props.asks);

    sortedAsks.forEach((b, i) => {
      const valueN =  b.price
        .multipliedBy(b.quantity.div(10 ** baseAsset.decimals));

      if (i === 0) {
        asks.push({
          price: b.price
            .toFixed(quoteAsset.nativeDecimals),
          value: bn(0).toFixed(quoteAsset.nativeDecimals),
        });
      }

      asks.push({
        price: b.price
          .toFixed(quoteAsset.nativeDecimals),
        value: valueN.plus(accumAsks).toFixed(quoteAsset.nativeDecimals),
      });
      accumAsks = valueN.plus(accumAsks);
    });

    this.askSeries.data = asks;
    this.askSeries.dataFields.valueX = 'value';
    this.askSeries.dataFields.valueY= 'price';

    const cp = batch.clearingPrice
      .div(10 ** quoteAsset.decimals)
      .toFixed(quoteAsset.nativeDecimals);

    this.clearingPriceSeries.data = [
      {
        price: cp,
        value: '0',
      },
      {
        price: cp,
        value: accumBids.isGreaterThan(accumAsks)
          ? accumBids.toFixed(quoteAsset.nativeDecimals)
          : accumAsks.toFixed(quoteAsset.nativeDecimals),
      },
    ];
    this.clearingPriceSeries.dataFields.valueX = 'value';
    this.clearingPriceSeries.dataFields.valueY= 'price';
    // this.clearingPriceSeries.tooltip

    this.clearingPriceSeries.showTooltipAtDataItem(this.clearingPriceSeries.dataItem);

    this.chart.invalidateRawData();
  }

  render() {
    return (
      <div style={{ height: '100%', width: '100%' }} id="batchchart-wrapper" />
    );
  }
}

function mapStateToProps (state: REDUX_STATE, ownProps: OwnProps): StateProps {
  const {
    exchange: {
      selectedMarket,
      markets,
    },
    assets: { assets, symbolToAssetId },
  } = state;

  const market = markets[selectedMarket];
  const { baseSymbol = '', quoteSymbol = '', batches = {} } = market || {};
  const baseAsset = assets[symbolToAssetId[baseSymbol]];
  const quoteAsset = assets[symbolToAssetId[quoteSymbol]];
  const batch = batches[ownProps.batchId];

  const { depths: bidDepth } = reduceDepthFromOrders(
    batch ? batch.bids : [],
    quoteAsset.decimals,
    quoteAsset.nativeDecimals,
  );
  const { depths: askDepth } = reduceDepthFromOrders(
    batch ? batch.asks : [],
    quoteAsset.decimals,
    quoteAsset.nativeDecimals,
  );

  // const { clearingPrice } = estimateBatch(bidDepth, askDepth, quoteAsset.decimals, quoteAsset.nativeDecimals);
  return {
    baseAsset,
    quoteAsset,
    bids: bidDepth,
    asks: askDepth,
    batch,
  };
}

export default connect(mapStateToProps)(BatchExecutionChart);
