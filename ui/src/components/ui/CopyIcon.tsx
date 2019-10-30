import React, {Component, ReactNode} from "react";
import copy from "copy-to-clipboard";
import Tooltip from "./Tooltip";
import Icon from "./Icon";
import CopyIconUrl from '../../assets/icons/copy.svg';

type State = {
  justCopied: boolean
}

type Props = {
  copyText: string
  icon?: string
}

export default class CopyIcon extends Component<Props, State> {
  timeout: any | null

  state = {
    justCopied: false,
  };

  componentWillUnmount(): void {
    if (this.timeout) {
      clearTimeout(this.timeout);
    }
  }

  copy = () => {
    if (this.timeout) {
      return;
    }

    copy(this.props.copyText);
    this.setState({ justCopied: true });
    this.timeout = setTimeout(() => {
      this.setState({
        justCopied: false,
      });

      this.timeout = null;
    }, 3000);
  };

  render (): ReactNode {
    const { icon } = this.props;

    return (
      <Tooltip content={this.state.justCopied ? 'Copied' : 'Copy'}>
        <Icon
          url={icon || CopyIconUrl}
          onClick={this.copy}
          tabIndex={0}
        />
      </Tooltip>
    )
  }
}