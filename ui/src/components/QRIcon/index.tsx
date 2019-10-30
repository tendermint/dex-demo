import React, {Component, ReactNode} from "react";
import QRCode from "qrcode";
import Tooltip from "../ui/Tooltip";
import Icon from "../ui/Icon";
import QRIconUrl from '../../assets/icons/qr.svg';
import {Modal} from "../../pages/Modals";
import CancelIcon from "../../assets/icons/cancel-black.svg";
import Button from "../ui/Button";
import "./qr-icon.scss";

type State = {
  isShowingModal: boolean
  dataUrl: string
}

type Props = {
  text: string
  icon?: string
}

export default class QRIcon extends Component<Props, State> {
  state = {
    isShowingModal: false,
    dataUrl: '',
  };

  openModal = () => this.setState({ isShowingModal: true });
  closeModal = () => this.setState({ isShowingModal: false });

  async componentDidMount() {
    const dataUrl = await QRCode.toDataURL(this.props.text);
    this.setState({ dataUrl });
  }

  render (): ReactNode {
    const { icon } = this.props;

    return (
      <Tooltip
        content="Show QR Code"
        className="qr-icon__tooltip"
      >
        <Icon
          url={icon || QRIconUrl}
          onClick={this.openModal}
          tabIndex={0}
        />
        { this.renderModal() }
      </Tooltip>
    )
  }

  renderModal(): ReactNode {
    const { dataUrl } = this.state;

    if (!this.state.isShowingModal) {
      return null;
    }

    return (
      <Modal
        className="qr-modal"
        onClose={this.closeModal}
      >
        <div
          className="qr-modal__container"
        >
          <div className="qr-modal__content">
            <div
              className="qr-modal__qrcode"
              style={{
                backgroundImage: `url(${dataUrl})`,
              }}
            />
          </div>
          <div className="qr-modal__actions">
            <Button
              type="primary"
              onClick={this.closeModal}
            >
              Close
            </Button>
          </div>
        </div>
      </Modal>
    );
  }
}