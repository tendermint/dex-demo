import React, {Component, ReactNode} from "react";
import {BaseModalProps, Modal} from "../../pages/Modals";
import Button from "../ui/Button";
import './confirm-modal.scss';

type Props = BaseModalProps & {
  title: string
  children: ReactNode | ReactNode[]
}

class ConfirmModal extends Component<Props> {
  render(): ReactNode {
    const { onClose, title, children } = this.props;

    return (
      <Modal
        className="confirm-modal"
        onClose={onClose}
      >
        <div className="confirm-modal__container">
          <div className="confirm-modal__title">{title}</div>
          <div className="confirm-modal__content">{children}</div>
          <div className="confirm-modal__actions">
            <Button
              type="secondary-reverse"
              onClick={onClose}
            >
              Cancel
            </Button>
            <Button
              type="alert"
            >
              Skip
            </Button>
          </div>
        </div>
      </Modal>
    )
  }
}

export default ConfirmModal;
