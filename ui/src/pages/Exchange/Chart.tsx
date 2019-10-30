import React, {Component, ReactNode} from 'react';
import {Dispatch} from 'redux';
import {connect} from 'react-redux';
import c from 'classnames';
import {Module, ModuleContent, ModuleHeader, ModuleHeaderButton,} from '../../components/Module';
import {ActionType} from "../../ducks/types";
import {CHART_TYPE, INTERVAL, selectBatch, setChartInterval, setChartType,} from "../../ducks/exchange";
import {REDUX_STATE} from "../../ducks";
import DepthChart from "../../components/DepthChart";
import TradingView from "../../components/TradingView";
import Icon from "../../components/ui/Icon";
import CancelBlue from "../../assets/icons/cancel-blue.svg";
import InfoIcon from "../../assets/icons/info-icon.svg";
import BatchExecutionChart from "../../components/BatchExecutionChart";
import InfoModal from "../../components/InfoModal";
const BatchInfoDoc = require("../../assets/docs/batch-info.md");

type StatePropTypes = {
  selectedInterval: INTERVAL
  selectedChartType: CHART_TYPE
  selectedBatch: string
  selectedMarket: string
}

type DispatchPropTypes = {
  setChartInterval: (i: INTERVAL) => void
  setChartType: (t: CHART_TYPE) => void
  selectBatch: (id: string) => void
}

type State = {
  isShowingBatchInfoModal: boolean
}

type PropTypes = StatePropTypes & DispatchPropTypes

class ChartView extends Component<PropTypes, State> {
  state = {
    isShowingBatchInfoModal: false,
  };

  openBatchInfoModal = () => this.setState({ isShowingBatchInfoModal: true });
  closeBatchInfoModal = () => this.setState({ isShowingBatchInfoModal: false });

  render() {
    return (
      <Module className="exchange__chart">
        { this.renderHeader() }
        { this.renderSubHeader() }
        <ModuleContent className="exchange__chart__content">
          { this.renderContent() }
        </ModuleContent>
        { this.renderBatchInfoModal() }
      </Module>
    )
  }

  renderBatchInfoModal (): ReactNode {
    const { isShowingBatchInfoModal } = this.state;

    return isShowingBatchInfoModal && (
      <InfoModal
        title="What is Frequent Batch Auction?"
        documentationUrl=""
        onClose={this.closeBatchInfoModal}
        markdownUrl={BatchInfoDoc}
      >
        hey
      </InfoModal>
    )
  }

  renderContent(): ReactNode | ReactNode[] {
    const { selectedChartType, selectedBatch } = this.props;

    switch (selectedChartType) {
      case CHART_TYPE.TradingView:
        return [
          <TradingView key="trading-view" />
        ];
      case CHART_TYPE.Depth:
        return [
          <TradingView
            key="trading-view"
            className="hidden"
          />,
          <DepthChart key="depth-chart" />,
        ];
      case CHART_TYPE.Batch:
        return [
          <TradingView
            key="trading-view"
            className="hidden"
          />,
          <BatchExecutionChart
            key="batch-chart"
            batchId={selectedBatch}
          />,
        ];
      default:
        return null;
    }
  }

  renderHeader(): ReactNode {
    const { selectedChartType, setChartType, selectedBatch, selectBatch, selectedMarket } = this.props;

    return (
      <ModuleHeader>
        <ModuleHeaderButton
          onClick={() => setChartType(CHART_TYPE.TradingView)}
          active={selectedChartType === CHART_TYPE.TradingView}
        >
          TradingView
        </ModuleHeaderButton>
        <ModuleHeaderButton
          onClick={() => setChartType(CHART_TYPE.Depth)}
          active={selectedChartType === CHART_TYPE.Depth}
        >
          Depth
        </ModuleHeaderButton>
        {
          selectedBatch && (
            <ModuleHeaderButton
              active={selectedChartType === CHART_TYPE.Batch}
              onClick={() => setChartType(CHART_TYPE.Batch)}
            >
              {`Batch # ${selectedMarket}-${selectedBatch}`}
              {
                selectedChartType === CHART_TYPE.Batch && (
                  <Icon
                    url={CancelBlue}
                    height={8}
                    width={8}
                    onClick={e => {
                      e && e.stopPropagation();
                      selectBatch('');
                    }}
                    tabIndex={0}
                  />
                )
              }
            </ModuleHeaderButton>
          )
        }
        <div className="exchange__chart__header-actions">
          <Icon
            className={c('exchange__chart__info-icon', {
              'exchange__chart__info-icon--show': selectedBatch,
            })}
            height={18}
            width={18}
            url={InfoIcon}
            tabIndex={selectedBatch ? 0 : undefined}
            onClick={this.openBatchInfoModal}
          />
        </div>
      </ModuleHeader>
    )
  }

  renderSubHeader(): ReactNode {
    const { setChartInterval, selectedInterval, selectedChartType } = this.props;

    if (selectedChartType !== CHART_TYPE.TradingView) {
      return null;
    }

    return (
      <ModuleHeader className="exchange__chart__subheader">
        <ModuleHeaderButton
          onClick={() => setChartInterval(INTERVAL['1m'])}
          active={selectedInterval === INTERVAL['1m']}
        >
          1m
        </ModuleHeaderButton>
        <ModuleHeaderButton
          onClick={() => setChartInterval(INTERVAL['5m'])}
          active={selectedInterval === INTERVAL['5m']}
        >
          5m
        </ModuleHeaderButton>
        <ModuleHeaderButton
          onClick={() => setChartInterval(INTERVAL['15m'])}
          active={selectedInterval === INTERVAL['15m']}
        >
          15m
        </ModuleHeaderButton>
        <ModuleHeaderButton
          onClick={() => setChartInterval(INTERVAL['1h'])}
          active={selectedInterval === INTERVAL['1h']}
        >
          1h
        </ModuleHeaderButton>
      </ModuleHeader>
    )
  }
}

function mapStateToProps(state: REDUX_STATE): StatePropTypes {
  const {
    selectedInterval,
    selectedChartType,
    selectedBatch,
    selectedMarket,
  } = state.exchange;

  return {
    selectedInterval,
    selectedChartType,
    selectedBatch,
    selectedMarket,
  };
}

function mapDispatchToProps(dispatch: Dispatch<ActionType<any>>): DispatchPropTypes {
  return {
    setChartInterval: (interval: INTERVAL) => dispatch(setChartInterval(interval)),
    setChartType: (type: CHART_TYPE) => dispatch(setChartType(type)),
    selectBatch: (blockId: string) => dispatch(selectBatch(blockId)),
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(ChartView);
