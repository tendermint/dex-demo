import React, {Component, MouseEvent} from 'react';
import ReactDOM from 'react-dom';
import c from 'classnames';
import "./modal.scss";

const modalRoot = document.getElementById('modal-root');

export type BaseModalProps = {
  onClose?: () => void
  className?: string
}

export class Modal extends Component<BaseModalProps> {
  el: HTMLDivElement
  constructor(props: BaseModalProps) {
    super(props);
    this.el = document.createElement('div');
    this.el.className = 'modal__root'
  }

  componentDidMount() {
    // The portal element is inserted in the DOM tree after
    // the Modal's children are mounted, meaning that children
    // will be mounted on a detached DOM node. If a child
    // component requires to be attached to the DOM tree
    // immediately when mounted, for example to measure a
    // DOM node, or uses 'autoFocus' in a descendant, add
    // state to Modal and only render the children when Modal
    // is inserted in the DOM tree.
    modalRoot && modalRoot.appendChild(this.el);
    window.addEventListener('keyup', this.onEscape);
  }

  componentWillUnmount() {
    modalRoot && modalRoot.removeChild(this.el);
    window.removeEventListener('keyup', this.onEscape);
  }

  onEscape = (e: KeyboardEvent) => {
    if (e.key === 'Escape' && this.props.onClose) {
      this.props.onClose();
    }
  };

  onOverlayClick = (e: MouseEvent) => {
    const { onClose } = this.props;
    e.stopPropagation();
    onClose && onClose();
  };

  render() {
    const {
      className,
      children,
    } = this.props;
    return ReactDOM.createPortal(
      <div className={c('modal', className)}>
        <div className="modal__overlay" onClick={this.onOverlayClick}/>
        <div className="modal__container">
          {children}
        </div>
      </div>,
      this.el,
    );
  }
}