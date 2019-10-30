import React, {Component, ReactNode} from "react";
import "./pill.scss";

type Props = {
  children: ReactNode
  type: 'success' | 'default' | 'error' | 'pending' | 'need_attention'
}

export default class Pill extends Component<Props> {
  render(): ReactNode {
    const { children, type } = this.props;

    return (
      <div className={`pill pill--${type}`}>
        {children}
      </div>
    )
  }
}