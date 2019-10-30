import React, {Component, ReactNode} from "react";
import "./transfer.scss";
import Icon from "../../components/ui/Icon";
import ArrowRight from "../../assets/icons/arrow-right.svg";
import Button from "../../components/ui/Button";
import {WALLET} from "../../constants/routes";
import {RouteComponentProps, withRouter} from "react-router";
import {connect} from "react-redux";
import {REDUX_STATE} from "../../ducks";
import {AssetType} from "../../ducks/assets";
import Dropdown from "../../components/ui/Dropdown";
import Input from "../../components/ui/Input";
import CheckIcon from "../../assets/icons/check-green.svg";

enum FeeType {
  Slow = 'slow',
  Normal = 'normal',
  Fast = 'fast',
}

const FeeOptions = [
  { label: 'Slow - 0.002 DEMO', value: FeeType.Slow},
  { label: 'Normal - 0.004 DEMO', value: FeeType.Normal},
  { label: 'Fast - 0.013 DEMO', value: FeeType.Fast},
];

type StateProps = {
  assets: {
    [k: string]: AssetType
  }
  address: string
}

type Props = StateProps & RouteComponentProps

type State = {
  selectedAsset: string | null
  recipientAddress: string
  selectedFee: FeeType
  amount: string
  isReviewing: boolean
  isCompleted: boolean
  isSending: boolean
}

class Transfer extends Component<Props, State> {
  state = {
    selectedAsset: '',
    recipientAddress: '',
    selectedFee: FeeType.Normal,
    amount: '',
    isReviewing: false,
    isCompleted: false,
    isSending: false,
  };

  isValid(): boolean {
    const {
      selectedFee,
      selectedAsset,
      amount,
      recipientAddress,
      isSending,
    } = this.state;

    return !!selectedAsset.length && !!selectedFee && !!amount && !!recipientAddress && !isSending;
  }

  send = () => {
    this.setState({ isSending: true });
    setTimeout(() => this.setState({
      isSending: false,
      isReviewing: false,
      isCompleted: true,
    }), 2000);
  };

  render(): ReactNode {
    return (
      <div className="wallet__content transfer">
        { this.renderHeader() }
        { this.renderContent() }
      </div>
    );
  }

  renderHeader(): ReactNode {
    const { history } = this.props;
    const { isReviewing, isCompleted } = this.state;

    if (isCompleted) {
      return (
        <div className="wallet__content__header">
          <span>My Wallet</span>
          <Icon className="wallet__content__header__arrow" url={ArrowRight} />
          <span>Transfer</span>
          <Icon className="wallet__content__header__arrow" url={ArrowRight} />
          <span>Review</span>
          <Icon className="wallet__content__header__arrow" url={ArrowRight} />
          <span>Completed</span>
        </div>
      );
    }

    if (isReviewing) {
      return (
        <div className="wallet__content__header">
          <Button type="link" onClick={() => history.push(WALLET)}>
            My Wallet
          </Button>
          <Icon className="wallet__content__header__arrow" url={ArrowRight} />
          <Button type="link" onClick={() => this.setState({ isReviewing: false })}>
            Transfer
          </Button>
          <Icon className="wallet__content__header__arrow" url={ArrowRight} />
          <span>Review</span>
        </div>
      );
    }

    return (
      <div className="wallet__content__header">
        <Button type="link" onClick={() => history.push(WALLET)}>
          My Wallet
        </Button>
        <Icon className="wallet__content__header__arrow" url={ArrowRight} />
        <span>Transfer</span>
      </div>
    );
  }

  renderContent(): ReactNode {
    const { isReviewing, isCompleted } = this.state;

    if (isCompleted) {
      return this.renderCompleted();
    }

    if (isReviewing) {
      return this.renderReview();
    }

    return this.renderSend();
  }

