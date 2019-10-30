import React, {ClipboardEvent, Component, ReactNode} from "react";
import {RouteComponentProps, withRouter} from "react-router";
import {connect} from "react-redux";
import c from 'classnames';
import "./confirm-seed.scss";
import Checkbox from "../../components/ui/Checkbox";
import Button from "../../components/ui/Button";
import CopyBlue from "../../assets/icons/copy-blue.svg";
import CopyIcon from "../../components/ui/CopyIcon";
import Input from "../../components/ui/Input";
import ConfirmModal from "../../components/ConfirmModal";

enum ConfirmSeedView {
  Entry,
  Download,
  Manual,
}

type StateProps = {

}

type DispatchProps = {

}

type Props = StateProps & DispatchProps & RouteComponentProps

type State = {
  currentView: ConfirmSeedView
  hasAcknowledgedDownload: boolean
  hasAcknowledgedManualCopy: boolean
  isShowingSeed: boolean
  isConfirmingSeed: boolean
  isShowingSkipModal: boolean
  seeds: string[]
}

class ConfirmSeedphrase extends Component<Props, State> {
  state = {
    currentView: ConfirmSeedView.Entry,
    hasAcknowledgedDownload: false,
    hasAcknowledgedManualCopy: false,
    isShowingSeed: false,
    isConfirmingSeed: false,
    isShowingSkipModal: false,
    seeds: [],
  };

  isValidSeed(): boolean {
    const expected = "manage broken nut cliff entire always course sorry pause avocado fiber slide";
    const { seeds } = this.state;

    return expected === seeds.join(' ');
  }

  onSeedChange (seed: string, index: number) {
    const { seeds } = this.state;

    // If seed contains space or newline, skip
    if ((/\s+/g).test(seed)) {
      return;
    }

    const newSeeds = Array(12)
      .fill('')
      .map((n, i) => seeds[i] ? seeds[i] : n);

    this.setState({
      seeds: newSeeds.map((s, i) => i === index ? seed : s),
    });
  };

  onSeedPaste = (e: ClipboardEvent<HTMLInputElement>, index: number) => {
    const { seeds } = this.state;
    const pasted = e.clipboardData.getData('Text');
    const splits = pasted.split(/\s+/g);

    if (splits.length < 2) {
      return;
    }

    const newSeeds = Array(12)
      .fill('')
      .map((n, i) => seeds[i] ? seeds[i] : n);

    newSeeds.splice(index, splits.length, ...splits);

    this.setState({
      seeds: newSeeds.slice(0, 12),
    });
  };

  render(): ReactNode {
    return (
      <div className="connect-wallet">
        <div className="connect-wallet__title">Back up Your Seedphrase</div>
        { this.renderContent() }
        { this.renderModal() }
      </div>
    );
  }

  renderModal(): ReactNode {
    const { isShowingSkipModal } = this.state;

    return isShowingSkipModal && (
      <ConfirmModal
        title="Are you sure?"
        onClose={() => this.setState({ isShowingSkipModal: false })}
      >
        You are at risk of losing your funds permanently if you do not back up your seed phrase.
      </ConfirmModal>
    );
  }

  renderContent(): ReactNode {
    switch (this.state.currentView) {
      case ConfirmSeedView.Entry:
        return this.renderEntry();
      case ConfirmSeedView.Download:
        return this.renderDownload();
      case ConfirmSeedView.Manual:
        return this.state.isConfirmingSeed
          ? this.renderConfirmCopy()
          : this.renderManualCopy();
      default:
        return this.renderEntry();
    }
  }

  renderEntry(): ReactNode {
    const {
      hasAcknowledgedDownload,
      hasAcknowledgedManualCopy,
    } = this.state;
    return (
      <div className="confirm-seed">
        <div className="confirm-seed__title">
          Please back up your seedphrase using one or more of the following options
        </div>
        <div className="confirm-seed__options">
          { this.renderOption(
            'Download Keystore',
            hasAcknowledgedDownload,
            ConfirmSeedView.Download,
          ) }
          { this.renderOption(
            'Manual Copy',
            hasAcknowledgedManualCopy,
            ConfirmSeedView.Manual,
          ) }
        </div>
        <div className="confirm-seed__actions">
          <Button
            type="link"
            onClick={() => this.setState({ isShowingSkipModal: true })}
          >
            Skip
          </Button>
          <Button
            type="primary"
            disabled={!hasAcknowledgedDownload && !hasAcknowledgedManualCopy}
          >
            Go to My Wallet
          </Button>
        </div>
      </div>
    );
  }

