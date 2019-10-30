import React, {Component, MouseEvent} from "react";
import c from 'classnames';
import './icon.scss';

type Props = {
  width?: number | string
  height?: number | string
  url: string
  className?: string
  onClick?: (e?: MouseEvent) => void
  tabIndex?: number
}

type IconStyle = {
  backgroundImage: string,
  width?: string,
  height?: string,
}

export default class Icon extends Component<Props> {
  render () {
    const { className, width, height, url, tabIndex, onClick } = this.props;
    let style: IconStyle = {
      backgroundImage: `url(${url})`,
    };

    if (width) style.width = maybePx(width);
    if (height) style.height = maybePx(height);

    return (
      <div
        className={c('icon', className, {
          'icon__tabbable': tabIndex && tabIndex > -1,
        })}
        style={style}
        onClick={onClick}
        onKeyPress={() => {
          if (typeof tabIndex !== "undefined") {
            onClick && onClick()
          }
        }}
        tabIndex={tabIndex}
      />
    );
  }
}

function maybePx(val: string | number | undefined) {
  if (typeof val === 'number') {
    return `${val}px`;
  }

  return val;
}
