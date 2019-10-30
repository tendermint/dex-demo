import React, {Component, ReactNode} from "react";
import {BaseModalProps, Modal} from "../../pages/Modals";
import marked from "marked";
import "./info-modal.scss";
import Button from "../ui/Button";
import Icon from "../ui/Icon";
import OpenLinkIcon from "../../assets/icons/open-link.svg";
import CancelIcon from "../../assets/icons/cancel.svg";

type Props = BaseModalProps & {
  title: string
  documentationUrl: string
  markdownUrl: string
}

type State = {
  markdownHTML: string
}

export default class InfoModal extends Component<Props, State> {
  state = {
    markdownHTML: '',
  };

  async componentWillMount() {
    const resp = await fetch(this.props.markdownUrl);
    const markdown = await resp.text();
    this.setState({
      markdownHTML: marked(markdown),
    });
  }

  render (): ReactNode {
    const { onClose, title } = this.props;
    const { markdownHTML } = this.state;
    return (
      <Modal
        className="info-modal"
        onClose={onClose}
      >
        <div
          className="info-modal__container"
        >
          <div className="info-modal__header">
            <div className="info-modal__title">
              <div>{title}</div>
              <Icon
                url={CancelIcon}
                onClick={onClose}
              />
            </div>
            <Button type="link">
              Go to Documentation
              <Icon url={OpenLinkIcon} />
            </Button>
          </div>
          <div
            className="info-modal__body"
            dangerouslySetInnerHTML={{
              __html: markdownHTML,
            }}
          />
        </div>
      </Modal>
    )
  }
}