  renderDownload(): ReactNode {
    const { hasAcknowledgedDownload } = this.state;

    return (
      <div className="confirm-seed confirm-seed--download">
        <div className="confirm-seed__title">Download Keystore</div>
        <div className="confirm-seed__info">
          Please download the keystore file and keep it somewhere secure. Your account can be recovered using your password and the keystore file.
        </div>
        <Button
          type="primary"
          className="confirm-seed__download-button"
        >
          Download Keystore File
        </Button>
        <div className="confirm-seed__acknowledgement">
          <Checkbox
            onChange={e => this.setState({
              hasAcknowledgedDownload: e.target.checked,
            })}
            checked={hasAcknowledgedDownload}
          >
            I have downloaded the keystore file and kept it somewhere secure.
          </Checkbox>
        </div>
        <div className="confirm-seed__actions">
          <Button
            type="secondary"
            onClick={() => this.setState({
              currentView: ConfirmSeedView.Entry,
              hasAcknowledgedDownload: false,
            })}
          >
            Back
          </Button>
          <Button
            type="primary"
            disabled={!hasAcknowledgedDownload}
            onClick={() => this.setState({
              currentView: ConfirmSeedView.Entry,
            })}
          >
            Complete
          </Button>
        </div>
      </div>
    );
  }

  renderManualCopy(): ReactNode {
    const { isShowingSeed } = this.state;

    const seeds = "manage broken nut cliff entire always course sorry pause avocado fiber slide";

    return (
      <div className="confirm-seed confirm-seed--download">
        <div className="confirm-seed__title">Copy Your Seed Phrase</div>
        <div className="confirm-seed__info">
          Please copy your seed phrase and store them securely (like LastPass).
        </div>
        {
          isShowingSeed && (
            <div className="confirm-seed__seeds-actions">
              <CopyIcon copyText={seeds} icon={CopyBlue} />
            </div>
          )
        }
        {
          isShowingSeed
            ? (
              <div className="confirm-seed__seeds">
                {seeds.split(' ').map((seed, i) => (
                  <div key={i} className="confirm-seed__seed">
                    <div className="confirm-seed__seed__number">
                      {i + 1}:
                    </div>
                    <div className="confirm-seed__seed__text">
                      {seed}
                    </div>
                  </div>
                ))}
              </div>
            )
            : (
              <div className="confirm-seed__hidden-seeds">
                <Button
                  type="secondary"
                  onClick={() => this.setState({ isShowingSeed: true })}
                >
                  I am safe - Show my seed phrase
                </Button>
              </div>
            )
        }
        <div className="confirm-seed__acknowledgement">
          You will be asked to confirm your seed phrase in the next step.
        </div>
        <div className="confirm-seed__actions">
          <Button
            type="secondary"
            onClick={() => this.setState({
              currentView: ConfirmSeedView.Entry,
              hasAcknowledgedManualCopy: false,
              isShowingSeed: false,
            })}
          >
            Back
          </Button>
          <Button
            type="primary"
            disabled={!isShowingSeed}
            onClick={() => this.setState({
              isConfirmingSeed: true,
            })}
          >
            Next
          </Button>
        </div>
      </div>
    );
  }

  renderConfirmCopy(): ReactNode {
    const { seeds } = this.state;

    return (
      <div className="confirm-seed confirm-seed--download">
        <div className="confirm-seed__title">Copy Your Seed Phrase</div>
        <div className="confirm-seed__info">
          Please confirm your seed phrase by typing them one by one, or paste the entire seed phrase in cell # 1.
        </div>
        <div className="confirm-seed__seeds confirm-seed__seeds--confirm">
          {Array(12).fill('').map((seed, i) => (
            <div key={i} className="confirm-seed__seed">
              <div className="confirm-seed__seed__number">
                {i + 1}:
              </div>
              <div className="confirm-seed__seed__text">
                <Input
                  type="text"
                  autoFocus={!i}
                  onPaste={e => this.onSeedPaste(e, i)}
                  onChange={e => this.onSeedChange(e.target.value, i)}
                  value={seeds[i] || ''}
                />
              </div>
            </div>
          ))}
        </div>
        <div className="confirm-seed__actions">
          <Button
            type="secondary"
            onClick={() => this.setState({
              isConfirmingSeed: false,
              seeds: [],
            })}
          >
            Back
          </Button>
          <Button
            type="primary"
            disabled={!this.isValidSeed()}
            onClick={() => this.setState({
              isConfirmingSeed: false,
              currentView: ConfirmSeedView.Entry,
              hasAcknowledgedManualCopy: true,
            })}
          >
            Complete
          </Button>
        </div>
      </div>
    );
  }

  renderOption(label: string, checked: boolean, to: ConfirmSeedView): ReactNode {
    return (
      <div
        className={c('confirm-seed__option', {
          'confirm-seed__option--completed': checked,
        })}
        tabIndex={checked ? -1 : 0}
        onClick={() => this.setState({ currentView: to })}
      >
        <Checkbox checked={checked}>
          <div className="confirm-seed__option__title">
            {label}
          </div>
        </Checkbox>
      </div>
    );
  }
}

export default withRouter(
  connect()(ConfirmSeedphrase)
);
