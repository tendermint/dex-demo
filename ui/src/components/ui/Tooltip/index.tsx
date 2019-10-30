import React, {Component, ReactNode} from 'react';
import "./tooltip.scss";

type Props = {
  children: ReactNode
  content: ReactNode
  className?: string
}

export default class Tooltip extends Component<Props> {
  render (): ReactNode {
    const { children, content, className } = this.props;

    return (
      <div className={`tooltip ${className}`}>
        { children }
        <div className="tooltip__content">
          { content }
        </div>
      </div>
    );
  }
}