  renderCompleted(): ReactNode {
    const { assets, history } = this.props;
    const {
      selectedAsset,
      recipientAddress,
      amount,
    } = this.state;
    const asset = assets[selectedAsset] || {};

    return (
      <div className="wallet__content__body deposit deposit--done">
        <div className="deposit__body">
          <div className="deposit__body__hero">
            <div className="deposit__body__icon-wrapper">
              <Icon url={CheckIcon} width={60} height={60} />
            </div>
          </div>
          <div className="deposit__body__text">
            {`Your have initiated a transfer of ${amount} ${asset.symbol} to ${recipientAddress}. The transaction will take about 5 minutes to complete.`}
          </div>
          <div className="deposit__footer">
            <Button
              type="primary"
              onClick={() => history.push(WALLET)}
            >
              Close
            </Button>
          </div>
        </div>
      </div>
    );
  }


  renderReview(): ReactNode {
    const { assets } = this.props;
    const {
      selectedAsset,
      selectedFee,
      recipientAddress,
      amount,
      isSending,
    } = this.state;
    const asset = assets[selectedAsset] || {};

    const feedIndex = FeeOptions.findIndex(({ value }) => value === selectedFee);

    return (
      <div className="transfer__form">
        { this.renderReviewRow('Asset Type', `${asset.symbol} - ${asset.name}`)}
        { this.renderReviewRow('Recipient Address', recipientAddress)}
        { this.renderReviewRow('Fee', FeeOptions[feedIndex].label)}
        { this.renderReviewRow('Amount', `${amount} ${asset.symbol}`)}
        <div className="transfer__form__actions">
          <Button
            type="primary"
            disabled={!this.isValid()}
            onClick={this.send}
            loading={isSending}
          >
            Transfer
          </Button>
        </div>
      </div>
    );
  }

  renderSend(): ReactNode {
    const { assets } = this.props;
    const {
      selectedAsset,
      selectedFee,
      recipientAddress,
      amount,
    } = this.state;
    const assetItems = Object.entries(assets)
      .filter(([ _, asset ]) => asset.sources.length)
      .map(([ assetId, asset ]) => ({
        label: `${asset.symbol} - ${asset.name}`,
        value: assetId,
      }));

    const currentIndex = assetItems.findIndex(({ value }) => value === selectedAsset);
    const feedIndex = FeeOptions.findIndex(({ value }) => value === selectedFee);

    return (
      <div className="transfer__form">
        <div className="transfer__form__row">
          <div className="transfer__form__label">Asset Type</div>
          <Dropdown
            className="transfer__form__content transfer__asset-dropdown"
            items={assetItems}
            currentIndex={currentIndex}
            onSelect={id => this.setState({ selectedAsset: id })}
          />
        </div>
        <div className="transfer__form__row">
          <div className="transfer__form__label">Recipient Address</div>
          <Input
            type="text"
            className="transfer__form__content transfer__input"
            onChange={e => this.setState({ recipientAddress: e.target.value })}
            value={recipientAddress}
            placeholder="cosmos1j689jv..."
          />
        </div>
        <div className="transfer__form__row">
          <div className="transfer__form__label">Amount</div>
          <Input
            type="number"
            className="transfer__form__content transfer__input"
            step="0.01"
            onChange={e => this.setState({ amount: e.target.value })}
            value={amount}
            placeholder="0.01"
          />
        </div>
        <div className="transfer__form__actions">
          <Button
            type="primary"
            disabled={!this.isValid()}
            onClick={() => this.setState({ isReviewing: true })}
          >
            Review
          </Button>
        </div>
      </div>
    );
  }

  renderReviewRow(label: string, value: string) {
    return (
      <div className="transfer__review__row">
        <div className="transfer__review__row__label">{label}</div>
        <div className="transfer__review__row__value">{value}</div>
      </div>
    )
  }
}

function mapStateToProps(state: REDUX_STATE): StateProps {
  const {
    assets: { assets },
    user: { address },
  } = state;

  return {
    assets,
    address,
  };
}

export default withRouter(
  connect(mapStateToProps)(Transfer)
);
